package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/mail"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

var (
	_ resource.Resource                 = &userResource{}
	_ resource.ResourceWithImportState  = &userResource{}
	_ resource.ResourceWithUpgradeState = &userResource{}
)

func newUser() resource.Resource { return &userResource{} }

type userResource struct {
	client *apiClient
}

type userNameModel struct {
	GivenName  types.String `tfsdk:"given_name"`
	FamilyName types.String `tfsdk:"family_name"`
}

type userResourceModel struct {
	Id            types.String            `tfsdk:"id"`
	PrimaryEmail  types.String            `tfsdk:"primary_email"`
	Name          *userNameModel          `tfsdk:"name"`
	OrgUnitPath   types.String            `tfsdk:"org_unit_path"`
	Suspended     types.Bool              `tfsdk:"suspended"`
	Archived      types.Bool              `tfsdk:"archived"`
	CustomSchemas []userCustomSchemaModel `tfsdk:"custom_schemas"`
}

type userCustomSchemaModel struct {
	SchemaName   types.String `tfsdk:"schema_name"`
	SchemaValues types.Map    `tfsdk:"schema_values"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"id":            rsId(),
			"primary_email": schema.StringAttribute{Required: true},
			"org_unit_path": schema.StringAttribute{Optional: true, Computed: true},
			"suspended": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"archived": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
		Blocks: map[string]schema.Block{
			"name": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"given_name":  schema.StringAttribute{Required: true},
					"family_name": schema.StringAttribute{Required: true},
				},
			},
			"custom_schemas": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"schema_name": schema.StringAttribute{Required: true},
						"schema_values": schema.MapAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *apiClient, got: %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	password, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError("UUID Error", fmt.Sprintf("Unable to generate password: %s", err))
		return
	}

	user := &directory.User{
		PrimaryEmail: plan.PrimaryEmail.ValueString(),
		Name: &directory.UserName{
			GivenName:  plan.Name.GivenName.ValueString(),
			FamilyName: plan.Name.FamilyName.ValueString(),
		},
		Password:                  password,
		ChangePasswordAtNextLogin: true,
		Suspended:                 plan.Suspended.ValueBool(),
		Archived:                  plan.Archived.ValueBool(),
		ForceSendFields:           []string{"Suspended", "Archived"},
	}
	customSchemas, customSchemasState, err := customSchemasFromTerraform(ctx, svc, r.client.customerID, plan.CustomSchemas)
	if err != nil {
		resp.Diagnostics.AddError("Custom Schemas Error", err.Error())
		return
	}
	if customSchemas != nil {
		user.CustomSchemas = customSchemas
		user.ForceSendFields = append(user.ForceSendFields, "CustomSchemas")
	}
	if !plan.OrgUnitPath.IsNull() && !plan.OrgUnitPath.IsUnknown() {
		user.OrgUnitPath = plan.OrgUnitPath.ValueString()
	}

	created, err := svc.Users.Insert(user).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create user: %s", err))
		return
	}

	plan.Id = types.StringValue(created.Id)
	plan.OrgUnitPath = types.StringValue(created.OrgUnitPath)
	plan.CustomSchemas = customSchemasState
	if len(created.CustomSchemas) > 0 {
		plan.CustomSchemas, err = terraformCustomSchemasFromAPI(ctx, svc, r.client.customerID, created.CustomSchemas)
		if err != nil {
			resp.Diagnostics.AddError("Custom Schemas Error", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	user, err := svc.Users.Get(state.Id.ValueString()).Projection("full").Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read user: %s", err))
		return
	}

	state.PrimaryEmail = types.StringValue(user.PrimaryEmail)
	state.OrgUnitPath = types.StringValue(user.OrgUnitPath)
	state.Suspended = types.BoolValue(user.Suspended)
	state.Archived = types.BoolValue(user.Archived)
	if len(user.CustomSchemas) > 0 {
		customSchemas, err := terraformCustomSchemasFromAPI(ctx, svc, r.client.customerID, user.CustomSchemas)
		if err != nil {
			resp.Diagnostics.AddError("Custom Schemas Error", err.Error())
			return
		}
		state.CustomSchemas = customSchemas
	} else {
		state.CustomSchemas = nil
	}
	if user.Name != nil {
		state.Name = &userNameModel{
			GivenName:  types.StringValue(user.Name.GivenName),
			FamilyName: types.StringValue(user.Name.FamilyName),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var config userResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	user := &directory.User{
		PrimaryEmail: plan.PrimaryEmail.ValueString(),
		Name: &directory.UserName{
			GivenName:  plan.Name.GivenName.ValueString(),
			FamilyName: plan.Name.FamilyName.ValueString(),
		},
		OrgUnitPath:     plan.OrgUnitPath.ValueString(),
		Suspended:       plan.Suspended.ValueBool(),
		Archived:        plan.Archived.ValueBool(),
		ForceSendFields: []string{"Suspended", "Archived"},
	}
	customSchemas, customSchemasState, err := customSchemasFromTerraform(ctx, svc, r.client.customerID, config.CustomSchemas)
	if err != nil {
		resp.Diagnostics.AddError("Custom Schemas Error", err.Error())
		return
	}
	if customSchemas != nil {
		user.CustomSchemas = customSchemas
		user.ForceSendFields = append(user.ForceSendFields, "CustomSchemas")
	} else if len(state.CustomSchemas) > 0 {
		user.CustomSchemas = map[string]googleapi.RawMessage{}
		user.ForceSendFields = append(user.ForceSendFields, "CustomSchemas")
	}

	updated, err := svc.Users.Update(plan.Id.ValueString(), user).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update user: %s", err))
		return
	}

	plan.OrgUnitPath = types.StringValue(updated.OrgUnitPath)
	if len(updated.CustomSchemas) > 0 {
		plan.CustomSchemas, err = terraformCustomSchemasFromAPI(ctx, svc, r.client.customerID, updated.CustomSchemas)
		if err != nil {
			resp.Diagnostics.AddError("Custom Schemas Error", err.Error())
			return
		}
	} else if len(config.CustomSchemas) > 0 {
		plan.CustomSchemas = customSchemasState
	} else {
		plan.CustomSchemas = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	err = svc.Users.Delete(state.Id.ValueString()).Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete user: %s", err))
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userResource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var raw map[string]json.RawMessage
				if err := json.Unmarshal(req.RawState.JSON, &raw); err != nil {
					resp.Diagnostics.AddError("State Upgrade Error", fmt.Sprintf("Unable to parse raw state: %s", err))
					return
				}

				var id, primaryEmail, orgUnitPath string
				var suspended, archived bool

				_ = json.Unmarshal(raw["id"], &id)
				_ = json.Unmarshal(raw["primary_email"], &primaryEmail)
				_ = json.Unmarshal(raw["org_unit_path"], &orgUnitPath)
				_ = json.Unmarshal(raw["suspended"], &suspended)
				_ = json.Unmarshal(raw["archived"], &archived)

				// The old provider stores name as a list: [{"given_name":..., "family_name":...}]
				var givenName, familyName string
				if nameRaw, ok := raw["name"]; ok {
					var nameList []map[string]string
					if err := json.Unmarshal(nameRaw, &nameList); err == nil && len(nameList) > 0 {
						givenName = nameList[0]["given_name"]
						familyName = nameList[0]["family_name"]
					}
				}

				state := userResourceModel{
					Id:           types.StringValue(id),
					PrimaryEmail: types.StringValue(primaryEmail),
					OrgUnitPath:  types.StringValue(orgUnitPath),
					Suspended:    types.BoolValue(suspended),
					Archived:     types.BoolValue(archived),
					Name: &userNameModel{
						GivenName:  types.StringValue(givenName),
						FamilyName: types.StringValue(familyName),
					},
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			},
		},
	}
}

func customSchemasFromTerraform(ctx context.Context, svc *directory.Service, customerID string, customSchemas []userCustomSchemaModel) (map[string]googleapi.RawMessage, []userCustomSchemaModel, error) {
	if len(customSchemas) == 0 {
		return nil, nil, nil
	}

	apiValues := make(map[string]googleapi.RawMessage, len(customSchemas))
	stateValues := make([]userCustomSchemaModel, 0, len(customSchemas))
	for _, customSchema := range customSchemas {
		schemaName := customSchema.SchemaName.ValueString()
		schemaDef, err := svc.Schemas.Get(customerID, schemaName).Context(ctx).Do()
		if err != nil {
			return nil, nil, fmt.Errorf("unable to read schema %q: %w", schemaName, err)
		}

		fieldDefs := schemaFieldSpecMap(schemaDef)
		var configuredValues map[string]string
		diags := customSchema.SchemaValues.ElementsAs(ctx, &configuredValues, false)
		if diags.HasError() {
			return nil, nil, fmt.Errorf("unable to read custom_schemas[%q].schema_values", schemaName)
		}

		apiSchemaValues := make(map[string]any, len(configuredValues))
		stateSchemaValues := make(map[string]string, len(configuredValues))
		for fieldName, jsonValue := range configuredValues {
			fieldDef, ok := fieldDefs[fieldName]
			if !ok {
				return nil, nil, fmt.Errorf("field %q is not defined in schema %q", fieldName, schemaName)
			}

			var value any
			if err := json.Unmarshal([]byte(jsonValue), &value); err != nil {
				return nil, nil, fmt.Errorf("custom_schemas[%q].schema_values[%q] must be valid JSON: %w", schemaName, fieldName, err)
			}
			if err := validateCustomSchemaFieldValue(fieldDef, value); err != nil {
				return nil, nil, fmt.Errorf("custom_schemas[%q].schema_values[%q]: %w", schemaName, fieldName, err)
			}

			apiSchemaValues[fieldName] = customSchemaFieldAPIValue(fieldDef, value)
			canonical, err := json.Marshal(value)
			if err != nil {
				return nil, nil, fmt.Errorf("unable to encode custom_schemas[%q].schema_values[%q]: %w", schemaName, fieldName, err)
			}
			stateSchemaValues[fieldName] = string(canonical)
		}

		raw, err := json.Marshal(apiSchemaValues)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to encode custom_schemas[%q]: %w", schemaName, err)
		}
		stateMap, diags := types.MapValueFrom(ctx, types.StringType, stateSchemaValues)
		if diags.HasError() {
			return nil, nil, fmt.Errorf("unable to build custom_schemas[%q] state", schemaName)
		}

		apiValues[schemaName] = googleapi.RawMessage(raw)
		stateValues = append(stateValues, userCustomSchemaModel{
			SchemaName:   types.StringValue(schemaName),
			SchemaValues: stateMap,
		})
	}

	return apiValues, stateValues, nil
}

func terraformCustomSchemasFromAPI(ctx context.Context, svc *directory.Service, customerID string, customSchemas map[string]googleapi.RawMessage) ([]userCustomSchemaModel, error) {
	schemaNames := make([]string, 0, len(customSchemas))
	for schemaName := range customSchemas {
		schemaNames = append(schemaNames, schemaName)
	}
	sort.Strings(schemaNames)

	stateValues := make([]userCustomSchemaModel, 0, len(schemaNames))
	for _, schemaName := range schemaNames {
		schemaDef, err := svc.Schemas.Get(customerID, schemaName).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("unable to read schema %q: %w", schemaName, err)
		}

		var apiSchemaValues map[string]any
		if err := json.Unmarshal(customSchemas[schemaName], &apiSchemaValues); err != nil {
			return nil, fmt.Errorf("custom_schemas[%q] returned invalid JSON: %w", schemaName, err)
		}

		fieldDefs := schemaFieldSpecMap(schemaDef)
		stateSchemaValues := make(map[string]string, len(apiSchemaValues))
		for fieldName, apiValue := range apiSchemaValues {
			fieldDef, ok := fieldDefs[fieldName]
			if !ok {
				return nil, fmt.Errorf("field %q is not defined in schema %q", fieldName, schemaName)
			}

			value, err := terraformCustomSchemaFieldValue(fieldDef, apiValue)
			if err != nil {
				return nil, fmt.Errorf("custom_schemas[%q].schema_values[%q]: %w", schemaName, fieldName, err)
			}

			raw, err := json.Marshal(value)
			if err != nil {
				return nil, fmt.Errorf("unable to encode custom_schemas[%q].schema_values[%q]: %w", schemaName, fieldName, err)
			}
			stateSchemaValues[fieldName] = string(raw)
		}

		stateMap, diags := types.MapValueFrom(ctx, types.StringType, stateSchemaValues)
		if diags.HasError() {
			return nil, fmt.Errorf("unable to build custom_schemas[%q] state", schemaName)
		}

		stateValues = append(stateValues, userCustomSchemaModel{
			SchemaName:   types.StringValue(schemaName),
			SchemaValues: stateMap,
		})
	}

	return stateValues, nil
}

func schemaFieldSpecMap(schemaDef *directory.Schema) map[string]*directory.SchemaFieldSpec {
	fields := make(map[string]*directory.SchemaFieldSpec, len(schemaDef.Fields))
	for _, field := range schemaDef.Fields {
		fields[field.FieldName] = field
	}
	return fields
}

func customSchemaFieldAPIValue(fieldDef *directory.SchemaFieldSpec, value any) any {
	if !fieldDef.MultiValued {
		return value
	}

	values := value.([]any)
	apiValues := make([]map[string]any, 0, len(values))
	for _, item := range values {
		apiValues = append(apiValues, map[string]any{"value": item})
	}
	return apiValues
}

func terraformCustomSchemaFieldValue(fieldDef *directory.SchemaFieldSpec, apiValue any) (any, error) {
	if !fieldDef.MultiValued {
		return convertCustomSchemaAPIValue(fieldDef.FieldType, apiValue)
	}

	values, ok := apiValue.([]any)
	if !ok {
		return nil, fmt.Errorf("expected array for multi-valued field")
	}
	result := make([]any, 0, len(values))
	for _, item := range values {
		itemMap, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object with value for multi-valued field item")
		}
		value, ok := itemMap["value"]
		if !ok {
			return nil, fmt.Errorf("expected value for multi-valued field item")
		}
		converted, err := convertCustomSchemaAPIValue(fieldDef.FieldType, value)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

func validateCustomSchemaFieldValue(fieldDef *directory.SchemaFieldSpec, value any) error {
	if !fieldDef.MultiValued {
		return validateCustomSchemaSingleValue(fieldDef.FieldType, value)
	}

	values, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected JSON array for multi-valued %s field", fieldDef.FieldType)
	}
	for _, item := range values {
		if err := validateCustomSchemaSingleValue(fieldDef.FieldType, item); err != nil {
			return err
		}
	}
	return nil
}

func validateCustomSchemaSingleValue(fieldType string, value any) error {
	switch fieldType {
	case "BOOL":
		if _, ok := value.(bool); ok {
			return nil
		}
	case "DATE":
		date, ok := value.(string)
		if ok {
			_, err := time.Parse("2006-01-02", date)
			if err == nil {
				return nil
			}
		}
	case "DOUBLE":
		if _, ok := value.(float64); ok {
			return nil
		}
	case "EMAIL":
		email, ok := value.(string)
		if ok {
			_, err := mail.ParseAddress(email)
			if err == nil {
				return nil
			}
		}
	case "INT64":
		number, ok := value.(float64)
		if ok && math.Trunc(number) == number {
			return nil
		}
	case "PHONE", "STRING":
		if _, ok := value.(string); ok {
			return nil
		}
	}
	return fmt.Errorf("expected %s value", fieldType)
}

func convertCustomSchemaAPIValue(fieldType string, value any) (any, error) {
	text, ok := value.(string)
	if !ok {
		return value, nil
	}

	switch fieldType {
	case "BOOL":
		return strconv.ParseBool(text)
	case "DOUBLE":
		return strconv.ParseFloat(text, 64)
	case "INT64":
		return strconv.ParseInt(text, 10, 64)
	default:
		return value, nil
	}
}

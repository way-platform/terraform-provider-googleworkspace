package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

var (
	_ resource.Resource                = &schemaResource{}
	_ resource.ResourceWithImportState = &schemaResource{}
)

func newSchema() resource.Resource { return &schemaResource{} }

type schemaResource struct {
	client *apiClient
}

type schemaResourceModel struct {
	Id          types.String       `tfsdk:"id"`
	SchemaName  types.String       `tfsdk:"schema_name"`
	DisplayName types.String       `tfsdk:"display_name"`
	Fields      []schemaFieldModel `tfsdk:"field"`
}

type schemaFieldModel struct {
	FieldName      types.String `tfsdk:"field_name"`
	FieldType      types.String `tfsdk:"field_type"`
	DisplayName    types.String `tfsdk:"display_name"`
	MultiValued    types.Bool   `tfsdk:"multi_valued"`
	Indexed        types.Bool   `tfsdk:"indexed"`
	ReadAccessType types.String `tfsdk:"read_access_type"`
}

func (r *schemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *schemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": rsId(),
			"schema_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{Optional: true, Computed: true},
		},
		Blocks: map[string]schema.Block{
			"field": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"field_name": schema.StringAttribute{Required: true},
						"field_type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("STRING", "INT64", "BOOL", "DOUBLE", "EMAIL", "PHONE", "DATE"),
							},
						},
						"display_name": schema.StringAttribute{Optional: true},
						"multi_valued": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"indexed": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(true),
						},
						"read_access_type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("ALL_DOMAIN_USERS"),
							Validators: []validator.String{
								stringvalidator.OneOf("ALL_DOMAIN_USERS", "ADMINS_AND_SELF"),
							},
						},
					},
				},
			},
		},
	}
}

func (r *schemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *schemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	created, err := svc.Schemas.Insert(r.client.customerID, schemaFromModel(plan)).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create schema: %s", err))
		return
	}

	updateSchemaModelFromAPI(&plan, created)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *schemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	s, err := svc.Schemas.Get(r.client.customerID, state.Id.ValueString()).Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read schema: %s", err))
		return
	}

	updateSchemaModelFromAPI(&state, s)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *schemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	updated, err := svc.Schemas.Update(r.client.customerID, plan.Id.ValueString(), schemaFromModel(plan)).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update schema: %s", err))
		return
	}

	updateSchemaModelFromAPI(&plan, updated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *schemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	err = svc.Schemas.Delete(r.client.customerID, state.Id.ValueString()).Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete schema: %s", err))
		return
	}
}

func (r *schemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func schemaFromModel(model schemaResourceModel) *directory.Schema {
	return &directory.Schema{
		SchemaName:  model.SchemaName.ValueString(),
		DisplayName: model.DisplayName.ValueString(),
		Fields:      schemaFieldsFromModel(model.Fields),
	}
}

func schemaFieldsFromModel(fields []schemaFieldModel) []*directory.SchemaFieldSpec {
	specs := make([]*directory.SchemaFieldSpec, 0, len(fields))
	for _, field := range fields {
		indexed := field.Indexed.ValueBool()
		spec := &directory.SchemaFieldSpec{
			FieldName:       field.FieldName.ValueString(),
			FieldType:       field.FieldType.ValueString(),
			DisplayName:     field.DisplayName.ValueString(),
			MultiValued:     field.MultiValued.ValueBool(),
			Indexed:         &indexed,
			ReadAccessType:  field.ReadAccessType.ValueString(),
			ForceSendFields: []string{"MultiValued"},
		}
		specs = append(specs, spec)
	}
	return specs
}

func updateSchemaModelFromAPI(model *schemaResourceModel, s *directory.Schema) {
	model.Id = types.StringValue(s.SchemaId)
	model.SchemaName = types.StringValue(s.SchemaName)
	model.DisplayName = types.StringValue(s.DisplayName)
	model.Fields = make([]schemaFieldModel, 0, len(s.Fields))
	for _, field := range s.Fields {
		indexed := types.BoolValue(true)
		if field.Indexed != nil {
			indexed = types.BoolValue(*field.Indexed)
		}
		readAccessType := "ALL_DOMAIN_USERS"
		if field.ReadAccessType != "" {
			readAccessType = field.ReadAccessType
		}
		model.Fields = append(model.Fields, schemaFieldModel{
			FieldName:      types.StringValue(field.FieldName),
			FieldType:      types.StringValue(field.FieldType),
			DisplayName:    optionalStringValue(field.DisplayName),
			MultiValued:    types.BoolValue(field.MultiValued),
			Indexed:        indexed,
			ReadAccessType: types.StringValue(readAccessType),
		})
	}
}

func optionalStringValue(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

package provider

import (
	"context"
	"fmt"

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
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
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
	Id           types.String   `tfsdk:"id"`
	PrimaryEmail types.String   `tfsdk:"primary_email"`
	Name         *userNameModel `tfsdk:"name"`
	OrgUnitPath  types.String   `tfsdk:"org_unit_path"`
	Suspended    types.Bool     `tfsdk:"suspended"`
	Archived     types.Bool     `tfsdk:"archived"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
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
	if !plan.OrgUnitPath.IsNull() && !plan.OrgUnitPath.IsUnknown() {
		user.OrgUnitPath = plan.OrgUnitPath.ValueString()
	}

	created, err := svc.Users.Insert(user).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create user: %s", err))
		return
	}

	plan.Id = types.StringValue(created.Id)
	plan.OrgUnitPath = types.StringValue(created.OrgUnitPath)

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

	user, err := svc.Users.Get(state.Id.ValueString()).Projection("full").Do()
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
		OrgUnitPath: plan.OrgUnitPath.ValueString(),
		Suspended:   plan.Suspended.ValueBool(),
		Archived:    plan.Archived.ValueBool(),
		ForceSendFields: []string{"Suspended", "Archived"},
	}

	updated, err := svc.Users.Update(plan.Id.ValueString(), user).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update user: %s", err))
		return
	}

	plan.OrgUnitPath = types.StringValue(updated.OrgUnitPath)

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

	err = svc.Users.Delete(state.Id.ValueString()).Do()
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

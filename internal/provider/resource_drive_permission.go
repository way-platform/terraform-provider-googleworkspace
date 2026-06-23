package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var (
	_ resource.Resource                = &drivePermissionResource{}
	_ resource.ResourceWithImportState = &drivePermissionResource{}
)

func newDrivePermission() resource.Resource { return &drivePermissionResource{} }

type drivePermissionResource struct {
	client *apiClient
}

type drivePermissionResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	PermissionId         types.String `tfsdk:"permission_id"`
	FileId               types.String `tfsdk:"file_id"`
	EmailAddress         types.String `tfsdk:"email_address"`
	Role                 types.String `tfsdk:"role"`
	Type                 types.String `tfsdk:"type"`
	UseDomainAdminAccess types.Bool   `tfsdk:"use_domain_admin_access"`
}

func (r *drivePermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_drive_permission"
}

func (r *drivePermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": rsId(),
			"permission_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email_address": schema.StringAttribute{
				Optional: true,
			},
			"role": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"use_domain_admin_access": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (r *drivePermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *drivePermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan drivePermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	perm := &drive.Permission{
		Role:         plan.Role.ValueString(),
		Type:         plan.Type.ValueString(),
		EmailAddress: plan.EmailAddress.ValueString(),
	}

	created, err := driveSvc.Permissions.Create(plan.FileId.ValueString(), perm).
		UseDomainAdminAccess(plan.UseDomainAdminAccess.ValueBool()).
		SupportsAllDrives(true).
		Fields("id").
		Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create permission: %s", err))
		return
	}

	plan.PermissionId = types.StringValue(created.Id)
	plan.Id = types.StringValue(plan.FileId.ValueString() + "/" + created.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *drivePermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state drivePermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	perm, err := driveSvc.Permissions.Get(state.FileId.ValueString(), state.PermissionId.ValueString()).
		UseDomainAdminAccess(state.UseDomainAdminAccess.ValueBool()).
		SupportsAllDrives(true).
		Fields("id,emailAddress,role,type").
		Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read permission: %s", err))
		return
	}

	state.PermissionId = types.StringValue(perm.Id)
	state.EmailAddress = types.StringValue(perm.EmailAddress)
	state.Role = types.StringValue(perm.Role)
	state.Type = types.StringValue(perm.Type)
	state.Id = types.StringValue(state.FileId.ValueString() + "/" + perm.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *drivePermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan drivePermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	perm := &drive.Permission{
		Role: plan.Role.ValueString(),
	}

	updated, err := driveSvc.Permissions.Update(plan.FileId.ValueString(), plan.PermissionId.ValueString(), perm).
		UseDomainAdminAccess(plan.UseDomainAdminAccess.ValueBool()).
		SupportsAllDrives(true).
		Fields("id,emailAddress,role,type").
		Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update permission: %s", err))
		return
	}

	plan.PermissionId = types.StringValue(updated.Id)
	plan.EmailAddress = types.StringValue(updated.EmailAddress)
	plan.Role = types.StringValue(updated.Role)
	plan.Type = types.StringValue(updated.Type)
	plan.Id = types.StringValue(plan.FileId.ValueString() + "/" + updated.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *drivePermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state drivePermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	err = driveSvc.Permissions.Delete(state.FileId.ValueString(), state.PermissionId.ValueString()).
		UseDomainAdminAccess(state.UseDomainAdminAccess.ValueBool()).
		SupportsAllDrives(true).
		Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete permission: %s", err))
		return
	}
}

func (r *drivePermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Format: use_domain_admin_access,file_id/permission_id
	parts := strings.SplitN(req.ID, ",", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Expected format: use_domain_admin_access,file_id/permission_id. Got: %q", req.ID))
		return
	}

	boolVal := parts[0] == "true"
	idParts := strings.SplitN(parts[1], "/", 2)
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Expected file_id/permission_id in %q", parts[1]))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("use_domain_admin_access"), boolVal)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("file_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("permission_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

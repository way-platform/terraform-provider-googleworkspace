package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var (
	_ resource.Resource                = &driveResource{}
	_ resource.ResourceWithImportState = &driveResource{}
)

var driveRestrictionRetryDelay = func(attempt int) time.Duration {
	return time.Duration(attempt+1) * 2 * time.Second
}

func newDrive() resource.Resource { return &driveResource{} }

type driveResource struct {
	client *apiClient
}

type driveRestrictionsModel struct {
	AdminManagedRestrictions                  types.Bool `tfsdk:"admin_managed_restrictions"`
	CopyRequiresWriterPermission              types.Bool `tfsdk:"copy_requires_writer_permission"`
	DomainUsersOnly                           types.Bool `tfsdk:"domain_users_only"`
	DriveMembersOnly                          types.Bool `tfsdk:"drive_members_only"`
	SharingFoldersRequiresOrganizerPermission types.Bool `tfsdk:"sharing_folders_requires_organizer_permission"`
}

type driveResourceModel struct {
	Restrictions         *driveRestrictionsModel `tfsdk:"restrictions"`
	Name                 types.String            `tfsdk:"name"`
	Id                   types.String            `tfsdk:"id"`
	UseDomainAdminAccess types.Bool              `tfsdk:"use_domain_admin_access"`
}

func (r *driveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_drive"
}

func (r *driveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": rsId(),
			"name": schema.StringAttribute{
				Required: true,
			},
			"use_domain_admin_access": schema.BoolAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"restrictions": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"admin_managed_restrictions": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"copy_requires_writer_permission": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"domain_users_only": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"drive_members_only": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"sharing_folders_requires_organizer_permission": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

func (r *driveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *driveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan driveResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError("UUID Error", fmt.Sprintf("Unable to generate request ID: %s", err))
		return
	}

	driveReq := &drive.Drive{
		Name: plan.Name.ValueString(),
	}

	created, err := driveSvc.Drives.Create(requestId, driveReq).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create drive: %s", err))
		return
	}

	plan.Id = types.StringValue(created.Id)

	if plan.Restrictions != nil {
		updateReq := &drive.Drive{
			Restrictions: &drive.DriveRestrictions{
				AdminManagedRestrictions:                  plan.Restrictions.AdminManagedRestrictions.ValueBool(),
				CopyRequiresWriterPermission:              plan.Restrictions.CopyRequiresWriterPermission.ValueBool(),
				DomainUsersOnly:                           plan.Restrictions.DomainUsersOnly.ValueBool(),
				DriveMembersOnly:                          plan.Restrictions.DriveMembersOnly.ValueBool(),
				SharingFoldersRequiresOrganizerPermission: plan.Restrictions.SharingFoldersRequiresOrganizerPermission.ValueBool(),
				ForceSendFields: []string{
					"AdminManagedRestrictions",
					"CopyRequiresWriterPermission",
					"DomainUsersOnly",
					"DriveMembersOnly",
					"SharingFoldersRequiresOrganizerPermission",
				},
			},
		}
		// Retry with backoff: the Drive API has eventual consistency after creation.
		for attempt := 0; ; attempt++ {
			_, err = driveSvc.Drives.Update(created.Id, updateReq).
				Fields("id,name,restrictions").
				Do()
			if err == nil {
				break
			}
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 && attempt < 5 {
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Context Cancelled", "Operation cancelled while waiting for drive to become available")
					return
				case <-time.After(driveRestrictionRetryDelay(attempt)):
				}
				continue
			}
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update drive restrictions: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *driveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state driveResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	d, err := driveSvc.Drives.Get(state.Id.ValueString()).
		UseDomainAdminAccess(state.UseDomainAdminAccess.ValueBool()).
		Fields("id,name,restrictions").
		Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read drive: %s", err))
		return
	}

	state.Name = types.StringValue(d.Name)
	state.Id = types.StringValue(d.Id)

	if d.Restrictions != nil {
		hasRestrictions := d.Restrictions.AdminManagedRestrictions ||
			d.Restrictions.CopyRequiresWriterPermission ||
			d.Restrictions.DomainUsersOnly ||
			d.Restrictions.DriveMembersOnly ||
			d.Restrictions.SharingFoldersRequiresOrganizerPermission
		if hasRestrictions || state.Restrictions != nil {
			state.Restrictions = &driveRestrictionsModel{
				AdminManagedRestrictions:                  types.BoolValue(d.Restrictions.AdminManagedRestrictions),
				CopyRequiresWriterPermission:              types.BoolValue(d.Restrictions.CopyRequiresWriterPermission),
				DomainUsersOnly:                           types.BoolValue(d.Restrictions.DomainUsersOnly),
				DriveMembersOnly:                          types.BoolValue(d.Restrictions.DriveMembersOnly),
				SharingFoldersRequiresOrganizerPermission: types.BoolValue(d.Restrictions.SharingFoldersRequiresOrganizerPermission),
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *driveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan driveResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	driveReq := &drive.Drive{
		Name: plan.Name.ValueString(),
	}
	if plan.Restrictions != nil {
		driveReq.Restrictions = &drive.DriveRestrictions{
			AdminManagedRestrictions:                  plan.Restrictions.AdminManagedRestrictions.ValueBool(),
			CopyRequiresWriterPermission:              plan.Restrictions.CopyRequiresWriterPermission.ValueBool(),
			DomainUsersOnly:                           plan.Restrictions.DomainUsersOnly.ValueBool(),
			DriveMembersOnly:                          plan.Restrictions.DriveMembersOnly.ValueBool(),
			SharingFoldersRequiresOrganizerPermission: plan.Restrictions.SharingFoldersRequiresOrganizerPermission.ValueBool(),
			ForceSendFields: []string{
				"AdminManagedRestrictions",
				"CopyRequiresWriterPermission",
				"DomainUsersOnly",
				"DriveMembersOnly",
				"SharingFoldersRequiresOrganizerPermission",
			},
		}
	}

	d, err := driveSvc.Drives.Update(plan.Id.ValueString(), driveReq).
		UseDomainAdminAccess(plan.UseDomainAdminAccess.ValueBool()).
		Fields("id,name,restrictions").
		Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update drive: %s", err))
		return
	}

	plan.Id = types.StringValue(d.Id)
	if d.Restrictions != nil {
		plan.Restrictions = &driveRestrictionsModel{
			AdminManagedRestrictions:                  types.BoolValue(d.Restrictions.AdminManagedRestrictions),
			CopyRequiresWriterPermission:              types.BoolValue(d.Restrictions.CopyRequiresWriterPermission),
			DomainUsersOnly:                           types.BoolValue(d.Restrictions.DomainUsersOnly),
			DriveMembersOnly:                          types.BoolValue(d.Restrictions.DriveMembersOnly),
			SharingFoldersRequiresOrganizerPermission: types.BoolValue(d.Restrictions.SharingFoldersRequiresOrganizerPermission),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *driveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state driveResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	err = driveSvc.Drives.Delete(state.Id.ValueString()).
		UseDomainAdminAccess(state.UseDomainAdminAccess.ValueBool()).
		Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete drive: %s", err))
		return
	}
}

func (r *driveResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSplitId(ctx, req, resp, "use_domain_admin_access", "id")
}

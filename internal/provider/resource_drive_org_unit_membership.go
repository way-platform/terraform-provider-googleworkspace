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
	cibeta "google.golang.org/api/cloudidentity/v1beta1"
	"google.golang.org/api/googleapi"
)

var (
	_ resource.Resource                = &driveOrgUnitMembershipResource{}
	_ resource.ResourceWithImportState = &driveOrgUnitMembershipResource{}
)

func newDriveOrgUnitMembership() resource.Resource { return &driveOrgUnitMembershipResource{} }

type driveOrgUnitMembershipResource struct {
	client *apiClient
}

type driveOrgUnitMembershipResourceModel struct {
	Id        types.String `tfsdk:"id"`
	DriveId   types.String `tfsdk:"drive_id"`
	OrgUnitId types.String `tfsdk:"org_unit_id"`
}

func (r *driveOrgUnitMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_drive_org_unit_membership"
}

func (r *driveOrgUnitMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": rsId(),
			"drive_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org_unit_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *driveOrgUnitMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *driveOrgUnitMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan driveOrgUnitMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.moveDriveToOU(ctx, plan.DriveId.ValueString(), plan.OrgUnitId.ValueString()); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to move drive to org unit: %s", err))
		return
	}

	plan.Id = types.StringValue(plan.DriveId.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *driveOrgUnitMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state driveOrgUnitMembershipResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	driveSvc, err := r.client.NewDriveService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
		return
	}

	d, err := driveSvc.Drives.Get(state.DriveId.ValueString()).
		UseDomainAdminAccess(true).
		Fields("orgUnitId").
		Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read drive: %s", err))
		return
	}

	if d.OrgUnitId == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.OrgUnitId = types.StringValue(d.OrgUnitId)
	state.Id = types.StringValue(state.DriveId.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *driveOrgUnitMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan driveOrgUnitMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.moveDriveToOU(ctx, plan.DriveId.ValueString(), plan.OrgUnitId.ValueString()); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to move drive to org unit: %s", err))
		return
	}

	plan.Id = types.StringValue(plan.DriveId.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *driveOrgUnitMembershipResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *driveOrgUnitMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("drive_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *driveOrgUnitMembershipResource) moveDriveToOU(ctx context.Context, driveId, orgUnitId string) error {
	ciSvc, err := r.client.NewCloudIdentityService(ctx)
	if err != nil {
		return fmt.Errorf("creating Cloud Identity service: %w", err)
	}

	orgUnitId = strings.TrimPrefix(orgUnitId, "id:")
	name := fmt.Sprintf("orgUnits/-/memberships/shared_drive;%s", driveId)
	moveReq := &cibeta.MoveOrgMembershipRequest{
		Customer:           fmt.Sprintf("customers/%s", r.client.customerID),
		DestinationOrgUnit: fmt.Sprintf("orgUnits/%s", orgUnitId),
	}

	_, err = cibeta.NewOrgUnitsMembershipsService(ciSvc).Move(name, moveReq).Context(ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

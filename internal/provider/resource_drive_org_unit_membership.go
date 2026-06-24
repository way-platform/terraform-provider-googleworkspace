package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	directory "google.golang.org/api/admin/directory/v1"
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

	driveId := state.DriveId.ValueString()
	orgUnitId := state.OrgUnitId.ValueString()

	var found bool
	var err error
	if orgUnitId != "" {
		orgUnitId, found, err = r.findDriveOrgUnitMembershipInOrgUnit(ctx, driveId, orgUnitId)
	}
	if err == nil && !found {
		orgUnitId, found, err = r.findDriveOrgUnitMembership(ctx, driveId)
	}
	if err != nil {
		if isGoogleNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read drive org unit membership: %s", err))
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	state.OrgUnitId = types.StringValue(orgUnitId)
	state.Id = types.StringValue(driveId)
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

func (r *driveOrgUnitMembershipResource) findDriveOrgUnitMembership(ctx context.Context, driveId string) (string, bool, error) {
	dirSvc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		return "", false, fmt.Errorf("creating Directory service: %w", err)
	}

	orgUnits, err := dirSvc.Orgunits.List(r.client.customerID).
		Type("allIncludingParent").
		Fields("organizationUnits(orgUnitId)").
		Context(ctx).
		Do()
	if err != nil {
		return "", false, err
	}

	for _, orgUnit := range orgUnits.OrganizationUnits {
		orgUnitId := orgUnitID(orgUnit)
		if orgUnitId == "" {
			continue
		}

		foundOrgUnitId, found, err := r.findDriveOrgUnitMembershipInOrgUnit(ctx, driveId, orgUnitId)
		if err != nil {
			if isGoogleNotFound(err) {
				continue
			}
			return "", false, err
		}
		if found {
			return foundOrgUnitId, true, nil
		}
	}

	return "", false, nil
}

func (r *driveOrgUnitMembershipResource) findDriveOrgUnitMembershipInOrgUnit(ctx context.Context, driveId, orgUnitId string) (string, bool, error) {
	ciSvc, err := r.client.NewCloudIdentityService(ctx)
	if err != nil {
		return "", false, fmt.Errorf("creating Cloud Identity service: %w", err)
	}

	membershipsSvc := cibeta.NewOrgUnitsMembershipsService(ciSvc)
	parent := fmt.Sprintf("orgUnits/%s", normalizeOrgUnitId(orgUnitId))
	err = membershipsSvc.List(parent).
		Customer(fmt.Sprintf("customers/%s", r.client.customerID)).
		Filter("type == 'shared_drive'").
		PageSize(100).
		Fields("orgMemberships(name,member,memberUri,type),nextPageToken").
		Pages(ctx, func(page *cibeta.ListOrgMembershipsResponse) error {
			for _, membership := range page.OrgMemberships {
				if !driveOrgMembershipMatches(membership, driveId) {
					continue
				}

				orgUnitId = orgUnitIDFromOrgMembershipName(membership.Name)
				if orgUnitId == "" {
					orgUnitId = strings.TrimPrefix(parent, "orgUnits/")
				}
				return errDriveOrgUnitMembershipFound
			}
			return nil
		})
	if errors.Is(err, errDriveOrgUnitMembershipFound) {
		return orgUnitId, true, nil
	}
	if err != nil {
		return "", false, err
	}

	return "", false, nil
}

func (r *driveOrgUnitMembershipResource) moveDriveToOU(ctx context.Context, driveId, orgUnitId string) error {
	ciSvc, err := r.client.NewCloudIdentityService(ctx)
	if err != nil {
		return fmt.Errorf("creating Cloud Identity service: %w", err)
	}

	orgUnitId = normalizeOrgUnitId(orgUnitId)
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

var errDriveOrgUnitMembershipFound = errors.New("drive org unit membership found")

func driveOrgMembershipMatches(membership *cibeta.OrgMembership, driveId string) bool {
	if membership == nil {
		return false
	}

	return strings.HasSuffix(membership.Name, "/memberships/shared_drive;"+driveId) ||
		membership.Member == "//drive.googleapis.com/drives/"+driveId ||
		membership.MemberUri == "https://drive.googleapis.com/drive/v3/drives/"+driveId
}

func orgUnitIDFromOrgMembershipName(name string) string {
	name = strings.TrimPrefix(name, "orgUnits/")
	before, _, ok := strings.Cut(name, "/memberships/")
	if !ok {
		return ""
	}

	return normalizeOrgUnitId(before)
}

func orgUnitID(orgUnit *directory.OrgUnit) string {
	if orgUnit == nil {
		return ""
	}

	return normalizeOrgUnitId(orgUnit.OrgUnitId)
}

func normalizeOrgUnitId(orgUnitId string) string {
	return strings.TrimPrefix(orgUnitId, "id:")
}

func isGoogleNotFound(err error) bool {
	gerr, ok := err.(*googleapi.Error)
	return ok && gerr.Code == 404
}

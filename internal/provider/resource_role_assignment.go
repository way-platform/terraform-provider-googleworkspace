package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

var (
	_ resource.Resource                 = &roleAssignmentResource{}
	_ resource.ResourceWithImportState  = &roleAssignmentResource{}
	_ resource.ResourceWithUpgradeState = &roleAssignmentResource{}
)

func newRoleAssignment() resource.Resource { return &roleAssignmentResource{} }

type roleAssignmentResource struct {
	client *apiClient
}

type roleAssignmentResourceModel struct {
	Id         types.String `tfsdk:"id"`
	RoleId     types.String `tfsdk:"role_id"`
	AssignedTo types.String `tfsdk:"assigned_to"`
	ScopeType  types.String `tfsdk:"scope_type"`
	OrgUnitId  types.String `tfsdk:"org_unit_id"`
}

func (r *roleAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_assignment"
}

func (r *roleAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	forceNew := []planmodifier.String{stringplanmodifier.RequiresReplace()}
	resp.Schema = schema.Schema{
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"id": rsId(),
			"role_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: forceNew,
			},
			"assigned_to": schema.StringAttribute{
				Required:      true,
				PlanModifiers: forceNew,
			},
			"scope_type": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: forceNew,
			},
			"org_unit_id": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: forceNew,
			},
		},
	}
}

func (r *roleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *roleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	roleId, err := strconv.ParseInt(plan.RoleId.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Value Error", fmt.Sprintf("Invalid role_id %q: %s", plan.RoleId.ValueString(), err))
		return
	}

	scopeType := strings.ToUpper(plan.ScopeType.ValueString())
	orgUnitId := strings.TrimPrefix(plan.OrgUnitId.ValueString(), "id:")
	if scopeType == "ORG_UNIT" && orgUnitId == "" {
		resp.Diagnostics.AddError("Value Error", "org_unit_id must be set when scope_type is ORG_UNIT")
		return
	}

	assignment := &directory.RoleAssignment{
		RoleId:     roleId,
		AssignedTo: plan.AssignedTo.ValueString(),
		ScopeType:  scopeType,
		OrgUnitId:  orgUnitId,
	}

	created, err := svc.RoleAssignments.Insert(r.client.customerID, assignment).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create role assignment: %s", err))
		return
	}

	plan.Id = types.StringValue(strconv.FormatInt(created.RoleAssignmentId, 10))
	plan.ScopeType = types.StringValue(created.ScopeType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *roleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	ra, err := svc.RoleAssignments.Get(r.client.customerID, state.Id.ValueString()).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read role assignment: %s", err))
		return
	}

	state.RoleId = types.StringValue(strconv.FormatInt(ra.RoleId, 10))
	state.AssignedTo = types.StringValue(ra.AssignedTo)
	state.ScopeType = types.StringValue(ra.ScopeType)
	state.OrgUnitId = normalizedOrgUnitID(ra.OrgUnitId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *roleAssignmentResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "All attributes are ForceNew; update should never be called.")
}

func (r *roleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	err = svc.RoleAssignments.Delete(r.client.customerID, state.Id.ValueString()).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete role assignment: %s", err))
		return
	}
}

func (r *roleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *roleAssignmentResource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var raw map[string]json.RawMessage
				if err := json.Unmarshal(req.RawState.JSON, &raw); err != nil {
					resp.Diagnostics.AddError("State Upgrade Error", fmt.Sprintf("Unable to parse raw state: %s", err))
					return
				}

				var id, roleId, assignedTo, scopeType, orgUnitId string
				_ = json.Unmarshal(raw["id"], &id)
				_ = json.Unmarshal(raw["role_id"], &roleId)
				_ = json.Unmarshal(raw["assigned_to"], &assignedTo)
				_ = json.Unmarshal(raw["scope_type"], &scopeType)
				_ = json.Unmarshal(raw["org_unit_id"], &orgUnitId)

				if scopeType == "" {
					scopeType = "CUSTOMER"
				}

				state := roleAssignmentResourceModel{
					Id:         types.StringValue(id),
					RoleId:     types.StringValue(roleId),
					AssignedTo: types.StringValue(assignedTo),
					ScopeType:  types.StringValue(scopeType),
					OrgUnitId:  normalizedOrgUnitID(orgUnitId),
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			},
		},
	}
}

func normalizedOrgUnitID(id string) types.String {
	id = strings.TrimPrefix(id, "id:")
	if id == "" {
		return types.StringNull()
	}
	return types.StringValue(id)
}

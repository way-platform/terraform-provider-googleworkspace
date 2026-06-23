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
	directory "google.golang.org/api/admin/directory/v1"
)

var (
	_ resource.Resource                = &orgUnitResource{}
	_ resource.ResourceWithImportState = &orgUnitResource{}
)

func newOrgUnit() resource.Resource { return &orgUnitResource{} }

type orgUnitResource struct {
	client *apiClient
}

type orgUnitResourceModel struct {
	Id                types.String `tfsdk:"id"`
	OrgUnitPath       types.String `tfsdk:"org_unit_path"`
	Name              types.String `tfsdk:"name"`
	ParentOrgUnitPath types.String `tfsdk:"parent_org_unit_path"`
	Description       types.String `tfsdk:"description"`
}

func (r *orgUnitResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_unit"
}

func (r *orgUnitResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            rsId(),
			"org_unit_path": schema.StringAttribute{Computed: true},
			"name":          schema.StringAttribute{Required: true},
			"parent_org_unit_path": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{Optional: true},
		},
	}
}

func (r *orgUnitResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *orgUnitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan orgUnitResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	ou := &directory.OrgUnit{
		Name:              plan.Name.ValueString(),
		ParentOrgUnitPath: plan.ParentOrgUnitPath.ValueString(),
		Description:       plan.Description.ValueString(),
	}

	created, err := svc.Orgunits.Insert(r.client.customerID, ou).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create org unit: %s", err))
		return
	}

	plan.Id = types.StringValue(strings.TrimPrefix(created.OrgUnitId, "id:"))
	plan.OrgUnitPath = types.StringValue(created.OrgUnitPath)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *orgUnitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state orgUnitResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	ou, err := svc.Orgunits.Get(r.client.customerID, orgUnitAPIPath(state.Id.ValueString())).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read org unit: %s", err))
		return
	}

	state.Name = types.StringValue(ou.Name)
	state.OrgUnitPath = types.StringValue(ou.OrgUnitPath)
	state.ParentOrgUnitPath = types.StringValue(ou.ParentOrgUnitPath)
	if ou.Description != "" {
		state.Description = types.StringValue(ou.Description)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *orgUnitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan orgUnitResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	ou := &directory.OrgUnit{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	updated, err := svc.Orgunits.Update(r.client.customerID, orgUnitAPIPath(plan.Id.ValueString()), ou).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update org unit: %s", err))
		return
	}

	plan.OrgUnitPath = types.StringValue(updated.OrgUnitPath)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *orgUnitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state orgUnitResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	err = svc.Orgunits.Delete(r.client.customerID, orgUnitAPIPath(state.Id.ValueString())).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete org unit: %s", err))
		return
	}
}

func (r *orgUnitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := strings.TrimPrefix(req.ID, "id:")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

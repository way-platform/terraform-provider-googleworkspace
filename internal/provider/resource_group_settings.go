package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/groupssettings/v1"
)

var (
	_ resource.Resource                = &groupSettingsResource{}
	_ resource.ResourceWithImportState = &groupSettingsResource{}
)

func newGroupSettings() resource.Resource { return &groupSettingsResource{} }

type groupSettingsResource struct {
	client *apiClient
}

type groupSettingsResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Email                types.String `tfsdk:"email"`
	WhoCanViewMembership types.String `tfsdk:"who_can_view_membership"`
}

func (r *groupSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_settings"
}

func (r *groupSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": rsId(),
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"who_can_view_membership": schema.StringAttribute{Optional: true},
		},
	}
}

func (r *groupSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewGroupsSettingsService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create GroupsSettings service: %s", err))
		return
	}

	plan.Id = plan.Email

	settings := &groupssettings.Groups{
		WhoCanViewMembership: plan.WhoCanViewMembership.ValueString(),
	}

	_, err = svc.Groups.Update(plan.Email.ValueString(), settings).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update group settings: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewGroupsSettingsService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create GroupsSettings service: %s", err))
		return
	}

	settings, err := svc.Groups.Get(state.Email.ValueString()).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read group settings: %s", err))
		return
	}

	if settings.WhoCanViewMembership != "" {
		state.WhoCanViewMembership = types.StringValue(settings.WhoCanViewMembership)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *groupSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewGroupsSettingsService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create GroupsSettings service: %s", err))
		return
	}

	plan.Id = plan.Email

	settings := &groupssettings.Groups{
		WhoCanViewMembership: plan.WhoCanViewMembership.ValueString(),
	}

	_, err = svc.Groups.Update(plan.Email.ValueString(), settings).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update group settings: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupSettingsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Group settings cannot be deleted; they exist with the group. No-op.
}

func (r *groupSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), req.ID)...)
}

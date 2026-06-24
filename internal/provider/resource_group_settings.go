package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/googleapi"
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
	Id                     types.String `tfsdk:"id"`
	Email                  types.String `tfsdk:"email"`
	WhoCanViewMembership   types.String `tfsdk:"who_can_view_membership"`
	WhoCanPostMessage      types.String `tfsdk:"who_can_post_message"`
	MessageModerationLevel types.String `tfsdk:"message_moderation_level"`
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
			// who_can_post_message is mutated server-side (e.g. archive_only flips it to
			// NONE_CAN_POST / ALL_MANAGERS_CAN_POST) and the API always returns a value, so it is
			// Computed to avoid perpetual diffs when left unset.
			"who_can_post_message": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"NONE_CAN_POST",
						"ALL_MANAGERS_CAN_POST",
						"ALL_MEMBERS_CAN_POST",
						"ALL_OWNERS_CAN_POST",
						"ALL_IN_DOMAIN_CAN_POST",
						"ANYONE_CAN_POST",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"message_moderation_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("MODERATE_NONE"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"MODERATE_ALL_MESSAGES",
						"MODERATE_NON_MEMBERS",
						"MODERATE_NEW_MEMBERS",
						"MODERATE_NONE",
					),
				},
			},
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
		WhoCanViewMembership:   plan.WhoCanViewMembership.ValueString(),
		WhoCanPostMessage:      plan.WhoCanPostMessage.ValueString(),
		MessageModerationLevel: plan.MessageModerationLevel.ValueString(),
	}

	updated, err := svc.Groups.Update(plan.Email.ValueString(), settings).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update group settings: %s", err))
		return
	}

	// who_can_post_message is Computed: populate it from the API response so a value left unset in
	// the configuration resolves to the server-assigned value rather than staying unknown.
	if updated.WhoCanPostMessage != "" {
		plan.WhoCanPostMessage = types.StringValue(updated.WhoCanPostMessage)
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
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read group settings: %s", err))
		return
	}

	if settings.WhoCanViewMembership != "" {
		state.WhoCanViewMembership = types.StringValue(settings.WhoCanViewMembership)
	}
	if settings.WhoCanPostMessage != "" {
		state.WhoCanPostMessage = types.StringValue(settings.WhoCanPostMessage)
	}
	if settings.MessageModerationLevel != "" {
		state.MessageModerationLevel = types.StringValue(settings.MessageModerationLevel)
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
		WhoCanViewMembership:   plan.WhoCanViewMembership.ValueString(),
		WhoCanPostMessage:      plan.WhoCanPostMessage.ValueString(),
		MessageModerationLevel: plan.MessageModerationLevel.ValueString(),
	}

	updated, err := svc.Groups.Update(plan.Email.ValueString(), settings).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update group settings: %s", err))
		return
	}

	if updated.WhoCanPostMessage != "" {
		plan.WhoCanPostMessage = types.StringValue(updated.WhoCanPostMessage)
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

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	AllowExternalMembers   types.Bool   `tfsdk:"allow_external_members"`
	ReplyTo                types.String `tfsdk:"reply_to"`
	CustomReplyTo          types.String `tfsdk:"custom_reply_to"`
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
			// The attributes below are Computed for the same reason as who_can_post_message: when left
			// unset, the value set in the Console is kept rather than producing a perpetual diff.
			"allow_external_members": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"reply_to": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"REPLY_TO_CUSTOM",
						"REPLY_TO_SENDER",
						"REPLY_TO_LIST",
						"REPLY_TO_OWNER",
						"REPLY_TO_IGNORE",
						"REPLY_TO_MANAGERS",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// custom_reply_to is required by the API when reply_to is REPLY_TO_CUSTOM.
			"custom_reply_to": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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

	updated, err := svc.Groups.Update(plan.Email.ValueString(), plan.toAPI()).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update group settings: %s", err))
		return
	}

	// Computed attributes: populate them from the API response so a value left unset in the
	// configuration resolves to the server-assigned value rather than staying unknown.
	plan.setComputedFromAPI(updated)

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
	if settings.MessageModerationLevel != "" {
		state.MessageModerationLevel = types.StringValue(settings.MessageModerationLevel)
	}
	state.setComputedFromAPI(settings)

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

	updated, err := svc.Groups.Update(plan.Email.ValueString(), plan.toAPI()).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update group settings: %s", err))
		return
	}

	plan.setComputedFromAPI(updated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupSettingsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Group settings cannot be deleted; they exist with the group. No-op.
}

func (r *groupSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), req.ID)...)
}

// toAPI builds the settings to send. Unset Computed attributes are unknown or null and send an empty
// string, which omitempty drops, so the Update keeps their current values.
func (m *groupSettingsResourceModel) toAPI() *groupssettings.Groups {
	settings := &groupssettings.Groups{
		WhoCanViewMembership:   m.WhoCanViewMembership.ValueString(),
		WhoCanPostMessage:      m.WhoCanPostMessage.ValueString(),
		MessageModerationLevel: m.MessageModerationLevel.ValueString(),
		ReplyTo:                m.ReplyTo.ValueString(),
		CustomReplyTo:          m.CustomReplyTo.ValueString(),
	}
	// The API encodes this boolean as the string "true" or "false", so false is not dropped.
	if !m.AllowExternalMembers.IsNull() && !m.AllowExternalMembers.IsUnknown() {
		settings.AllowExternalMembers = strconv.FormatBool(m.AllowExternalMembers.ValueBool())
	}
	return settings
}

// setComputedFromAPI copies the Computed attributes from an API response. A value the response
// does not include must still be known after apply, so it resolves to null.
func (m *groupSettingsResourceModel) setComputedFromAPI(settings *groupssettings.Groups) {
	if settings.WhoCanPostMessage != "" {
		m.WhoCanPostMessage = types.StringValue(settings.WhoCanPostMessage)
	}
	switch {
	case settings.AllowExternalMembers != "":
		m.AllowExternalMembers = types.BoolValue(settings.AllowExternalMembers == "true")
	case m.AllowExternalMembers.IsUnknown():
		m.AllowExternalMembers = types.BoolNull()
	}
	m.ReplyTo = computedString(m.ReplyTo, settings.ReplyTo)
	m.CustomReplyTo = computedString(m.CustomReplyTo, settings.CustomReplyTo)
}

func computedString(current types.String, fromAPI string) types.String {
	switch {
	case fromAPI != "":
		return types.StringValue(fromAPI)
	case current.IsUnknown():
		return types.StringNull()
	default:
		return current
	}
}

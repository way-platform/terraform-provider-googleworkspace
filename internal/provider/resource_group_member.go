package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

var (
	_ resource.Resource                = &groupMemberResource{}
	_ resource.ResourceWithImportState = &groupMemberResource{}
)

func newGroupMember() resource.Resource { return &groupMemberResource{} }

// groupMemberResource manages one member of a group. Unlike group_members, it is not
// authoritative: members it does not manage are never read or removed. Never use both resources
// on the same group, because group_members removes every member it does not list.
type groupMemberResource struct {
	client *apiClient
}

type groupMemberResourceModel struct {
	Id      types.String `tfsdk:"id"`
	GroupId types.String `tfsdk:"group_id"`
	Email   types.String `tfsdk:"email"`
	Role    types.String `tfsdk:"role"`
	Type    types.String `tfsdk:"type"`
}

func (r *groupMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_member"
}

func (r *groupMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": rsId(),
			"group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("OWNER", "MANAGER", "MEMBER"),
				},
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("USER"),
				Validators: []validator.String{
					stringvalidator.OneOf("USER", "GROUP", "CUSTOMER", "EXTERNAL"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *groupMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	member := &directory.Member{
		Email: plan.Email.ValueString(),
		Role:  plan.Role.ValueString(),
		Type:  plan.Type.ValueString(),
	}
	_, err = svc.Members.Insert(plan.GroupId.ValueString(), member).Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 {
			resp.Diagnostics.AddError("Member Already Exists", fmt.Sprintf(
				"%q is already a member of group %q. Import it with ID %q.",
				plan.Email.ValueString(), plan.GroupId.ValueString(), groupMemberId(plan.GroupId.ValueString(), plan.Email.ValueString())))
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to insert member %q: %s", plan.Email.ValueString(), err))
		return
	}

	plan.Id = types.StringValue(groupMemberId(plan.GroupId.ValueString(), plan.Email.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	member, err := svc.Members.Get(state.GroupId.ValueString(), state.Email.ValueString()).Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read member %q: %s", state.Email.ValueString(), err))
		return
	}

	// email is kept from state: the API returns it lowercased, which would diff against a
	// configuration in mixed case.
	state.Role = types.StringValue(member.Role)
	state.Type = types.StringValue(member.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update changes only the role: every other attribute requires replacement.
func (r *groupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	member := &directory.Member{
		Role: plan.Role.ValueString(),
	}
	_, err = svc.Members.Update(plan.GroupId.ValueString(), plan.Email.ValueString(), member).Context(ctx).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update member %q: %s", plan.Email.ValueString(), err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	err = svc.Members.Delete(state.GroupId.ValueString(), state.Email.ValueString()).Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete member %q: %s", state.Email.ValueString(), err))
	}
}

// ImportState accepts <group_id>/<email>, with an optional groups/ prefix.
func (r *groupMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	groupId, email, ok := strings.Cut(strings.TrimPrefix(req.ID, "groups/"), "/")
	if !ok || groupId == "" || email == "" {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Expected format: group_id/email. Got: %q", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), groupMemberId(groupId, email))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), email)...)
}

func groupMemberId(groupId, email string) string {
	return groupId + "/" + email
}

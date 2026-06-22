package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	directory "google.golang.org/api/admin/directory/v1"
)

var (
	_ resource.Resource                = &groupMembersResource{}
	_ resource.ResourceWithImportState = &groupMembersResource{}
)

func newGroupMembers() resource.Resource { return &groupMembersResource{} }

type groupMembersResource struct {
	client *apiClient
}

type groupMemberModel struct {
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
	Type  types.String `tfsdk:"type"`
}

type groupMembersResourceModel struct {
	Id      types.String      `tfsdk:"id"`
	GroupId types.String      `tfsdk:"group_id"`
	Members []groupMemberModel `tfsdk:"members"`
}

func (r *groupMembersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_members"
}

func (r *groupMembersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": rsId(),
			"group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"members": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{Required: true},
						"role":  schema.StringAttribute{Required: true},
						"type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("USER"),
						},
					},
				},
			},
		},
	}
}

func (r *groupMembersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupMembersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupMembersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	plan.Id = plan.GroupId

	for _, m := range plan.Members {
		member := &directory.Member{
			Email: m.Email.ValueString(),
			Role:  m.Role.ValueString(),
			Type:  m.Type.ValueString(),
		}
		_, err := svc.Members.Insert(plan.GroupId.ValueString(), member).Do()
		if err != nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to insert member %q: %s", m.Email.ValueString(), err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupMembersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupMembersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	var members []groupMemberModel
	err = svc.Members.List(state.GroupId.ValueString()).Pages(ctx, func(page *directory.Members) error {
		for _, m := range page.Members {
			members = append(members, groupMemberModel{
				Email: types.StringValue(m.Email),
				Role:  types.StringValue(m.Role),
				Type:  types.StringValue(m.Type),
			})
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to list members: %s", err))
		return
	}

	state.Members = members

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *groupMembersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupMembersResourceModel
	var state groupMembersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	plan.Id = plan.GroupId

	currentMembers := make(map[string]groupMemberModel)
	for _, m := range state.Members {
		currentMembers[m.Email.ValueString()] = m
	}

	desiredMembers := make(map[string]groupMemberModel)
	for _, m := range plan.Members {
		desiredMembers[m.Email.ValueString()] = m
	}

	// Delete removed members
	for email := range currentMembers {
		if _, ok := desiredMembers[email]; !ok {
			err := svc.Members.Delete(plan.GroupId.ValueString(), email).Do()
			if err != nil {
				resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete member %q: %s", email, err))
				return
			}
		}
	}

	// Insert new or update changed members
	for email, desired := range desiredMembers {
		current, exists := currentMembers[email]
		if !exists {
			member := &directory.Member{
				Email: email,
				Role:  desired.Role.ValueString(),
				Type:  desired.Type.ValueString(),
			}
			_, err := svc.Members.Insert(plan.GroupId.ValueString(), member).Do()
			if err != nil {
				resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to insert member %q: %s", email, err))
				return
			}
		} else if desired.Role.ValueString() != current.Role.ValueString() {
			member := &directory.Member{
				Role: desired.Role.ValueString(),
			}
			_, err := svc.Members.Update(plan.GroupId.ValueString(), email, member).Do()
			if err != nil {
				resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update member %q: %s", email, err))
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupMembersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupMembersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	for _, m := range state.Members {
		err := svc.Members.Delete(state.GroupId.ValueString(), m.Email.ValueString()).Do()
		if err != nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete member %q: %s", m.Email.ValueString(), err))
			return
		}
	}
}

func (r *groupMembersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Accept format: groups/<group_id> or just <group_id>
	id := strings.TrimPrefix(req.ID, "groups/")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), id)...)
}

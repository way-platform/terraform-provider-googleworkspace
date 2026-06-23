package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	directory "google.golang.org/api/admin/directory/v1"
)

var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

func newGroup() resource.Resource { return &groupResource{} }

type groupResource struct {
	client *apiClient
}

type groupResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Aliases     types.List   `tfsdk:"aliases"`
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          rsId(),
			"email":       schema.StringAttribute{Required: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"aliases": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	group := &directory.Group{
		Email:       plan.Email.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	created, err := svc.Groups.Insert(group).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create group: %s", err))
		return
	}

	plan.Id = types.StringValue(created.Id)

	var desiredAliases []string
	resp.Diagnostics.Append(plan.Aliases.ElementsAs(ctx, &desiredAliases, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	for _, a := range desiredAliases {
		_, err := svc.Groups.Aliases.Insert(created.Id, &directory.Alias{Alias: a}).Do()
		if err != nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to add alias %q: %s", a, err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	group, err := svc.Groups.Get(state.Id.ValueString()).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read group: %s", err))
		return
	}

	state.Email = types.StringValue(group.Email)
	state.Name = types.StringValue(group.Name)
	if group.Description != "" {
		state.Description = types.StringValue(group.Description)
	} else {
		state.Description = types.StringNull()
	}

	if len(group.Aliases) > 0 {
		aliases, diags := types.ListValueFrom(ctx, types.StringType, group.Aliases)
		resp.Diagnostics.Append(diags...)
		state.Aliases = aliases
	} else {
		state.Aliases = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	group := &directory.Group{
		Email:       plan.Email.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	_, err = svc.Groups.Update(plan.Id.ValueString(), group).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update group: %s", err))
		return
	}

	// Reconcile aliases
	current, err := svc.Groups.Get(plan.Id.ValueString()).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read group after update: %s", err))
		return
	}

	var desiredAliases []string
	resp.Diagnostics.Append(plan.Aliases.ElementsAs(ctx, &desiredAliases, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	currentSet := make(map[string]bool)
	for _, a := range current.Aliases {
		currentSet[a] = true
	}
	desiredSet := make(map[string]bool)
	for _, a := range desiredAliases {
		desiredSet[a] = true
	}

	for _, a := range desiredAliases {
		if !currentSet[a] {
			_, err := svc.Groups.Aliases.Insert(plan.Id.ValueString(), &directory.Alias{Alias: a}).Do()
			if err != nil {
				resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to add alias %q: %s", a, err))
				return
			}
		}
	}
	for _, a := range current.Aliases {
		if !desiredSet[a] {
			err := svc.Groups.Aliases.Delete(plan.Id.ValueString(), a).Do()
			if err != nil {
				resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to remove alias %q: %s", a, err))
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := r.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	err = svc.Groups.Delete(state.Id.ValueString()).Do()
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete group: %s", err))
		return
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

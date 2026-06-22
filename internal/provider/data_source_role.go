package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	directory "google.golang.org/api/admin/directory/v1"
)

var _ datasource.DataSource = &roleDataSource{}

func newRoleDataSource() datasource.DataSource { return &roleDataSource{} }

type roleDataSource struct {
	client *apiClient
}

type roleDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Required: true},
		},
	}
}

func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*apiClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type", fmt.Sprintf("Expected *apiClient, got: %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data roleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, err := d.client.NewDirectoryService(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Directory service: %s", err))
		return
	}

	name := data.Name.ValueString()
	var found *directory.Role

	err = svc.Roles.List(d.client.customerID).Pages(ctx, func(page *directory.Roles) error {
		for _, role := range page.Items {
			if role.RoleName == name {
				found = role
			}
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to list roles: %s", err))
		return
	}

	if found == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Role %q not found", name))
		return
	}

	data.Id = types.StringValue(strconv.FormatInt(found.RoleId, 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

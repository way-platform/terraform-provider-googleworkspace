package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

var _ provider.Provider = &googleworkspaceProvider{}

// testAPIClient is set by tests to bypass authentication and inject a mock client.
var testAPIClient *apiClient

type googleworkspaceProvider struct {
	version string
}

type googleworkspaceProviderModel struct {
	AccessToken      types.String `tfsdk:"access_token"`
	ServiceAccount   types.String `tfsdk:"service_account"`
	ImpersonatedUser types.String `tfsdk:"impersonated_user_email"`
	CustomerID       types.String `tfsdk:"customer_id"`
	OAuthScopes      types.List   `tfsdk:"oauth_scopes"`
	RetryOn          types.List   `tfsdk:"retry_on"`
}

func (p *googleworkspaceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "googleworkspace"
	resp.Version = p.version
}

func (p *googleworkspaceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"access_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: `A pre-minted OAuth2 access token. When set, the provider uses
impersonate.CredentialsTokenSource with this token as the base credential to
impersonate the service_account with Domain-Wide Delegation.`,
			},
			"service_account": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: `The email address of the Service Account to impersonate for Domain-Wide Delegation.`,
			},
			"impersonated_user_email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: `The email address of the Workspace user to impersonate with Domain-Wide Delegation.`,
			},
			"customer_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: `The Google Workspace customer ID. Required for Admin SDK resources (users, groups, OUs, roles).`,
			},
			"oauth_scopes": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				MarkdownDescription: `List of OAuth scopes for the API client.
If unset, uses a default set covering Admin SDK and Drive.`,
			},
			"retry_on": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				MarkdownDescription: `HTTP error codes to retry on. Defaults to 502.
Always retries 429, 403 rate-limit errors, and 5xx (except 501).`,
			},
		},
	}
}

func (p *googleworkspaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data googleworkspaceProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// In tests, skip authentication and use the injected mock client.
	if testAPIClient != nil {
		resp.DataSourceData = testAPIClient
		resp.ResourceData = testAPIClient
		return
	}

	accessToken := data.AccessToken.ValueString()
	serviceAccount := data.ServiceAccount.ValueString()
	if serviceAccount == "" {
		serviceAccount = os.Getenv("SERVICE_ACCOUNT")
	}
	subject := data.ImpersonatedUser.ValueString()
	if subject == "" {
		subject = os.Getenv("SUBJECT")
	}
	if subject == "" {
		resp.Diagnostics.AddError("Configuration Error", "impersonated_user_email must be set")
		return
	}
	customerID := data.CustomerID.ValueString()
	if customerID == "" {
		customerID = os.Getenv("GOOGLEWORKSPACE_CUSTOMER_ID")
	}

	var scopes []string
	if data.OAuthScopes.IsNull() {
		scopes = []string{
			"https://www.googleapis.com/auth/admin.directory.group",
			"https://www.googleapis.com/auth/admin.directory.user",
			"https://www.googleapis.com/auth/admin.directory.orgunit",
			"https://www.googleapis.com/auth/admin.directory.rolemanagement",
			"https://www.googleapis.com/auth/apps.groups.settings",
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/cloud-identity.orgunits",
		}
	} else {
		resp.Diagnostics.Append(data.OAuthScopes.ElementsAs(ctx, &scopes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var retryOn []int
	if data.RetryOn.IsNull() {
		retryOn = []int{502}
	} else {
		resp.Diagnostics.Append(data.RetryOn.ElementsAs(ctx, &retryOn, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var tokenSource oauth2.TokenSource

	if serviceAccount == "" {
		resp.Diagnostics.AddError("Configuration Error", "service_account is required")
		return
	}

	switch {
	case accessToken != "":
		token := &oauth2.Token{AccessToken: accessToken}
		ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: serviceAccount,
			Scopes:          scopes,
			Subject:         subject,
		}, option.WithTokenSource(oauth2.StaticTokenSource(token)))
		if err != nil {
			resp.Diagnostics.AddError("Configuration Error", fmt.Sprintf("Unable to create impersonated credentials: %s", err))
			return
		}
		tokenSource = ts

	default:
		ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: serviceAccount,
			Scopes:          scopes,
			Subject:         subject,
		})
		if err != nil {
			resp.Diagnostics.AddError("Configuration Error", fmt.Sprintf("Unable to create impersonated credentials from ADC: %s", err))
			return
		}
		tokenSource = ts
	}

	httpClient := newRetryableClient(retryOn)
	httpClient.Transport = &oauth2.Transport{
		Source: tokenSource,
		Base:   httpClient.Transport,
	}

	client := &apiClient{
		client:     httpClient,
		customerID: customerID,
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *googleworkspaceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newDrive,
		newDriveOrgUnitMembership,
		newDrivePermission,
		newOrgUnit,
		newGroup,
		newGroupMembers,
		newGroupSettings,
		newRoleAssignment,
		newSchema,
		newUser,
	}
}

func (p *googleworkspaceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newRoleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &googleworkspaceProvider{
			version: version,
		}
	}
}

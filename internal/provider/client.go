package provider

import (
	"context"
	"net/http"

	directory "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/groupssettings/v1"
	"google.golang.org/api/option"
)

type apiClient struct {
	client     *http.Client
	customerID string
	basePath   string // empty in production; set to test server URL in tests
}

func (c *apiClient) clientOptions() []option.ClientOption {
	opts := []option.ClientOption{option.WithHTTPClient(c.client)}
	if c.basePath != "" {
		opts = append(opts, option.WithEndpoint(c.basePath))
	}
	return opts
}

func (c *apiClient) NewDriveService(ctx context.Context) (*drive.Service, error) {
	return drive.NewService(ctx, c.clientOptions()...)
}

func (c *apiClient) NewDirectoryService(ctx context.Context) (*directory.Service, error) {
	return directory.NewService(ctx, c.clientOptions()...)
}

func (c *apiClient) NewGroupsSettingsService(ctx context.Context) (*groupssettings.Service, error) {
	return groupssettings.NewService(ctx, c.clientOptions()...)
}

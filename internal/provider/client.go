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
}

func (c *apiClient) NewDriveService(ctx context.Context) (*drive.Service, error) {
	return drive.NewService(ctx, option.WithHTTPClient(c.client))
}

func (c *apiClient) NewDirectoryService(ctx context.Context) (*directory.Service, error) {
	return directory.NewService(ctx, option.WithHTTPClient(c.client))
}

func (c *apiClient) NewGroupsSettingsService(ctx context.Context) (*groupssettings.Service, error) {
	return groupssettings.NewService(ctx, option.WithHTTPClient(c.client))
}

# Authentication & Service Creation

## Creating a Service

```go
import (
    "context"
    "google.golang.org/api/drive/v3"
    "google.golang.org/api/option"
)

// Application Default Credentials (all scopes)
svc, err := drive.NewService(ctx)

// With explicit options
svc, err := drive.NewService(ctx,
    option.WithHTTPClient(httpClient),
    option.WithScopes(drive.DriveScope),
)

// With API key (limited; no user-specific data)
svc, err := drive.NewService(ctx, option.WithAPIKey("AIza..."))

// With OAuth token
svc, err := drive.NewService(ctx, option.WithTokenSource(tokenSource))

// With custom endpoint (testing)
svc, err := drive.NewService(ctx,
    option.WithHTTPClient(testServer.Client()),
    option.WithEndpoint(testServer.URL),
)
```

## This Provider's Pattern

The provider wraps service creation in `apiClient`:

```go
func (c *apiClient) NewDriveService(ctx context.Context) (*drive.Service, error) {
    return drive.NewService(ctx, c.clientOptions()...)
}

func (c *apiClient) clientOptions() []option.ClientOption {
    opts := []option.ClientOption{option.WithHTTPClient(c.client)}
    if c.basePath != "" {
        opts = append(opts, option.WithEndpoint(c.basePath))
    }
    return opts
}
```

In resource methods:

```go
driveSvc, err := r.client.NewDriveService(ctx)
if err != nil {
    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Drive service: %s", err))
    return
}
// Use driveSvc.Files, driveSvc.Drives, driveSvc.Permissions, etc.
```

## Service Structure

`drive.Service` exposes sub-services as fields:

```go
type Service struct {
    About       *AboutService
    Changes     *ChangesService
    Channels    *ChannelsService
    Comments    *CommentsService
    Drives      *DrivesService
    Files       *FilesService
    Permissions *PermissionsService
    Replies     *RepliesService
    Revisions   *RevisionsService
}
```

## Scopes

Choose the narrowest scope that fits:

| Scope                              | Access Level                          |
| ---------------------------------- | ------------------------------------- |
| `drive.DriveScope`                 | Full read/write/delete on all files   |
| `drive.DriveReadonlyScope`         | Read-only on all files                |
| `drive.DriveFileScope`             | Only files created/opened by this app |
| `drive.DriveMetadataScope`         | Metadata read/write, no content       |
| `drive.DriveMetadataReadonlyScope` | Metadata read-only                    |

This provider uses `DriveScope` (full access) since it manages drives and permissions as an admin.

Full scope details: `references/api/scopes.md`

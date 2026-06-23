# PermissionsService

```go
type PermissionsService struct{}

func NewPermissionsService(s *Service) *PermissionsService
func (r *PermissionsService) Create(fileId string, permission *Permission) *PermissionsCreateCall
func (r *PermissionsService) Delete(fileId string, permissionId string) *PermissionsDeleteCall
func (r *PermissionsService) Get(fileId string, permissionId string) *PermissionsGetCall
func (r *PermissionsService) List(fileId string) *PermissionsListCall
func (r *PermissionsService) Update(fileId string, permissionId string, permission *Permission) *PermissionsUpdateCall
```

## PermissionsCreateCall

```go
func (c *PermissionsCreateCall) Context(ctx context.Context) *PermissionsCreateCall
func (c *PermissionsCreateCall) Do(opts ...googleapi.CallOption) (*Permission, error)
func (c *PermissionsCreateCall) EmailMessage(v string) *PermissionsCreateCall
func (c *PermissionsCreateCall) EnforceSingleParent(v bool) *PermissionsCreateCall
func (c *PermissionsCreateCall) Fields(s ...googleapi.Field) *PermissionsCreateCall
func (c *PermissionsCreateCall) MoveToNewOwnersRoot(v bool) *PermissionsCreateCall
func (c *PermissionsCreateCall) SendNotificationEmail(v bool) *PermissionsCreateCall
func (c *PermissionsCreateCall) SupportsAllDrives(v bool) *PermissionsCreateCall
func (c *PermissionsCreateCall) TransferOwnership(v bool) *PermissionsCreateCall
func (c *PermissionsCreateCall) UseDomainAdminAccess(v bool) *PermissionsCreateCall
```

### Options

| Method                  | Description                                                  |
| ----------------------- | ------------------------------------------------------------ |
| `EmailMessage`          | Custom message in the notification email                     |
| `SendNotificationEmail` | Whether to send notification (default true for user/group)   |
| `TransferOwnership`     | Required true when setting role=owner                        |
| `MoveToNewOwnersRoot`   | Move file to new owner's My Drive root on ownership transfer |
| `UseDomainAdminAccess`  | Act as domain admin (requires admin privileges)              |
| `SupportsAllDrives`     | Required true for shared drive items                         |

## PermissionsGetCall

```go
func (c *PermissionsGetCall) Context(ctx context.Context) *PermissionsGetCall
func (c *PermissionsGetCall) Do(opts ...googleapi.CallOption) (*Permission, error)
func (c *PermissionsGetCall) Fields(s ...googleapi.Field) *PermissionsGetCall
func (c *PermissionsGetCall) IfNoneMatch(entityTag string) *PermissionsGetCall
func (c *PermissionsGetCall) SupportsAllDrives(v bool) *PermissionsGetCall
func (c *PermissionsGetCall) UseDomainAdminAccess(v bool) *PermissionsGetCall
```

## PermissionsListCall

```go
func (c *PermissionsListCall) Context(ctx context.Context) *PermissionsListCall
func (c *PermissionsListCall) Do(opts ...googleapi.CallOption) (*PermissionList, error)
func (c *PermissionsListCall) Fields(s ...googleapi.Field) *PermissionsListCall
func (c *PermissionsListCall) IfNoneMatch(entityTag string) *PermissionsListCall
func (c *PermissionsListCall) IncludePermissionsForView(v string) *PermissionsListCall
func (c *PermissionsListCall) PageSize(v int64) *PermissionsListCall
func (c *PermissionsListCall) PageToken(v string) *PermissionsListCall
func (c *PermissionsListCall) Pages(ctx context.Context, f func(*PermissionList) error) error
func (c *PermissionsListCall) SupportsAllDrives(v bool) *PermissionsListCall
func (c *PermissionsListCall) UseDomainAdminAccess(v bool) *PermissionsListCall
```

## PermissionsUpdateCall

```go
func (c *PermissionsUpdateCall) Context(ctx context.Context) *PermissionsUpdateCall
func (c *PermissionsUpdateCall) Do(opts ...googleapi.CallOption) (*Permission, error)
func (c *PermissionsUpdateCall) Fields(s ...googleapi.Field) *PermissionsUpdateCall
func (c *PermissionsUpdateCall) RemoveExpiration(v bool) *PermissionsUpdateCall
func (c *PermissionsUpdateCall) SupportsAllDrives(v bool) *PermissionsUpdateCall
func (c *PermissionsUpdateCall) TransferOwnership(v bool) *PermissionsUpdateCall
func (c *PermissionsUpdateCall) UseDomainAdminAccess(v bool) *PermissionsUpdateCall
```

## PermissionsDeleteCall

```go
func (c *PermissionsDeleteCall) Context(ctx context.Context) *PermissionsDeleteCall
func (c *PermissionsDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *PermissionsDeleteCall) Fields(s ...googleapi.Field) *PermissionsDeleteCall
func (c *PermissionsDeleteCall) SupportsAllDrives(v bool) *PermissionsDeleteCall
func (c *PermissionsDeleteCall) UseDomainAdminAccess(v bool) *PermissionsDeleteCall
```

## PermissionList

```go
type PermissionList struct {
    Kind            string        `json:"kind,omitempty"`
    NextPageToken   string        `json:"nextPageToken,omitempty"`
    Permissions     []*Permission `json:"permissions,omitempty"`
}
```

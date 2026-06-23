# DrivesService (Shared Drives)

## Drive Struct

Representation of a shared drive. The drive ID is also the ID of the top-level folder.

```go
type Drive struct {
    BackgroundImageFile *DriveBackgroundImageFile `json:"backgroundImageFile,omitempty"`
    BackgroundImageLink string                   `json:"backgroundImageLink,omitempty"`
    Capabilities        *DriveCapabilities       `json:"capabilities,omitempty"`
    ColorRgb            string                   `json:"colorRgb,omitempty"`
    CreatedTime         string                   `json:"createdTime,omitempty"`
    Hidden              bool                     `json:"hidden,omitempty"`
    Id                  string                   `json:"id,omitempty"`
    Kind                string                   `json:"kind,omitempty"`
    Name                string                   `json:"name,omitempty"`
    OrgUnitId           string                   `json:"orgUnitId,omitempty"`
    Restrictions        *DriveRestrictions       `json:"restrictions,omitempty"`
    ThemeId             string                   `json:"themeId,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## DriveRestrictions

```go
type DriveRestrictions struct {
    AdminManagedRestrictions                  bool `json:"adminManagedRestrictions,omitempty"`
    CopyRequiresWriterPermission              bool `json:"copyRequiresWriterPermission,omitempty"`
    DomainUsersOnly                           bool `json:"domainUsersOnly,omitempty"`
    DriveMembersOnly                          bool `json:"driveMembersOnly,omitempty"`
    SharingFoldersRequiresOrganizerPermission bool `json:"sharingFoldersRequiresOrganizerPermission,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

| Field                                       | Description                                             |
| ------------------------------------------- | ------------------------------------------------------- |
| `AdminManagedRestrictions`                  | Admin privileges required to modify restrictions        |
| `CopyRequiresWriterPermission`              | Disable copy/print/download for readers/commenters      |
| `DomainUsersOnly`                           | Restrict access to domain users                         |
| `DriveMembersOnly`                          | Restrict access to drive members                        |
| `SharingFoldersRequiresOrganizerPermission` | Only organizers can share folders (not file organizers) |

## DriveCapabilities (output only)

Key fields:

- `CanAddChildren`, `CanDeleteChildren`, `CanTrashChildren`
- `CanComment`, `CanCopy`, `CanDownload`, `CanEdit`
- `CanDeleteDrive`, `CanRenameDrive`
- `CanListChildren`, `CanReadRevisions`, `CanRename`, `CanShare`
- `CanManageMembers`
- `CanResetDriveRestrictions`
- `CanChangeCopyRequiresWriterPermissionRestriction`
- `CanChangeDomainUsersOnlyRestriction`
- `CanChangeDriveMembersOnlyRestriction`
- `CanChangeDriveBackground`
- `CanChangeSharingFoldersRequiresOrganizerPermissionRestriction`

## DrivesService

```go
type DrivesService struct{}

func NewDrivesService(s *Service) *DrivesService
func (r *DrivesService) Create(requestId string, drive *Drive) *DrivesCreateCall
func (r *DrivesService) Delete(driveId string) *DrivesDeleteCall
func (r *DrivesService) Get(driveId string) *DrivesGetCall
func (r *DrivesService) Hide(driveId string) *DrivesHideCall
func (r *DrivesService) List() *DrivesListCall
func (r *DrivesService) Unhide(driveId string) *DrivesUnhideCall
func (r *DrivesService) Update(driveId string, drive *Drive) *DrivesUpdateCall
```

**Important:** `Create` requires a `requestId` (UUID) for idempotency. Generate with `uuid.GenerateUUID()`.

## DrivesCreateCall

```go
func (c *DrivesCreateCall) Context(ctx context.Context) *DrivesCreateCall
func (c *DrivesCreateCall) Do(opts ...googleapi.CallOption) (*Drive, error)
func (c *DrivesCreateCall) Fields(s ...googleapi.Field) *DrivesCreateCall
```

## DrivesGetCall

```go
func (c *DrivesGetCall) Context(ctx context.Context) *DrivesGetCall
func (c *DrivesGetCall) Do(opts ...googleapi.CallOption) (*Drive, error)
func (c *DrivesGetCall) Fields(s ...googleapi.Field) *DrivesGetCall
func (c *DrivesGetCall) IfNoneMatch(entityTag string) *DrivesGetCall
func (c *DrivesGetCall) UseDomainAdminAccess(v bool) *DrivesGetCall
```

## DrivesListCall

```go
func (c *DrivesListCall) Context(ctx context.Context) *DrivesListCall
func (c *DrivesListCall) Do(opts ...googleapi.CallOption) (*DriveList, error)
func (c *DrivesListCall) Fields(s ...googleapi.Field) *DrivesListCall
func (c *DrivesListCall) IfNoneMatch(entityTag string) *DrivesListCall
func (c *DrivesListCall) PageSize(v int64) *DrivesListCall
func (c *DrivesListCall) PageToken(v string) *DrivesListCall
func (c *DrivesListCall) Pages(ctx context.Context, f func(*DriveList) error) error
func (c *DrivesListCall) Q(q string) *DrivesListCall
func (c *DrivesListCall) UseDomainAdminAccess(v bool) *DrivesListCall
```

## DrivesUpdateCall

```go
func (c *DrivesUpdateCall) Context(ctx context.Context) *DrivesUpdateCall
func (c *DrivesUpdateCall) Do(opts ...googleapi.CallOption) (*Drive, error)
func (c *DrivesUpdateCall) Fields(s ...googleapi.Field) *DrivesUpdateCall
func (c *DrivesUpdateCall) UseDomainAdminAccess(v bool) *DrivesUpdateCall
```

## DrivesDeleteCall

```go
func (c *DrivesDeleteCall) Context(ctx context.Context) *DrivesDeleteCall
func (c *DrivesDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *DrivesDeleteCall) Fields(s ...googleapi.Field) *DrivesDeleteCall
func (c *DrivesDeleteCall) UseDomainAdminAccess(v bool) *DrivesDeleteCall
```

## DriveList

```go
type DriveList struct {
    Drives        []*Drive `json:"drives,omitempty"`
    Kind          string   `json:"kind,omitempty"`
    NextPageToken string   `json:"nextPageToken,omitempty"`
}
```

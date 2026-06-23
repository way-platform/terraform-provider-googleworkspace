# ChangesService

## Change Struct

```go
type Change struct {
    ChangeType  string  `json:"changeType,omitempty"`
    Drive       *Drive  `json:"drive,omitempty"`
    DriveId     string  `json:"driveId,omitempty"`
    File        *File   `json:"file,omitempty"`
    FileId      string  `json:"fileId,omitempty"`
    Kind        string  `json:"kind,omitempty"`
    Removed     bool    `json:"removed,omitempty"`
    Time        string  `json:"time,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

| Field        | Description                                                          |
| ------------ | -------------------------------------------------------------------- |
| `ChangeType` | `"file"` or `"drive"`                                                |
| `File`       | Updated file state (present if type=file and not removed)            |
| `FileId`     | ID of the changed file                                               |
| `Drive`      | Updated drive state (present if type=drive and user is still member) |
| `DriveId`    | ID of the changed shared drive                                       |
| `Removed`    | Whether file/drive has been removed from this change list            |
| `Time`       | RFC 3339 timestamp of the change                                     |

## StartPageToken

```go
type StartPageToken struct {
    Kind           string `json:"kind,omitempty"`
    StartPageToken string `json:"startPageToken,omitempty"`
}
```

## ChangeList

```go
type ChangeList struct {
    Changes           []*Change `json:"changes,omitempty"`
    Kind              string    `json:"kind,omitempty"`
    NewStartPageToken string    `json:"newStartPageToken,omitempty"`
    NextPageToken     string    `json:"nextPageToken,omitempty"`
}
```

When `NewStartPageToken` is present, the current change set is complete. Use it as the page token for the next poll.

## ChangesService

```go
type ChangesService struct{}

func NewChangesService(s *Service) *ChangesService
func (r *ChangesService) GetStartPageToken() *ChangesGetStartPageTokenCall
func (r *ChangesService) List(pageToken string) *ChangesListCall
func (r *ChangesService) Watch(pageToken string, channel *Channel) *ChangesWatchCall
```

## ChangesGetStartPageTokenCall

```go
func (c *ChangesGetStartPageTokenCall) Context(ctx context.Context) *ChangesGetStartPageTokenCall
func (c *ChangesGetStartPageTokenCall) Do(opts ...googleapi.CallOption) (*StartPageToken, error)
func (c *ChangesGetStartPageTokenCall) DriveId(v string) *ChangesGetStartPageTokenCall
func (c *ChangesGetStartPageTokenCall) Fields(s ...googleapi.Field) *ChangesGetStartPageTokenCall
func (c *ChangesGetStartPageTokenCall) SupportsAllDrives(v bool) *ChangesGetStartPageTokenCall
```

## ChangesListCall

```go
func (c *ChangesListCall) Context(ctx context.Context) *ChangesListCall
func (c *ChangesListCall) Do(opts ...googleapi.CallOption) (*ChangeList, error)
func (c *ChangesListCall) DriveId(v string) *ChangesListCall
func (c *ChangesListCall) Fields(s ...googleapi.Field) *ChangesListCall
func (c *ChangesListCall) IfNoneMatch(entityTag string) *ChangesListCall
func (c *ChangesListCall) IncludeCorpusRemovals(v bool) *ChangesListCall
func (c *ChangesListCall) IncludeItemsFromAllDrives(v bool) *ChangesListCall
func (c *ChangesListCall) IncludeLabels(v string) *ChangesListCall
func (c *ChangesListCall) IncludePermissionsForView(v string) *ChangesListCall
func (c *ChangesListCall) IncludeRemoved(v bool) *ChangesListCall
func (c *ChangesListCall) PageSize(v int64) *ChangesListCall
func (c *ChangesListCall) RestrictToMyDrive(v bool) *ChangesListCall
func (c *ChangesListCall) Spaces(v string) *ChangesListCall
func (c *ChangesListCall) SupportsAllDrives(v bool) *ChangesListCall
```

## ChangesWatchCall

Subscribe to push notifications for changes.

```go
func (c *ChangesWatchCall) Context(ctx context.Context) *ChangesWatchCall
func (c *ChangesWatchCall) Do(opts ...googleapi.CallOption) (*Channel, error)
func (c *ChangesWatchCall) DriveId(v string) *ChangesWatchCall
func (c *ChangesWatchCall) Fields(s ...googleapi.Field) *ChangesWatchCall
func (c *ChangesWatchCall) IncludeCorpusRemovals(v bool) *ChangesWatchCall
func (c *ChangesWatchCall) IncludeItemsFromAllDrives(v bool) *ChangesWatchCall
func (c *ChangesWatchCall) IncludeLabels(v string) *ChangesWatchCall
func (c *ChangesWatchCall) IncludePermissionsForView(v string) *ChangesWatchCall
func (c *ChangesWatchCall) IncludeRemoved(v bool) *ChangesWatchCall
func (c *ChangesWatchCall) PageSize(v int64) *ChangesWatchCall
func (c *ChangesWatchCall) RestrictToMyDrive(v bool) *ChangesWatchCall
func (c *ChangesWatchCall) Spaces(v string) *ChangesWatchCall
func (c *ChangesWatchCall) SupportsAllDrives(v bool) *ChangesWatchCall
```

## Channel Struct

Used by Watch calls and ChannelsService.Stop.

```go
type Channel struct {
    Address    string            `json:"address,omitempty"`
    Expiration int64             `json:"expiration,omitempty,string"`
    Id         string            `json:"id,omitempty"`
    Kind       string            `json:"kind,omitempty"`
    Params     map[string]string `json:"params,omitempty"`
    Payload    bool              `json:"payload,omitempty"`
    ResourceId string            `json:"resourceId,omitempty"`
    ResourceUri string           `json:"resourceUri,omitempty"`
    Token      string            `json:"token,omitempty"`
    Type       string            `json:"type,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

| Field        | Description                                       |
| ------------ | ------------------------------------------------- |
| `Id`         | UUID identifying this channel                     |
| `Type`       | Delivery mechanism (typically `"web_hook"`)       |
| `Address`    | URL where notifications are delivered             |
| `Token`      | Arbitrary string delivered with each notification |
| `Expiration` | Unix timestamp in milliseconds                    |
| `ResourceId` | Opaque ID of watched resource (output)            |

## ChannelsService

```go
type ChannelsService struct{}

func NewChannelsService(s *Service) *ChannelsService
func (r *ChannelsService) Stop(channel *Channel) *ChannelsStopCall
```

Stop a push notification channel by providing the channel's `Id` and `ResourceId`.

# FilesService

```go
type FilesService struct{}

func NewFilesService(s *Service) *FilesService
func (r *FilesService) Copy(fileId string, file *File) *FilesCopyCall
func (r *FilesService) Create(file *File) *FilesCreateCall
func (r *FilesService) Delete(fileId string) *FilesDeleteCall
func (r *FilesService) EmptyTrash() *FilesEmptyTrashCall
func (r *FilesService) Export(fileId string, mimeType string) *FilesExportCall
func (r *FilesService) GenerateIds() *FilesGenerateIdsCall
func (r *FilesService) Get(fileId string) *FilesGetCall
func (r *FilesService) List() *FilesListCall
func (r *FilesService) ListLabels(fileId string) *FilesListLabelsCall
func (r *FilesService) ModifyLabels(fileId string, req *ModifyLabelsRequest) *FilesModifyLabelsCall
func (r *FilesService) Update(fileId string, file *File) *FilesUpdateCall
func (r *FilesService) Watch(fileId string, channel *Channel) *FilesWatchCall
```

## FilesCreateCall

```go
func (c *FilesCreateCall) Context(ctx context.Context) *FilesCreateCall
func (c *FilesCreateCall) Do(opts ...googleapi.CallOption) (*File, error)
func (c *FilesCreateCall) EnforceSingleParent(v bool) *FilesCreateCall
func (c *FilesCreateCall) Fields(s ...googleapi.Field) *FilesCreateCall
func (c *FilesCreateCall) IgnoreDefaultVisibility(v bool) *FilesCreateCall
func (c *FilesCreateCall) IncludeLabels(v string) *FilesCreateCall
func (c *FilesCreateCall) IncludePermissionsForView(v string) *FilesCreateCall
func (c *FilesCreateCall) KeepRevisionForever(v bool) *FilesCreateCall
func (c *FilesCreateCall) Media(r io.Reader, options ...googleapi.MediaOption) *FilesCreateCall
func (c *FilesCreateCall) OcrLanguage(v string) *FilesCreateCall
func (c *FilesCreateCall) ProgressUpdater(pu googleapi.ProgressUpdater) *FilesCreateCall
func (c *FilesCreateCall) ResumableMedia(ctx context.Context, r io.ReaderAt, size int64, mediaType string) *FilesCreateCall
func (c *FilesCreateCall) SupportsAllDrives(v bool) *FilesCreateCall
func (c *FilesCreateCall) UseContentAsIndexableText(v bool) *FilesCreateCall
```

## FilesGetCall

```go
func (c *FilesGetCall) AcknowledgeAbuse(v bool) *FilesGetCall
func (c *FilesGetCall) Context(ctx context.Context) *FilesGetCall
func (c *FilesGetCall) Do(opts ...googleapi.CallOption) (*File, error)
func (c *FilesGetCall) Download(opts ...googleapi.CallOption) (*http.Response, error)
func (c *FilesGetCall) Fields(s ...googleapi.Field) *FilesGetCall
func (c *FilesGetCall) IfNoneMatch(entityTag string) *FilesGetCall
func (c *FilesGetCall) IncludeLabels(v string) *FilesGetCall
func (c *FilesGetCall) IncludePermissionsForView(v string) *FilesGetCall
func (c *FilesGetCall) SupportsAllDrives(v bool) *FilesGetCall
```

## FilesListCall

```go
func (c *FilesListCall) Context(ctx context.Context) *FilesListCall
func (c *FilesListCall) Corpora(v string) *FilesListCall
func (c *FilesListCall) Do(opts ...googleapi.CallOption) (*FileList, error)
func (c *FilesListCall) DriveId(v string) *FilesListCall
func (c *FilesListCall) Fields(s ...googleapi.Field) *FilesListCall
func (c *FilesListCall) IfNoneMatch(entityTag string) *FilesListCall
func (c *FilesListCall) IncludeItemsFromAllDrives(v bool) *FilesListCall
func (c *FilesListCall) IncludeLabels(v string) *FilesListCall
func (c *FilesListCall) IncludePermissionsForView(v string) *FilesListCall
func (c *FilesListCall) OrderBy(v string) *FilesListCall
func (c *FilesListCall) PageSize(v int64) *FilesListCall
func (c *FilesListCall) PageToken(v string) *FilesListCall
func (c *FilesListCall) Pages(ctx context.Context, f func(*FileList) error) error
func (c *FilesListCall) Q(q string) *FilesListCall
func (c *FilesListCall) Spaces(v string) *FilesListCall
func (c *FilesListCall) SupportsAllDrives(v bool) *FilesListCall
```

## FilesUpdateCall

```go
func (c *FilesUpdateCall) AddParents(v string) *FilesUpdateCall
func (c *FilesUpdateCall) Context(ctx context.Context) *FilesUpdateCall
func (c *FilesUpdateCall) Do(opts ...googleapi.CallOption) (*File, error)
func (c *FilesUpdateCall) Fields(s ...googleapi.Field) *FilesUpdateCall
func (c *FilesUpdateCall) IncludeLabels(v string) *FilesUpdateCall
func (c *FilesUpdateCall) IncludePermissionsForView(v string) *FilesUpdateCall
func (c *FilesUpdateCall) KeepRevisionForever(v bool) *FilesUpdateCall
func (c *FilesUpdateCall) Media(r io.Reader, options ...googleapi.MediaOption) *FilesUpdateCall
func (c *FilesUpdateCall) OcrLanguage(v string) *FilesUpdateCall
func (c *FilesUpdateCall) ProgressUpdater(pu googleapi.ProgressUpdater) *FilesUpdateCall
func (c *FilesUpdateCall) RemoveParents(v string) *FilesUpdateCall
func (c *FilesUpdateCall) ResumableMedia(ctx context.Context, r io.ReaderAt, size int64, mediaType string) *FilesUpdateCall
func (c *FilesUpdateCall) SupportsAllDrives(v bool) *FilesUpdateCall
func (c *FilesUpdateCall) UseContentAsIndexableText(v bool) *FilesUpdateCall
```

## FilesCopyCall

```go
func (c *FilesCopyCall) Context(ctx context.Context) *FilesCopyCall
func (c *FilesCopyCall) Do(opts ...googleapi.CallOption) (*File, error)
func (c *FilesCopyCall) EnforceSingleParent(v bool) *FilesCopyCall
func (c *FilesCopyCall) Fields(s ...googleapi.Field) *FilesCopyCall
func (c *FilesCopyCall) IgnoreDefaultVisibility(v bool) *FilesCopyCall
func (c *FilesCopyCall) IncludeLabels(v string) *FilesCopyCall
func (c *FilesCopyCall) IncludePermissionsForView(v string) *FilesCopyCall
func (c *FilesCopyCall) KeepRevisionForever(v bool) *FilesCopyCall
func (c *FilesCopyCall) OcrLanguage(v string) *FilesCopyCall
func (c *FilesCopyCall) SupportsAllDrives(v bool) *FilesCopyCall
```

## FilesDeleteCall

```go
func (c *FilesDeleteCall) Context(ctx context.Context) *FilesDeleteCall
func (c *FilesDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *FilesDeleteCall) EnforceSingleParent(v bool) *FilesDeleteCall
func (c *FilesDeleteCall) Fields(s ...googleapi.Field) *FilesDeleteCall
func (c *FilesDeleteCall) SupportsAllDrives(v bool) *FilesDeleteCall
```

## FilesExportCall

Exports a Google Workspace document to a specified MIME type. Returns file content, not metadata.

```go
func (c *FilesExportCall) Context(ctx context.Context) *FilesExportCall
func (c *FilesExportCall) Do(opts ...googleapi.CallOption) error
func (c *FilesExportCall) Download(opts ...googleapi.CallOption) (*http.Response, error)
func (c *FilesExportCall) Fields(s ...googleapi.Field) *FilesExportCall
```

## FilesGenerateIdsCall

Pre-generate file IDs for use in create requests.

```go
func (c *FilesGenerateIdsCall) Context(ctx context.Context) *FilesGenerateIdsCall
func (c *FilesGenerateIdsCall) Count(v int64) *FilesGenerateIdsCall
func (c *FilesGenerateIdsCall) Do(opts ...googleapi.CallOption) (*GeneratedIds, error)
func (c *FilesGenerateIdsCall) Fields(s ...googleapi.Field) *FilesGenerateIdsCall
func (c *FilesGenerateIdsCall) Space(v string) *FilesGenerateIdsCall
func (c *FilesGenerateIdsCall) Type(v string) *FilesGenerateIdsCall
```

## FilesWatchCall

Subscribe to push notifications for changes to a file.

```go
func (c *FilesWatchCall) AcknowledgeAbuse(v bool) *FilesWatchCall
func (c *FilesWatchCall) Context(ctx context.Context) *FilesWatchCall
func (c *FilesWatchCall) Do(opts ...googleapi.CallOption) (*Channel, error)
func (c *FilesWatchCall) Fields(s ...googleapi.Field) *FilesWatchCall
func (c *FilesWatchCall) SupportsAllDrives(v bool) *FilesWatchCall
```

## FileList

```go
type FileList struct {
    Files             []*File `json:"files,omitempty"`
    IncompleteSearch  bool    `json:"incompleteSearch,omitempty"`
    Kind              string  `json:"kind,omitempty"`
    NextPageToken     string  `json:"nextPageToken,omitempty"`
}
```

## Query Syntax (Q parameter)

Used with `FilesListCall.Q()`:

| Operator   | Example                                            |
| ---------- | -------------------------------------------------- |
| `=`        | `name = 'hello.txt'`                               |
| `!=`       | `mimeType != 'application/vnd.google-apps.folder'` |
| `contains` | `name contains 'hello'`                            |
| `in`       | `'parent-id' in parents`                           |
| `and`      | `name = 'hello' and mimeType = 'text/plain'`       |
| `or`       | `name = 'hello' or name = 'world'`                 |
| `not`      | `not name contains 'temp'`                         |

Common queries:

```
// Files in a specific folder
'folder-id' in parents

// Non-trashed files in a folder
'folder-id' in parents and trashed = false

// Files by MIME type
mimeType = 'application/vnd.google-apps.folder'

// Shared drive files
driveId = 'drive-id' and trashed = false
```

## OrderBy Values

Comma-separated, each optionally followed by `desc`:

- `createdTime`
- `folder`
- `modifiedByMeTime`
- `modifiedTime`
- `name`
- `name_natural`
- `quotaBytesUsed`
- `recency`
- `sharedWithMeTime`
- `starred`
- `viewedByMeTime`

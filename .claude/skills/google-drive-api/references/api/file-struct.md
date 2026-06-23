# File Struct

The metadata for a file. Some resource methods (such as `files.update`) require a `fileId`. Use the `files.list` method to retrieve the ID for a file.

```go
type File struct {
    AppProperties                map[string]string      `json:"appProperties,omitempty"`
    Capabilities                 *FileCapabilities      `json:"capabilities,omitempty"`
    ContentHints                 *FileContentHints      `json:"contentHints,omitempty"`
    ContentRestrictions          []*ContentRestriction  `json:"contentRestrictions,omitempty"`
    CopyRequiresWriterPermission bool                   `json:"copyRequiresWriterPermission,omitempty"`
    CreatedTime                  string                 `json:"createdTime,omitempty"`
    Description                  string                 `json:"description,omitempty"`
    DriveId                      string                 `json:"driveId,omitempty"`
    ExplicitlyTrashed            bool                   `json:"explicitlyTrashed,omitempty"`
    ExportLinks                  map[string]string      `json:"exportLinks,omitempty"`
    FileExtension                string                 `json:"fileExtension,omitempty"`
    FolderColorRgb               string                 `json:"folderColorRgb,omitempty"`
    FullFileExtension            string                 `json:"fullFileExtension,omitempty"`
    HasAugmentedPermissions      bool                   `json:"hasAugmentedPermissions,omitempty"`
    HasThumbnail                 bool                   `json:"hasThumbnail,omitempty"`
    HeadRevisionId               string                 `json:"headRevisionId,omitempty"`
    IconLink                     string                 `json:"iconLink,omitempty"`
    Id                           string                 `json:"id,omitempty"`
    ImageMediaMetadata           *FileImageMediaMetadata `json:"imageMediaMetadata,omitempty"`
    IsAppAuthorized              bool                   `json:"isAppAuthorized,omitempty"`
    Kind                         string                 `json:"kind,omitempty"`
    LabelInfo                    *FileLabelInfo         `json:"labelInfo,omitempty"`
    LastModifyingUser            *User                  `json:"lastModifyingUser,omitempty"`
    LinkShareMetadata            *FileLinkShareMetadata `json:"linkShareMetadata,omitempty"`
    Md5Checksum                  string                 `json:"md5Checksum,omitempty"`
    MimeType                     string                 `json:"mimeType,omitempty"`
    ModifiedByMe                 bool                   `json:"modifiedByMe,omitempty"`
    ModifiedByMeTime             string                 `json:"modifiedByMeTime,omitempty"`
    ModifiedTime                 string                 `json:"modifiedTime,omitempty"`
    Name                         string                 `json:"name,omitempty"`
    OriginalFilename             string                 `json:"originalFilename,omitempty"`
    OwnedByMe                    bool                   `json:"ownedByMe,omitempty"`
    Owners                       []*User                `json:"owners,omitempty"`
    Parents                      []string               `json:"parents,omitempty"`
    PermissionIds                []string               `json:"permissionIds,omitempty"`
    Permissions                  []*Permission          `json:"permissions,omitempty"`
    Properties                   map[string]string      `json:"properties,omitempty"`
    QuotaBytesUsed               int64                  `json:"quotaBytesUsed,omitempty,string"`
    ResourceKey                  string                 `json:"resourceKey,omitempty"`
    Sha1Checksum                 string                 `json:"sha1Checksum,omitempty"`
    Sha256Checksum               string                 `json:"sha256Checksum,omitempty"`
    Shared                       bool                   `json:"shared,omitempty"`
    SharedWithMeTime             string                 `json:"sharedWithMeTime,omitempty"`
    SharingUser                  *User                  `json:"sharingUser,omitempty"`
    ShortcutDetails              *FileShortcutDetails   `json:"shortcutDetails,omitempty"`
    Size                         int64                  `json:"size,omitempty,string"`
    Spaces                       []string               `json:"spaces,omitempty"`
    Starred                      bool                   `json:"starred,omitempty"`
    ThumbnailLink                string                 `json:"thumbnailLink,omitempty"`
    ThumbnailVersion             int64                  `json:"thumbnailVersion,omitempty,string"`
    Trashed                      bool                   `json:"trashed,omitempty"`
    TrashedTime                  string                 `json:"trashedTime,omitempty"`
    TrashingUser                 *User                  `json:"trashingUser,omitempty"`
    Version                      int64                  `json:"version,omitempty,string"`
    VideoMediaMetadata           *FileVideoMediaMetadata `json:"videoMediaMetadata,omitempty"`
    ViewedByMe                   bool                   `json:"viewedByMe,omitempty"`
    ViewedByMeTime               string                 `json:"viewedByMeTime,omitempty"`
    WebContentLink               string                 `json:"webContentLink,omitempty"`
    WebViewLink                  string                 `json:"webViewLink,omitempty"`
    WritersCanShare              bool                   `json:"writersCanShare,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## Field Notes

| Field                          | Writable          | Notes                                                                                |
| ------------------------------ | ----------------- | ------------------------------------------------------------------------------------ |
| `Id`                           | No (output)       | Server-generated unique ID                                                           |
| `Name`                         | Yes               | Not necessarily unique within a folder                                               |
| `MimeType`                     | Yes (on create)   | Auto-detected if not set; cannot change after creation unless uploading new revision |
| `Parents`                      | Yes (create only) | On update, use `AddParents`/`RemoveParents` on the call                              |
| `Description`                  | Yes               | Short description                                                                    |
| `Properties`                   | Yes               | Visible to all apps                                                                  |
| `AppProperties`                | Yes               | Private to requesting app (requires authenticated request)                           |
| `Starred`                      | Yes               | User star                                                                            |
| `Trashed`                      | Yes               | Only owner can trash                                                                 |
| `CopyRequiresWriterPermission` | Yes               | Disable copy/print/download for readers                                              |
| `WritersCanShare`              | Yes               | Allow writers to modify permissions                                                  |
| `FolderColorRgb`               | Yes               | RGB hex string for folder color                                                      |
| `ContentHints`                 | Yes (write-only)  | Never populated in responses                                                         |
| `DriveId`                      | No (output)       | Only for shared drive items                                                          |
| `Size`                         | No (output)       | Bytes; not populated for folders/shortcuts                                           |
| `CreatedTime`                  | No (output)       | RFC 3339                                                                             |
| `ModifiedTime`                 | Yes               | Setting this also updates ModifiedByMeTime                                           |
| `Capabilities`                 | No (output)       | What the current user can do                                                         |
| `Permissions`                  | No (output)       | Only if user can share; use PermissionsService to manage                             |

## Google Workspace MIME Types

| MIME Type                                  | File Type     |
| ------------------------------------------ | ------------- |
| `application/vnd.google-apps.document`     | Google Docs   |
| `application/vnd.google-apps.spreadsheet`  | Google Sheets |
| `application/vnd.google-apps.presentation` | Google Slides |
| `application/vnd.google-apps.folder`       | Folder        |
| `application/vnd.google-apps.shortcut`     | Shortcut      |
| `application/vnd.google-apps.form`         | Google Forms  |

## FileCapabilities (output only)

Key boolean fields indicating what the current user can do:

- `CanEdit`, `CanComment`, `CanCopy`, `CanDownload`
- `CanDelete`, `CanTrash`, `CanUntrash`
- `CanRename`, `CanShare`
- `CanAddChildren`, `CanRemoveChildren` (folders)
- `CanMoveItemWithinDrive`, `CanMoveItemOutOfDrive`
- `CanModifyContent`, `CanModifyLabels`
- `CanReadRevisions`, `CanReadDrive`, `CanReadLabels`

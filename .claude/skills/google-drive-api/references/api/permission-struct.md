# Permission Struct

A permission grants a user, group, domain, or the world access to a file or a folder hierarchy. Some resource methods (such as `permissions.update`) require a `permissionId`. Use the `permissions.list` method to retrieve the ID.

```go
type Permission struct {
    AllowFileDiscovery bool     `json:"allowFileDiscovery,omitempty"`
    Deleted            bool     `json:"deleted,omitempty"`
    DisplayName        string   `json:"displayName,omitempty"`
    Domain             string   `json:"domain,omitempty"`
    EmailAddress       string   `json:"emailAddress,omitempty"`
    ExpirationTime     string   `json:"expirationTime,omitempty"`
    Id                 string   `json:"id,omitempty"`
    Kind               string   `json:"kind,omitempty"`
    PendingOwner       bool     `json:"pendingOwner,omitempty"`
    PhotoLink          string   `json:"photoLink,omitempty"`
    Role               string   `json:"role,omitempty"`
    Type               string   `json:"type,omitempty"`
    View               string   `json:"view,omitempty"`

    PermissionDetails          []*PermissionPermissionDetails          `json:"permissionDetails,omitempty"`
    TeamDrivePermissionDetails []*PermissionTeamDrivePermissionDetails `json:"teamDrivePermissionDetails,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## Field Details

### Type (required on create)

| Value    | Description     | Required Fields |
| -------- | --------------- | --------------- |
| `user`   | Individual user | `emailAddress`  |
| `group`  | Google Group    | `emailAddress`  |
| `domain` | Entire domain   | `domain`        |
| `anyone` | Public access   | (none)          |

### Role (required)

| Value           | Description                                               |
| --------------- | --------------------------------------------------------- |
| `owner`         | Full ownership (transfer with `TransferOwnership` option) |
| `organizer`     | Shared drive: manage members + all content                |
| `fileOrganizer` | Shared drive: manage content but not members              |
| `writer`        | Edit files                                                |
| `commenter`     | View + comment                                            |
| `reader`        | View only                                                 |

### Other Fields

| Field                | Notes                                                                         |
| -------------------- | ----------------------------------------------------------------------------- |
| `AllowFileDiscovery` | Only for `domain` or `anyone` types; whether file appears in search           |
| `ExpirationTime`     | RFC 3339; only for `user` and `group`; max 1 year in the future               |
| `PendingOwner`       | Whether user is pending owner; only for `user` type on non-shared-drive files |
| `DisplayName`        | Output only; human-readable name of grantee                                   |
| `Deleted`            | Output only; whether the account has been deleted                             |
| `View`               | Only value: `"published"`                                                     |

## PermissionPermissionDetails

Output only, for shared drive items. Shows whether permissions are inherited or direct.

```go
type PermissionPermissionDetails struct {
    Inherited           bool   `json:"inherited,omitempty"`
    InheritedFrom       string `json:"inheritedFrom,omitempty"`
    PermissionType      string `json:"permissionType,omitempty"`
    Role                string `json:"role,omitempty"`
    ForceSendFields     []string `json:"-"`
    NullFields          []string `json:"-"`
}
```

# Group Struct

```go
type Group struct {
    // --- Writable fields ---

    Email       string `json:"email,omitempty"`       // Required for creation; must be unique
    Name        string `json:"name,omitempty"`        // Display name
    Description string `json:"description,omitempty"` // Max 4,096 characters

    // --- Read-only fields ---

    Id                 string   `json:"id,omitempty"`                 // Unique group ID (usable as groupKey)
    AdminCreated       bool     `json:"adminCreated,omitempty"`       // true if created by admin
    Aliases            []string `json:"aliases,omitempty"`
    NonEditableAliases []string `json:"nonEditableAliases,omitempty"` // Outside primary domain
    DirectMembersCount int64    `json:"directMembersCount,omitempty"`
    Etag               string   `json:"etag,omitempty"`
    Kind               string   `json:"kind,omitempty"`

    // --- JSON control ---
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## Groups (List Response)

```go
type Groups struct {
    Groups        []*Group `json:"groups,omitempty"`
    NextPageToken string   `json:"nextPageToken,omitempty"`
    Etag          string   `json:"etag,omitempty"`
    Kind          string   `json:"kind,omitempty"`
}
```

## Field Behavior Notes

| Field                | Behavior                                                     |
| -------------------- | ------------------------------------------------------------ |
| `Email`              | Required on Insert; used as groupKey                         |
| `Name`               | Display name shown in admin console                          |
| `Description`        | Optional; max 4096 chars                                     |
| `Id`                 | Stable unique ID; use as `groupKey` for permanent references |
| `DirectMembersCount` | Only counts direct members, not nested group members         |
| `AdminCreated`       | Distinguishes admin-created from user-created groups         |

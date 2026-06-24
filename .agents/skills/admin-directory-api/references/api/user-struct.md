# User Struct

```go
type User struct {
    // --- Writable fields ---

    PrimaryEmail              string      `json:"primaryEmail,omitempty"`   // Required for creation; must be unique
    Name                      *UserName   `json:"name,omitempty"`           // Required for creation
    Password                  string      `json:"password,omitempty"`       // Plaintext or hashed (see HashFunction)
    HashFunction              string      `json:"hashFunction,omitempty"`   // "MD5", "SHA-1", or "crypt"
    OrgUnitPath               string      `json:"orgUnitPath,omitempty"`    // Default "/"
    Suspended                 bool        `json:"suspended,omitempty"`
    Archived                  bool        `json:"archived,omitempty"`
    ChangePasswordAtNextLogin bool        `json:"changePasswordAtNextLogin,omitempty"`
    IncludeInGlobalAddressList bool       `json:"includeInGlobalAddressList,omitempty"`
    IpWhitelisted             bool        `json:"ipWhitelisted,omitempty"`  // Deprecated allowlist
    RecoveryEmail             string      `json:"recoveryEmail,omitempty"`
    RecoveryPhone             string      `json:"recoveryPhone,omitempty"`  // E.164 format: +16506661212

    // Structured data fields (interface{} — pass []map[string]interface{})
    Addresses     interface{} `json:"addresses,omitempty"`     // Max 10KB
    Emails        interface{} `json:"emails,omitempty"`        // Max 10KB
    Phones        interface{} `json:"phones,omitempty"`        // Max 1KB
    Organizations interface{} `json:"organizations,omitempty"` // Max 10KB
    Relations     interface{} `json:"relations,omitempty"`     // Max 2KB
    ExternalIds   interface{} `json:"externalIds,omitempty"`   // Max 2KB
    Websites      interface{} `json:"websites,omitempty"`      // Max 2KB
    Ims           interface{} `json:"ims,omitempty"`           // Max 2KB
    Keywords      interface{} `json:"keywords,omitempty"`      // Max 1KB
    Languages     interface{} `json:"languages,omitempty"`     // Max 1KB
    Locations     interface{} `json:"locations,omitempty"`     // Max 10KB
    Gender        interface{} `json:"gender,omitempty"`        // Max 1KB
    Notes         interface{} `json:"notes,omitempty"`
    PosixAccounts interface{} `json:"posixAccounts,omitempty"`
    SshPublicKeys interface{} `json:"sshPublicKeys,omitempty"`

    // Custom schemas (key = schema name, value = field map)
    CustomSchemas map[string]googleapi.RawMessage `json:"customSchemas,omitempty"`

    // --- Read-only fields ---

    Id                 string   `json:"id,omitempty"`                 // Unique user ID (usable as userKey)
    CustomerId         string   `json:"customerId,omitempty"`
    CreationTime       string   `json:"creationTime,omitempty"`
    LastLoginTime      string   `json:"lastLoginTime,omitempty"`
    DeletionTime       string   `json:"deletionTime,omitempty"`
    Aliases            []string `json:"aliases,omitempty"`
    NonEditableAliases []string `json:"nonEditableAliases,omitempty"` // Outside primary domain
    IsAdmin            bool     `json:"isAdmin,omitempty"`            // Use MakeAdmin to set
    IsDelegatedAdmin   bool     `json:"isDelegatedAdmin,omitempty"`
    IsEnforcedIn2Sv    bool     `json:"isEnforcedIn2Sv,omitempty"`
    IsEnrolledIn2Sv    bool     `json:"isEnrolledIn2Sv,omitempty"`
    IsMailboxSetup     bool     `json:"isMailboxSetup,omitempty"`
    AgreedToTerms      bool     `json:"agreedToTerms,omitempty"`
    SuspensionReason   string   `json:"suspensionReason,omitempty"`
    ThumbnailPhotoUrl  string   `json:"thumbnailPhotoUrl,omitempty"`
    ThumbnailPhotoEtag string   `json:"thumbnailPhotoEtag,omitempty"`
    Etag               string   `json:"etag,omitempty"`
    Kind               string   `json:"kind,omitempty"`

    // --- JSON control ---
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## UserName Struct

```go
type UserName struct {
    GivenName   string `json:"givenName,omitempty"`   // First name; required for creation
    FamilyName  string `json:"familyName,omitempty"`  // Last name; required for creation
    FullName    string `json:"fullName,omitempty"`    // Read-only; concatenated
    DisplayName string `json:"displayName,omitempty"` // Max 256 characters

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## Users (List Response)

```go
type Users struct {
    Users         []*User `json:"users,omitempty"`
    NextPageToken string  `json:"nextPageToken,omitempty"`
    Etag          string  `json:"etag,omitempty"`
    Kind          string  `json:"kind,omitempty"`
}
```

## UserMakeAdmin

```go
type UserMakeAdmin struct {
    Status bool `json:"status,omitempty"` // true = make admin, false = revoke
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## UserUndelete

```go
type UserUndelete struct {
    OrgUnitPath     string   `json:"orgUnitPath,omitempty"` // Where to restore user
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## UserPhoto

```go
type UserPhoto struct {
    Etag      string `json:"etag,omitempty"`
    Height    int64  `json:"height,omitempty"`
    Id        string `json:"id,omitempty"`
    Kind      string `json:"kind,omitempty"`
    MimeType  string `json:"mimeType,omitempty"`  // JPEG, PNG, GIF, BMP, TIFF
    PhotoData string `json:"photoData,omitempty"` // Base64 URL-safe encoded
    Width     int64  `json:"width,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## Alias

```go
type Alias struct {
    Alias        string `json:"alias,omitempty"`        // The alias email
    Etag         string `json:"etag,omitempty"`
    Id           string `json:"id,omitempty"`
    Kind         string `json:"kind,omitempty"`
    PrimaryEmail string `json:"primaryEmail,omitempty"` // Owner's primary email

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## Field Behavior Notes

| Field              | Behavior                                                          |
| ------------------ | ----------------------------------------------------------------- |
| `PrimaryEmail`     | Required on Insert; used as userKey                               |
| `Name.GivenName`   | Required on Insert                                                |
| `Name.FamilyName`  | Required on Insert                                                |
| `Password`         | Required on Insert; not returned on Get                           |
| `IsAdmin`          | Read-only; use `MakeAdmin()` to change                            |
| `IsDelegatedAdmin` | Read-only; set via Admin Console                                  |
| `CustomSchemas`    | Only returned with `Projection("full")` or `Projection("custom")` |
| `Id`               | Stable unique ID; use as `userKey` for permanent references       |

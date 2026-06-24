# Users

## Creating a User

```go
user := &admin.User{
    PrimaryEmail: "newuser@example.com",
    Name: &admin.UserName{
        GivenName:  "Jane",
        FamilyName: "Doe",
    },
    Password:     "initialPassword123",
    OrgUnitPath:  "/Engineering",
    ChangePasswordAtNextLogin: true,
    ForceSendFields: []string{"ChangePasswordAtNextLogin"},
}
created, err := svc.Users.Insert(user).Context(ctx).Do()
```

Required fields: `PrimaryEmail`, `Name.GivenName`, `Name.FamilyName`, `Password`.

## Reading a User

```go
user, err := svc.Users.Get("user@example.com").
    Projection("full").
    Context(ctx).
    Do()
```

Use `Projection("full")` to get all fields including custom schemas. The `userKey` parameter accepts primary email, alias email, or unique user ID.

## Updating a User

```go
user := &admin.User{
    Suspended:       false,
    OrgUnitPath:     "/Sales",
    ForceSendFields: []string{"Suspended"},
}
updated, err := svc.Users.Update("user@example.com", user).Context(ctx).Do()
```

`Update` replaces the entire resource. Use `Patch` for partial updates (only sends provided fields):

```go
updated, err := svc.Users.Patch("user@example.com", user).Context(ctx).Do()
```

## Deleting a User

```go
err := svc.Users.Delete("user@example.com").Context(ctx).Do()
```

Deleted users can be recovered within 20 days using `Undelete`.

## Listing Users

```go
var allUsers []*admin.User
err := svc.Users.List().
    Customer("my_customer").
    Projection("full").
    MaxResults(500).
    Pages(ctx, func(resp *admin.Users) error {
        allUsers = append(allUsers, resp.Users...)
        return nil
    })
```

Either `Customer` or `Domain` must be provided. Use `Query` for filtering:

```go
svc.Users.List().
    Customer("my_customer").
    Query("orgUnitPath='/Engineering'").
    Do()
```

## Custom Schemas on Users

Custom schema data lives in `User.CustomSchemas` as `map[string]googleapi.RawMessage`:

```go
user, err := svc.Users.Get(userKey).Projection("full").Do()
// user.CustomSchemas["MySchema"] contains raw JSON

// To set custom schema values:
import "encoding/json"

schemaData, _ := json.Marshal(map[string]interface{}{
    "field1": "value1",
    "field2": 42,
})
user := &admin.User{
    CustomSchemas: map[string]googleapi.RawMessage{
        "MySchema": schemaData,
    },
}
svc.Users.Patch(userKey, user).Context(ctx).Do()
```

Use `CustomFieldMask` with `Projection("custom")` to read specific schemas:

```go
svc.Users.Get(userKey).
    Projection("custom").
    CustomFieldMask("MySchema,OtherSchema").
    Do()
```

## User Aliases

```go
// Add alias
alias := &admin.Alias{Alias: "newalias@example.com"}
_, err := svc.Users.Aliases.Insert("user@example.com", alias).Context(ctx).Do()

// List aliases
aliases, err := svc.Users.Aliases.List("user@example.com").Context(ctx).Do()

// Delete alias
err := svc.Users.Aliases.Delete("user@example.com", "newalias@example.com").Context(ctx).Do()
```

## User Photos

```go
// Get photo
photo, err := svc.Users.Photos.Get("user@example.com").Context(ctx).Do()
// photo.PhotoData is base64 URL-safe encoded

// Update photo
photo := &admin.UserPhoto{
    PhotoData: base64URLEncodedData,
    MimeType:  "image/jpeg",
}
_, err := svc.Users.Photos.Update("user@example.com", photo).Context(ctx).Do()

// Delete photo
err := svc.Users.Photos.Delete("user@example.com").Context(ctx).Do()
```

## Admin Operations

```go
// Make user a super admin
err := svc.Users.MakeAdmin("user@example.com", &admin.UserMakeAdmin{
    Status: true,
}).Context(ctx).Do()

// Sign out user (invalidate sessions)
err := svc.Users.SignOut("user@example.com").Context(ctx).Do()

// Undelete a recently deleted user (within 20 days)
err := svc.Users.Undelete(userId, &admin.UserUndelete{
    OrgUnitPath: "/",
}).Context(ctx).Do()
```

Note: `Undelete` requires the user's unique ID (not email, which is released on deletion).

## ViewType

Controls which fields are visible:

| Value             | Behavior                                   |
| ----------------- | ------------------------------------------ |
| `"admin_view"`    | Default. All fields visible to admin.      |
| `"domain_public"` | Only fields visible to other domain users. |

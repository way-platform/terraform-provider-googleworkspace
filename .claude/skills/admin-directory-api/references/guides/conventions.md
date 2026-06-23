# Go SDK Conventions

## Call Chain Pattern

Every API operation follows the same builder pattern:

```go
result, err := svc.Users.Get(userKey).
    Projection("full").
    Fields("primaryEmail,name,suspended").
    Context(ctx).
    Do()
```

- Required arguments go in `Method()`
- Optional parameters are chained as method calls
- `Fields()` controls which response fields to include (partial response)
- `Context()` sets the request context
- `Do()` executes the request and returns the result

## ForceSendFields

Google API structs use `omitempty` on all JSON fields. Booleans set to `false`, integers set to `0`, and empty strings are omitted from the request body. To explicitly send these values, list the Go field names in `ForceSendFields`:

```go
user := &admin.User{
    Suspended:       false,
    Archived:        false,
    IncludeInGlobalAddressList: true,
    ForceSendFields: []string{
        "Suspended",
        "Archived",
        "IncludeInGlobalAddressList",
    },
}
```

Rule: always include boolean fields in `ForceSendFields` when you need to set them to `false`. Without this, `false` values are silently dropped from the request.

## NullFields

To explicitly send a JSON `null` (clearing a field), use `NullFields`:

```go
user := &admin.User{
    RecoveryEmail: "",
    NullFields:    []string{"RecoveryEmail"},
}
```

`ForceSendFields` sends the zero value; `NullFields` sends `null`. They are mutually exclusive per field.

## Fields() for Partial Responses

Reduce response size by requesting only needed fields:

```go
user, err := svc.Users.Get(userKey).Fields("primaryEmail,name,suspended").Do()

users, err := svc.Users.List().
    Customer("my_customer").
    Fields("users(primaryEmail,name),nextPageToken").
    Do()
```

Always include `nextPageToken` in list call field masks if paginating.

## Projection (Users only)

User responses have three projection levels:

| Value      | Returns                                               |
| ---------- | ----------------------------------------------------- |
| `"basic"`  | Default. Core fields only (no custom schemas)         |
| `"custom"` | Core fields + custom schema fields in CustomFieldMask |
| `"full"`   | All fields including all custom schemas               |

```go
svc.Users.Get(userKey).Projection("full").Do()
svc.Users.List().Customer("my_customer").Projection("full").Do()
```

## Pagination

### Manual (PageToken)

```go
var allMembers []*admin.Member
pageToken := ""
for {
    call := svc.Members.List(groupKey).MaxResults(200)
    if pageToken != "" {
        call = call.PageToken(pageToken)
    }
    result, err := call.Do()
    if err != nil {
        return err
    }
    allMembers = append(allMembers, result.Members...)
    if result.NextPageToken == "" {
        break
    }
    pageToken = result.NextPageToken
}
```

### Automatic (Pages)

```go
var allMembers []*admin.Member
err := svc.Members.List(groupKey).MaxResults(200).
    Pages(ctx, func(resp *admin.Members) error {
        allMembers = append(allMembers, resp.Members...)
        return nil
    })
```

`Pages()` handles the token loop automatically. Return an error from the callback to stop early.

## Error Handling

```go
import "google.golang.org/api/googleapi"

user, err := svc.Users.Get(userKey).Do()
if err != nil {
    if gerr, ok := err.(*googleapi.Error); ok {
        switch gerr.Code {
        case 404:
            // Resource not found
        case 403:
            // Permission denied
        case 409:
            // Entity already exists (Insert conflict)
        case 429:
            // Rate limited
        }
    }
    return err
}
```

In this provider, 404 on Read means the resource was deleted externally:

```go
if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
    resp.State.RemoveResource(ctx)
    return
}
```

## customerId Parameter

Many services require `customerId` as the first argument:

- OrgunitsService: `svc.Orgunits.List(customerId)`
- RolesService: `svc.Roles.List(customer)`
- RoleAssignmentsService: `svc.RoleAssignments.List(customer)`
- SchemasService: `svc.Schemas.List(customerId)`
- DomainsService: `svc.Domains.List(customer)`
- DomainAliasesService: `svc.DomainAliases.List(customer)`
- ChromeosdevicesService: `svc.Chromeosdevices.List(customerId)`
- MobiledevicesService: `svc.Mobiledevices.List(customerId)`

The string `"my_customer"` is a valid alias for the authenticated account's customer ID.

## userKey / groupKey

Any parameter named `userKey` accepts:

- User's primary email address (e.g., `user@example.com`)
- An alias email address
- The unique user ID string (from `User.Id`)

Any parameter named `groupKey` accepts:

- Group's email address
- An alias email address
- The unique group ID string (from `Group.Id`)

## Context

Always pass context for cancellation/timeout support:

```go
svc.Users.Get(userKey).Context(ctx).Do()
```

For `Pages()`, context is the first argument.

## Insert vs Create

The Admin Directory API uses `Insert` for creation operations (unlike Drive API which uses `Create`):

```go
svc.Users.Insert(user).Do()     // not Create
svc.Groups.Insert(group).Do()   // not Create
svc.Members.Insert(groupKey, member).Do()
```

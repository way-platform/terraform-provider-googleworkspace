# Go SDK Conventions

## Call Chain Pattern

Every API operation follows the same builder pattern:

```go
result, err := svc.Service.Method(requiredArgs...).
    Option1(value).
    Option2(value).
    Fields("field1,field2").
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
req := &drive.DriveRestrictions{
    AdminManagedRestrictions: false,
    DomainUsersOnly:          true,
    DriveMembersOnly:         false,
    ForceSendFields: []string{
        "AdminManagedRestrictions",
        "DomainUsersOnly",
        "DriveMembersOnly",
    },
}
```

Rule: always include boolean fields in `ForceSendFields` when you need to set them to `false`. Without this, `false` values are silently dropped from the request.

## NullFields

To explicitly send a JSON `null` (clearing a field), use `NullFields`:

```go
req := &drive.File{
    Description: "",
    NullFields:  []string{"Description"},
}
```

`ForceSendFields` sends the zero value; `NullFields` sends `null`. They are mutually exclusive per field.

## Fields() for Partial Responses

Reduce response size by requesting only needed fields:

```go
// Only get id and name
file, err := svc.Files.Get(fileId).Fields("id,name").Do()

// Nested fields use dot notation
drive, err := svc.Drives.Get(driveId).Fields("id,name,restrictions").Do()

// List calls: fields on the wrapper AND nested items
files, err := svc.Files.List().Fields("files(id,name,mimeType),nextPageToken").Do()
```

Always include `nextPageToken` in list call field masks if paginating.

## Pagination

### Manual (PageToken)

```go
var allFiles []*drive.File
pageToken := ""
for {
    call := svc.Files.List().PageSize(100).Q(query)
    if pageToken != "" {
        call = call.PageToken(pageToken)
    }
    result, err := call.Do()
    if err != nil {
        return err
    }
    allFiles = append(allFiles, result.Files...)
    if result.NextPageToken == "" {
        break
    }
    pageToken = result.NextPageToken
}
```

### Automatic (Pages)

```go
var allFiles []*drive.File
err := svc.Files.List().PageSize(100).Q(query).
    Pages(ctx, func(list *drive.FileList) error {
        allFiles = append(allFiles, list.Files...)
        return nil
    })
```

`Pages()` handles the token loop automatically. Return an error from the callback to stop early.

## Error Handling

```go
import "google.golang.org/api/googleapi"

result, err := svc.Files.Get(fileId).Do()
if err != nil {
    if gerr, ok := err.(*googleapi.Error); ok {
        switch gerr.Code {
        case 404:
            // Resource not found
        case 403:
            // Permission denied
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

## SupportsAllDrives

Any operation that might touch shared drive items needs `.SupportsAllDrives(true)`. Without it, shared drive files/permissions return 404 or are silently excluded from list results.

```go
// Required for shared drive items
svc.Permissions.Create(fileId, perm).SupportsAllDrives(true).Do()
svc.Permissions.Get(fileId, permId).SupportsAllDrives(true).Do()
svc.Files.Get(fileId).SupportsAllDrives(true).Do()
svc.Files.List().IncludeItemsFromAllDrives(true).SupportsAllDrives(true).Do()
```

## UseDomainAdminAccess

Allows a domain administrator to act on shared drives they don't directly have access to:

```go
svc.Drives.Get(driveId).UseDomainAdminAccess(true).Do()
svc.Drives.Update(driveId, req).UseDomainAdminAccess(true).Do()
svc.Permissions.Create(fileId, perm).UseDomainAdminAccess(true).Do()
```

## Context

Always pass context for cancellation/timeout support:

```go
svc.Files.Get(fileId).Context(ctx).Do()
```

Most calls accept `Context()` in the chain. For `Pages()`, context is the first argument.

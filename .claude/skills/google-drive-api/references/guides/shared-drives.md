# Shared Drives

The `DrivesService` manages shared drives (formerly Team Drives). A shared drive's ID is also the ID of its top-level folder.

## Creating a Shared Drive

`Create` requires a `requestId` (UUID) for idempotency:

```go
import "github.com/hashicorp/go-uuid"

requestId, err := uuid.GenerateUUID()
if err != nil {
    return err
}

driveReq := &drive.Drive{
    Name: "Engineering",
}
created, err := svc.Drives.Create(requestId, driveReq).Do()
```

Restrictions cannot be set on create; update immediately after:

```go
updateReq := &drive.Drive{
    Restrictions: &drive.DriveRestrictions{
        DriveMembersOnly: true,
        ForceSendFields:  []string{"DriveMembersOnly"},
    },
}
_, err = svc.Drives.Update(created.Id, updateReq).
    UseDomainAdminAccess(true).
    Fields("id,name,restrictions").
    Do()
```

## Reading a Shared Drive

```go
d, err := svc.Drives.Get(driveId).
    UseDomainAdminAccess(true).
    Fields("id,name,restrictions").
    Do()
```

## Listing Shared Drives

```go
var all []*drive.Drive
err := svc.Drives.List().
    PageSize(100).
    UseDomainAdminAccess(true).
    Pages(ctx, func(list *drive.DriveList) error {
        all = append(all, list.Drives...)
        return nil
    })
```

Filter with `Q()`:

```go
svc.Drives.List().Q("name contains 'Engineering'").Do()
```

## Updating a Shared Drive

```go
req := &drive.Drive{
    Name: "New Name",
    Restrictions: &drive.DriveRestrictions{
        AdminManagedRestrictions:                  true,
        CopyRequiresWriterPermission:              false,
        DomainUsersOnly:                           false,
        DriveMembersOnly:                          true,
        SharingFoldersRequiresOrganizerPermission: false,
        ForceSendFields: []string{
            "AdminManagedRestrictions",
            "CopyRequiresWriterPermission",
            "DomainUsersOnly",
            "DriveMembersOnly",
            "SharingFoldersRequiresOrganizerPermission",
        },
    },
}
updated, err := svc.Drives.Update(driveId, req).
    UseDomainAdminAccess(true).
    Fields("id,name,restrictions").
    Do()
```

## Deleting a Shared Drive

The drive must be empty (no untrashed items):

```go
err := svc.Drives.Delete(driveId).UseDomainAdminAccess(true).Do()
```

## Hide/Unhide

Hide a shared drive from the user's default view:

```go
svc.Drives.Hide(driveId).Do()
svc.Drives.Unhide(driveId).Do()
```

## Restrictions

| Field                                       | Effect                                                  |
| ------------------------------------------- | ------------------------------------------------------- |
| `AdminManagedRestrictions`                  | Only admins can modify these restrictions               |
| `CopyRequiresWriterPermission`              | Readers/commenters cannot copy/print/download           |
| `DomainUsersOnly`                           | Only domain users can access                            |
| `DriveMembersOnly`                          | Only drive members can access items                     |
| `SharingFoldersRequiresOrganizerPermission` | Only organizers (not file organizers) can share folders |

All are booleans defaulting to false. Always use `ForceSendFields` when setting restrictions to ensure `false` values are sent.

## This Provider's Pattern

The `resource_drive` creates the drive first, then updates restrictions in a separate call (the API doesn't support setting restrictions on create):

```go
created, err := driveSvc.Drives.Create(requestId, driveReq).Do()
// Then:
_, err = driveSvc.Drives.Update(created.Id, updateReq).
    UseDomainAdminAccess(plan.UseDomainAdminAccess.ValueBool()).
    Fields("id,name,restrictions").
    Do()
```

## Full API Reference

- Drive struct, DriveRestrictions, DriveCapabilities: `references/api/drives-service.md`

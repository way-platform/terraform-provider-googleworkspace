# Permissions

The `PermissionsService` manages sharing: who can access a file/folder/drive and what they can do.

## Creating a Permission

```go
perm := &drive.Permission{
    Role:         "writer",
    Type:         "user",
    EmailAddress: "user@example.com",
}

created, err := svc.Permissions.Create(fileId, perm).
    SupportsAllDrives(true).
    SendNotificationEmail(false).
    Fields("id,emailAddress,role,type").
    Do()
```

### Required fields for Create

| Type     | Required        |
| -------- | --------------- |
| `user`   | `EmailAddress`  |
| `group`  | `EmailAddress`  |
| `domain` | `Domain`        |
| `anyone` | (nothing extra) |

### Options on Create

| Method                        | Default               | Notes                           |
| ----------------------------- | --------------------- | ------------------------------- |
| `SendNotificationEmail(bool)` | true (for user/group) | Send email to grantee           |
| `EmailMessage(string)`        | —                     | Custom message in notification  |
| `TransferOwnership(bool)`     | false                 | Required true for role=owner    |
| `UseDomainAdminAccess(bool)`  | false                 | Act as domain admin             |
| `SupportsAllDrives(bool)`     | false                 | Required for shared drive items |

## Reading a Permission

```go
perm, err := svc.Permissions.Get(fileId, permissionId).
    SupportsAllDrives(true).
    UseDomainAdminAccess(true).
    Fields("id,emailAddress,role,type").
    Do()
```

## Listing Permissions

```go
list, err := svc.Permissions.List(fileId).
    SupportsAllDrives(true).
    UseDomainAdminAccess(true).
    Fields("permissions(id,emailAddress,role,type)").
    PageSize(100).
    Do()
```

Supports pagination with `Pages()`:

```go
var all []*drive.Permission
err := svc.Permissions.List(fileId).
    SupportsAllDrives(true).
    Pages(ctx, func(list *drive.PermissionList) error {
        all = append(all, list.Permissions...)
        return nil
    })
```

## Updating a Permission

Only `Role` can be changed (type/email are immutable):

```go
perm := &drive.Permission{
    Role: "reader",
}
updated, err := svc.Permissions.Update(fileId, permissionId, perm).
    SupportsAllDrives(true).
    UseDomainAdminAccess(true).
    Fields("id,emailAddress,role,type").
    Do()
```

To remove an expiration: `.RemoveExpiration(true)`

## Deleting a Permission

```go
err := svc.Permissions.Delete(fileId, permissionId).
    SupportsAllDrives(true).
    UseDomainAdminAccess(true).
    Do()
```

## This Provider's Pattern

The `resource_drive_permission` uses a composite ID (`fileId/permissionId`) and always sets `SupportsAllDrives(true)`:

```go
created, err := driveSvc.Permissions.Create(plan.FileId.ValueString(), perm).
    UseDomainAdminAccess(plan.UseDomainAdminAccess.ValueBool()).
    SupportsAllDrives(true).
    Fields("id").
    Do()

plan.PermissionId = types.StringValue(created.Id)
plan.Id = types.StringValue(plan.FileId.ValueString() + "/" + created.Id)
```

Import format: `use_domain_admin_access,file_id/permission_id` (e.g., `true,0ABC123/12345`)

## Roles Summary

| Role            | My Drive       | Shared Drive                 |
| --------------- | -------------- | ---------------------------- |
| `owner`         | Full control   | N/A (use organizer)          |
| `organizer`     | N/A            | Manage members + all content |
| `fileOrganizer` | N/A            | Manage content, not members  |
| `writer`        | Edit           | Edit                         |
| `commenter`     | View + comment | View + comment               |
| `reader`        | View           | View                         |

## Full API Reference

- Permission struct fields: `references/api/permission-struct.md`
- All service methods and call options: `references/api/permissions-service.md`

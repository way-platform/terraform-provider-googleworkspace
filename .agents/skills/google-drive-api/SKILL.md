---
name: google-drive-api
description: >
  Google Drive API v3 Go SDK reference (google.golang.org/api/drive/v3). Use
  this skill whenever implementing or modifying code that calls the Drive API:
  files, permissions, shared drives, changes, comments. Also use when the user
  asks about Drive API call patterns, query syntax, struct fields, or
  ForceSendFields. Always use this skill alongside terraform-provider-dev when
  adding new Drive-based Terraform resources.
---

# Google Drive API v3 (Go SDK)

Package: `google.golang.org/api/drive/v3`

## Service Overview

Access services via `svc.Files`, `svc.Drives`, etc. after creating a service:

```go
svc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
```

| Service           | Handles                                                |
| ----------------- | ------------------------------------------------------ |
| `svc.Files`       | File/folder CRUD, upload, download, list, copy, export |
| `svc.Permissions` | Sharing: grant/revoke access, roles                    |
| `svc.Drives`      | Shared drive CRUD, restrictions, hide/unhide           |
| `svc.Changes`     | Track file/drive modifications, push notifications     |
| `svc.Comments`    | File comments                                          |
| `svc.Replies`     | Replies to comments                                    |
| `svc.Revisions`   | File revision history                                  |
| `svc.Channels`    | Stop push notification channels                        |
| `svc.About`       | User/storage quota info                                |

## Universal Call Pattern

Every operation is a builder chain:

```go
result, err := svc.Files.Get(fileId).
    SupportsAllDrives(true).
    Fields("id,name,mimeType").
    Context(ctx).
    Do()
```

## Critical Conventions

### ForceSendFields (boolean pitfall)

All struct fields use `omitempty`. Boolean `false` is silently dropped unless listed:

```go
req := &drive.DriveRestrictions{
    DriveMembersOnly: false,
    ForceSendFields:  []string{"DriveMembersOnly"},
}
```

### SupportsAllDrives

Required for any operation that might touch shared drive items. Without it, shared drive files return 404 or are excluded from list results.

### Fields() for partial responses

Request only the fields you need. For list calls, include `nextPageToken`:

```go
.Fields("files(id,name,mimeType),nextPageToken")
```

### UseDomainAdminAccess

Enables domain admin to act on shared drives they aren't a direct member of.

## Routing Table

| Working on...                                 | Read                                    |
| --------------------------------------------- | --------------------------------------- |
| Service creation, scopes, auth                | `references/guides/authentication.md`   |
| ForceSendFields, Fields(), pagination, errors | `references/guides/conventions.md`      |
| File/folder CRUD, upload, download, queries   | `references/guides/files.md`            |
| Sharing, permissions, roles                   | `references/guides/permissions.md`      |
| Shared drives, restrictions                   | `references/guides/shared-drives.md`    |
| Change tracking, webhooks                     | `references/guides/changes.md`          |
| Comments, replies                             | `references/guides/comments-replies.md` |

## Detailed API Reference

For complete struct fields and all call options:

| Need                                      | Read                                    |
| ----------------------------------------- | --------------------------------------- |
| File struct (50+ fields)                  | `references/api/file-struct.md`         |
| FilesService methods + call options       | `references/api/files-service.md`       |
| Permission struct                         | `references/api/permission-struct.md`   |
| PermissionsService methods + call options | `references/api/permissions-service.md` |
| Drive struct + DrivesService              | `references/api/drives-service.md`      |
| Change struct + ChangesService            | `references/api/changes-service.md`     |
| Comment/Reply + services                  | `references/api/comments-service.md`    |
| OAuth2 scopes                             | `references/api/scopes.md`              |

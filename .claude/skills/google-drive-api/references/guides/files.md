# Files

The `FilesService` handles file and folder CRUD, content upload/download, and querying.

## Creating Files

### Metadata-only (folders, shortcuts)

```go
folder := &drive.File{
    Name:     "My Folder",
    MimeType: "application/vnd.google-apps.folder",
    Parents:  []string{parentFolderId},
}
created, err := svc.Files.Create(folder).Fields("id,name").Do()
```

### With content upload

```go
file := &drive.File{
    Name:    "report.pdf",
    Parents: []string{folderId},
}
f, _ := os.Open("report.pdf")
defer f.Close()

created, err := svc.Files.Create(file).
    Media(f, googleapi.ContentType("application/pdf")).
    Fields("id,name,size").
    Do()
```

For large files, use `ResumableMedia` for resumable uploads:

```go
created, err := svc.Files.Create(file).
    ResumableMedia(ctx, f, fileSize, "application/pdf").
    ProgressUpdater(func(current, total int64) { /* ... */ }).
    Do()
```

## Reading File Metadata

```go
file, err := svc.Files.Get(fileId).
    Fields("id,name,mimeType,size,parents,permissions").
    SupportsAllDrives(true).
    Do()
```

## Downloading File Content

```go
resp, err := svc.Files.Get(fileId).Download()
if err != nil {
    return err
}
defer resp.Body.Close()
io.Copy(dst, resp.Body)
```

For Google Workspace docs, use Export:

```go
resp, err := svc.Files.Export(fileId, "application/pdf").Download()
```

## Listing and Querying

```go
result, err := svc.Files.List().
    Q("'folder-id' in parents and trashed = false").
    Fields("files(id,name,mimeType),nextPageToken").
    PageSize(100).
    OrderBy("name").
    SupportsAllDrives(true).
    IncludeItemsFromAllDrives(true).
    Do()
```

### Query syntax

```
// By parent
'parent-id' in parents

// By type
mimeType = 'application/vnd.google-apps.folder'

// By name (exact or contains)
name = 'budget.xlsx'
name contains 'budget'

// Combine with and/or/not
mimeType != 'application/vnd.google-apps.folder' and trashed = false

// In a shared drive
driveId = 'drive-id'
```

### Pagination

```go
var all []*drive.File
err := svc.Files.List().Q(query).PageSize(100).
    Pages(ctx, func(list *drive.FileList) error {
        all = append(all, list.Files...)
        return nil
    })
```

### Corpora

Controls which files are searched:

| Value         | Scope                                                 |
| ------------- | ----------------------------------------------------- |
| `"user"`      | Files owned by or shared with the user (default)      |
| `"drive"`     | Files in a specific shared drive (requires `DriveId`) |
| `"allDrives"` | All drives the user has access to                     |
| `"domain"`    | Files shared to the user's domain                     |

```go
svc.Files.List().
    Corpora("drive").
    DriveId(driveId).
    IncludeItemsFromAllDrives(true).
    SupportsAllDrives(true).
    Do()
```

## Updating Files

### Metadata only

```go
update := &drive.File{
    Name: "new-name.txt",
}
updated, err := svc.Files.Update(fileId, update).
    Fields("id,name").
    Do()
```

### Moving between folders

```go
updated, err := svc.Files.Update(fileId, &drive.File{}).
    AddParents(newParentId).
    RemoveParents(oldParentId).
    Fields("id,parents").
    Do()
```

### Updating content

```go
updated, err := svc.Files.Update(fileId, &drive.File{}).
    Media(reader, googleapi.ContentType("text/plain")).
    Do()
```

## Deleting Files

Permanently deletes (skips trash):

```go
err := svc.Files.Delete(fileId).SupportsAllDrives(true).Do()
```

To trash instead, update the `Trashed` field:

```go
svc.Files.Update(fileId, &drive.File{Trashed: true, ForceSendFields: []string{"Trashed"}}).Do()
```

## Copying Files

```go
copy := &drive.File{
    Name:    "Copy of document",
    Parents: []string{destFolderId},
}
copied, err := svc.Files.Copy(sourceFileId, copy).
    Fields("id,name").
    SupportsAllDrives(true).
    Do()
```

## Full API Reference

- File struct fields: `references/api/file-struct.md`
- All service methods and call options: `references/api/files-service.md`

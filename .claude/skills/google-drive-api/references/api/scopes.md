# OAuth2 Scopes

```go
import "google.golang.org/api/drive/v3"
```

| Constant                     | Value                                                     | Description                                                                              |
| ---------------------------- | --------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `DriveScope`                 | `https://www.googleapis.com/auth/drive`                   | See, edit, create, and delete all of your Google Drive files                             |
| `DriveAppdataScope`          | `https://www.googleapis.com/auth/drive.appdata`           | See, create, and delete its own configuration data in your Google Drive                  |
| `DriveFileScope`             | `https://www.googleapis.com/auth/drive.file`              | See, edit, create, and delete only the specific Google Drive files you use with this app |
| `DriveMetadataScope`         | `https://www.googleapis.com/auth/drive.metadata`          | View and manage metadata of files in your Google Drive                                   |
| `DriveMetadataReadonlyScope` | `https://www.googleapis.com/auth/drive.metadata.readonly` | See information about your Google Drive files                                            |
| `DrivePhotosReadonlyScope`   | `https://www.googleapis.com/auth/drive.photos.readonly`   | View the photos, videos and albums in your Google Photos                                 |
| `DriveReadonlyScope`         | `https://www.googleapis.com/auth/drive.readonly`          | See and download all your Google Drive files                                             |
| `DriveScriptsScope`          | `https://www.googleapis.com/auth/drive.scripts`           | Modify your Google Apps Script scripts' behavior                                         |

## Scope Selection

- `DriveScope` — full access, use for admin/service-account scenarios
- `DriveFileScope` — limited to files the app created or the user opened with the app
- `DriveMetadataScope` — metadata only, no file content access
- `DriveReadonlyScope` — read-only access to all files

Restrict scopes with `option.WithScopes`:

```go
svc, err := drive.NewService(ctx, option.WithScopes(drive.DriveScope))
```

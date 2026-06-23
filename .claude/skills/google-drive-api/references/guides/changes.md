# Changes (Change Tracking)

The `ChangesService` tracks modifications to files and shared drives. Use it for sync/audit scenarios.

## Polling Pattern

1. Get a start page token (marks "now")
2. Periodically list changes since that token
3. When `NewStartPageToken` appears in the response, use it for the next poll

```go
// Step 1: Get start token
startResp, err := svc.Changes.GetStartPageToken().
    SupportsAllDrives(true).
    Do()
savedToken := startResp.StartPageToken

// Step 2: Poll for changes
pageToken := savedToken
for {
    result, err := svc.Changes.List(pageToken).
        SupportsAllDrives(true).
        IncludeItemsFromAllDrives(true).
        PageSize(100).
        Fields("changes(fileId,changeType,file(id,name,trashed),removed,time),newStartPageToken,nextPageToken").
        Do()
    if err != nil {
        return err
    }

    for _, change := range result.Changes {
        // Process change...
    }

    if result.NewStartPageToken != "" {
        // All current changes consumed; save for next poll
        savedToken = result.NewStartPageToken
        break
    }
    pageToken = result.NextPageToken
}
```

## Change Types

| `ChangeType` | `File` field             | `Drive` field                  | Meaning                   |
| ------------ | ------------------------ | ------------------------------ | ------------------------- |
| `"file"`     | Present (if not removed) | —                              | File was modified/created |
| `"drive"`    | —                        | Present (if user still member) | Shared drive was modified |

When `Removed` is true, the file/drive was deleted or access was lost.

## Shared Drive Changes

To track changes in a specific shared drive:

```go
svc.Changes.GetStartPageToken().DriveId(driveId).SupportsAllDrives(true).Do()
svc.Changes.List(token).DriveId(driveId).SupportsAllDrives(true).Do()
```

## Push Notifications (Watch)

Instead of polling, subscribe to push notifications:

```go
channel := &drive.Channel{
    Id:      "unique-channel-id",  // UUID you generate
    Type:    "web_hook",
    Address: "https://example.com/webhook",
    Token:   "optional-verification-token",
}

resp, err := svc.Changes.Watch(pageToken, channel).
    SupportsAllDrives(true).
    IncludeItemsFromAllDrives(true).
    Do()
// resp.ResourceId and resp.Expiration are set by the server
```

Stop a channel:

```go
svc.Channels.Stop(&drive.Channel{
    Id:         channelId,
    ResourceId: resourceId,
}).Do()
```

## List Options

| Option                            | Description                                   |
| --------------------------------- | --------------------------------------------- |
| `IncludeItemsFromAllDrives(true)` | Include shared drive items                    |
| `SupportsAllDrives(true)`         | Required with shared drives                   |
| `IncludeRemoved(bool)`            | Include removal changes (default true)        |
| `IncludeCorpusRemovals(bool)`     | Include changes for items removed from corpus |
| `RestrictToMyDrive(bool)`         | Only My Drive changes                         |
| `Spaces(string)`                  | Filter by space: `"drive"`, `"appDataFolder"` |

## Full API Reference

- Change/ChangeList/Channel structs, all call options: `references/api/changes-service.md`

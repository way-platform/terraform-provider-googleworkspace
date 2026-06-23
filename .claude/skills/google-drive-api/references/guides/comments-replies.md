# Comments & Replies

The `CommentsService` and `RepliesService` manage discussions on files.

## Creating a Comment

```go
comment := &drive.Comment{
    Content: "This section needs revision.",
}
created, err := svc.Comments.Create(fileId, comment).
    Fields("id,content,author,createdTime").
    Do()
```

With an anchor (document region):

```go
comment := &drive.Comment{
    Content: "Typo here.",
    Anchor:  `{"r":"head","a":[{"line":{"n":5,"l":1}}]}`,
}
```

## Listing Comments

```go
var all []*drive.Comment
err := svc.Comments.List(fileId).
    PageSize(100).
    IncludeDeleted(false).
    Fields("comments(id,content,author,resolved,replies),nextPageToken").
    Pages(ctx, func(list *drive.CommentList) error {
        all = append(all, list.Comments...)
        return nil
    })
```

Filter by modification time:

```go
svc.Comments.List(fileId).StartModifiedTime("2024-01-01T00:00:00Z").Do()
```

## Getting a Comment

```go
comment, err := svc.Comments.Get(fileId, commentId).
    Fields("id,content,author,resolved,replies").
    IncludeDeleted(false).
    Do()
```

## Updating a Comment

Only `Content` can be updated:

```go
update := &drive.Comment{
    Content: "Updated comment text.",
}
updated, err := svc.Comments.Update(fileId, commentId, update).
    Fields("id,content,modifiedTime").
    Do()
```

## Deleting a Comment

```go
err := svc.Comments.Delete(fileId, commentId).Do()
```

## Creating a Reply

```go
reply := &drive.Reply{
    Content: "Fixed, thanks!",
}
created, err := svc.Replies.Create(fileId, commentId, reply).
    Fields("id,content,author").
    Do()
```

### Resolving/Reopening a Comment

Use `Action` instead of `Content`:

```go
// Resolve
svc.Replies.Create(fileId, commentId, &drive.Reply{Action: "resolve"}).Do()

// Reopen
svc.Replies.Create(fileId, commentId, &drive.Reply{Action: "reopen"}).Do()
```

## Listing Replies

```go
var all []*drive.Reply
err := svc.Replies.List(fileId, commentId).
    PageSize(100).
    Pages(ctx, func(list *drive.ReplyList) error {
        all = append(all, list.Replies...)
        return nil
    })
```

## Updating a Reply

```go
update := &drive.Reply{Content: "Actually, let me reconsider."}
svc.Replies.Update(fileId, commentId, replyId, update).Do()
```

## Deleting a Reply

```go
svc.Replies.Delete(fileId, commentId, replyId).Do()
```

## Full API Reference

- Comment/Reply structs, all call options: `references/api/comments-service.md`

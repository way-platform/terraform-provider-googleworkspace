# CommentsService & RepliesService

## Comment Struct

```go
type Comment struct {
    Anchor            string                    `json:"anchor,omitempty"`
    Author            *User                     `json:"author,omitempty"`
    Content           string                    `json:"content,omitempty"`
    CreatedTime       string                    `json:"createdTime,omitempty"`
    Deleted           bool                      `json:"deleted,omitempty"`
    HtmlContent       string                    `json:"htmlContent,omitempty"`
    Id                string                    `json:"id,omitempty"`
    Kind              string                    `json:"kind,omitempty"`
    ModifiedTime      string                    `json:"modifiedTime,omitempty"`
    QuotedFileContent *CommentQuotedFileContent `json:"quotedFileContent,omitempty"`
    Replies           []*Reply                  `json:"replies,omitempty"`
    Resolved          bool                      `json:"resolved,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

| Field         | Writable    | Notes                                    |
| ------------- | ----------- | ---------------------------------------- |
| `Content`     | Yes         | Plain text; use this for setting content |
| `HtmlContent` | No (output) | HTML-formatted version for display       |
| `Anchor`      | Yes         | JSON string defining document region     |
| `Resolved`    | No (output) | Set via reply with `Action: "resolve"`   |

## Reply Struct

```go
type Reply struct {
    Action      string `json:"action,omitempty"`
    Author      *User  `json:"author,omitempty"`
    Content     string `json:"content,omitempty"`
    CreatedTime string `json:"createdTime,omitempty"`
    Deleted     bool   `json:"deleted,omitempty"`
    HtmlContent string `json:"htmlContent,omitempty"`
    Id          string `json:"id,omitempty"`
    Kind        string `json:"kind,omitempty"`
    ModifiedTime string `json:"modifiedTime,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

| Field     | Notes                                                         |
| --------- | ------------------------------------------------------------- |
| `Action`  | `"resolve"` or `"reopen"` (mutually exclusive with `Content`) |
| `Content` | Plain text; required on create if no `Action`                 |

## CommentsService

```go
type CommentsService struct{}

func NewCommentsService(s *Service) *CommentsService
func (r *CommentsService) Create(fileId string, comment *Comment) *CommentsCreateCall
func (r *CommentsService) Delete(fileId string, commentId string) *CommentsDeleteCall
func (r *CommentsService) Get(fileId string, commentId string) *CommentsGetCall
func (r *CommentsService) List(fileId string) *CommentsListCall
func (r *CommentsService) Update(fileId string, commentId string, comment *Comment) *CommentsUpdateCall
```

## CommentsListCall

```go
func (c *CommentsListCall) Context(ctx context.Context) *CommentsListCall
func (c *CommentsListCall) Do(opts ...googleapi.CallOption) (*CommentList, error)
func (c *CommentsListCall) Fields(s ...googleapi.Field) *CommentsListCall
func (c *CommentsListCall) IfNoneMatch(entityTag string) *CommentsListCall
func (c *CommentsListCall) IncludeDeleted(v bool) *CommentsListCall
func (c *CommentsListCall) PageSize(v int64) *CommentsListCall
func (c *CommentsListCall) PageToken(v string) *CommentsListCall
func (c *CommentsListCall) Pages(ctx context.Context, f func(*CommentList) error) error
func (c *CommentsListCall) StartModifiedTime(v string) *CommentsListCall
```

## RepliesService

```go
type RepliesService struct{}

func NewRepliesService(s *Service) *RepliesService
func (r *RepliesService) Create(fileId string, commentId string, reply *Reply) *RepliesCreateCall
func (r *RepliesService) Delete(fileId string, commentId string, replyId string) *RepliesDeleteCall
func (r *RepliesService) Get(fileId string, commentId string, replyId string) *RepliesGetCall
func (r *RepliesService) List(fileId string, commentId string) *RepliesListCall
func (r *RepliesService) Update(fileId string, commentId string, replyId string, reply *Reply) *RepliesUpdateCall
```

## RepliesListCall

```go
func (c *RepliesListCall) Context(ctx context.Context) *RepliesListCall
func (c *RepliesListCall) Do(opts ...googleapi.CallOption) (*ReplyList, error)
func (c *RepliesListCall) Fields(s ...googleapi.Field) *RepliesListCall
func (c *RepliesListCall) IfNoneMatch(entityTag string) *RepliesListCall
func (c *RepliesListCall) IncludeDeleted(v bool) *RepliesListCall
func (c *RepliesListCall) PageSize(v int64) *RepliesListCall
func (c *RepliesListCall) PageToken(v string) *RepliesListCall
func (c *RepliesListCall) Pages(ctx context.Context, f func(*ReplyList) error) error
```

## CommentList / ReplyList

```go
type CommentList struct {
    Comments      []*Comment `json:"comments,omitempty"`
    Kind          string     `json:"kind,omitempty"`
    NextPageToken string     `json:"nextPageToken,omitempty"`
}

type ReplyList struct {
    Kind          string   `json:"kind,omitempty"`
    NextPageToken string   `json:"nextPageToken,omitempty"`
    Replies       []*Reply `json:"replies,omitempty"`
}
```

# Member Struct & MembersService

## Member Struct

```go
type Member struct {
    // --- Writable fields ---

    Email            string `json:"email,omitempty"`            // Required for Insert; user or group email
    Role             string `json:"role,omitempty"`             // "OWNER", "MANAGER", "MEMBER"
    DeliverySettings string `json:"deliverySettings,omitempty"` // Mail delivery preference
    Type             string `json:"type,omitempty"`             // "USER", "GROUP", "EXTERNAL"

    // --- Read-only fields ---

    Id     string `json:"id,omitempty"`     // Unique member ID (usable as memberKey)
    Status string `json:"status,omitempty"` // Membership status
    Etag   string `json:"etag,omitempty"`
    Kind   string `json:"kind,omitempty"`

    // --- JSON control ---
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

### DeliverySettings Values

| Value        | Behavior                             |
| ------------ | ------------------------------------ |
| `"ALL_MAIL"` | All messages delivered               |
| `"DAILY"`    | Daily digest                         |
| `"DIGEST"`   | Up to 25 messages bundled            |
| `"DISABLED"` | No email delivery                    |
| `"NONE"`     | No delivery preference set (default) |

### Role Values

| Value       | Permissions                              |
| ----------- | ---------------------------------------- |
| `"OWNER"`   | Full control including member management |
| `"MANAGER"` | Can manage members and settings          |
| `"MEMBER"`  | Standard group membership                |

## Members (List Response)

```go
type Members struct {
    Members       []*Member `json:"members,omitempty"`
    NextPageToken string    `json:"nextPageToken,omitempty"`
    Etag          string    `json:"etag,omitempty"`
    Kind          string    `json:"kind,omitempty"`
}
```

## MembersHasMember (Response)

```go
type MembersHasMember struct {
    IsMember bool `json:"isMember,omitempty"` // true if member exists
}
```

---

## MembersService

```go
type MembersService struct{}

func NewMembersService(s *Service) *MembersService
func (r *MembersService) Delete(groupKey string, memberKey string) *MembersDeleteCall
func (r *MembersService) Get(groupKey string, memberKey string) *MembersGetCall
func (r *MembersService) HasMember(groupKey string, memberKey string) *MembersHasMemberCall
func (r *MembersService) Insert(groupKey string, member *Member) *MembersInsertCall
func (r *MembersService) List(groupKey string) *MembersListCall
func (r *MembersService) Patch(groupKey string, memberKey string, member *Member) *MembersPatchCall
func (r *MembersService) Update(groupKey string, memberKey string, member *Member) *MembersUpdateCall
```

## MembersGetCall

```go
func (c *MembersGetCall) Context(ctx context.Context) *MembersGetCall
func (c *MembersGetCall) Do(opts ...googleapi.CallOption) (*Member, error)
func (c *MembersGetCall) Fields(s ...googleapi.Field) *MembersGetCall
func (c *MembersGetCall) IfNoneMatch(entityTag string) *MembersGetCall
```

## MembersInsertCall

```go
func (c *MembersInsertCall) Context(ctx context.Context) *MembersInsertCall
func (c *MembersInsertCall) Do(opts ...googleapi.CallOption) (*Member, error)
func (c *MembersInsertCall) Fields(s ...googleapi.Field) *MembersInsertCall
```

## MembersListCall

```go
func (c *MembersListCall) Context(ctx context.Context) *MembersListCall
func (c *MembersListCall) Do(opts ...googleapi.CallOption) (*Members, error)
func (c *MembersListCall) Fields(s ...googleapi.Field) *MembersListCall
func (c *MembersListCall) IfNoneMatch(entityTag string) *MembersListCall
func (c *MembersListCall) IncludeDerivedMembership(v bool) *MembersListCall  // Include indirect members
func (c *MembersListCall) MaxResults(maxResults int64) *MembersListCall      // 1-200, default 200
func (c *MembersListCall) PageToken(pageToken string) *MembersListCall
func (c *MembersListCall) Pages(ctx context.Context, f func(*Members) error) error
func (c *MembersListCall) Roles(roles string) *MembersListCall               // Filter: "OWNER", "MANAGER", "MEMBER"
```

## MembersUpdateCall / MembersPatchCall

```go
func (c *MembersUpdateCall) Context(ctx context.Context) *MembersUpdateCall
func (c *MembersUpdateCall) Do(opts ...googleapi.CallOption) (*Member, error)
func (c *MembersUpdateCall) Fields(s ...googleapi.Field) *MembersUpdateCall

func (c *MembersPatchCall) Context(ctx context.Context) *MembersPatchCall
func (c *MembersPatchCall) Do(opts ...googleapi.CallOption) (*Member, error)
func (c *MembersPatchCall) Fields(s ...googleapi.Field) *MembersPatchCall
```

## MembersDeleteCall

```go
func (c *MembersDeleteCall) Context(ctx context.Context) *MembersDeleteCall
func (c *MembersDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *MembersDeleteCall) Fields(s ...googleapi.Field) *MembersDeleteCall
```

## MembersHasMemberCall

```go
func (c *MembersHasMemberCall) Context(ctx context.Context) *MembersHasMemberCall
func (c *MembersHasMemberCall) Do(opts ...googleapi.CallOption) (*MembersHasMember, error)
func (c *MembersHasMemberCall) Fields(s ...googleapi.Field) *MembersHasMemberCall
func (c *MembersHasMemberCall) IfNoneMatch(entityTag string) *MembersHasMemberCall
```

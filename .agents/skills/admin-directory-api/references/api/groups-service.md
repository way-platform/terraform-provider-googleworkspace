# GroupsService

```go
type GroupsService struct {
    Aliases *GroupsAliasesService
}

func NewGroupsService(s *Service) *GroupsService
func (r *GroupsService) Delete(groupKey string) *GroupsDeleteCall
func (r *GroupsService) Get(groupKey string) *GroupsGetCall
func (r *GroupsService) Insert(group *Group) *GroupsInsertCall
func (r *GroupsService) List() *GroupsListCall
func (r *GroupsService) Patch(groupKey string, group *Group) *GroupsPatchCall
func (r *GroupsService) Update(groupKey string, group *Group) *GroupsUpdateCall
```

## GroupsGetCall

```go
func (c *GroupsGetCall) Context(ctx context.Context) *GroupsGetCall
func (c *GroupsGetCall) Do(opts ...googleapi.CallOption) (*Group, error)
func (c *GroupsGetCall) Fields(s ...googleapi.Field) *GroupsGetCall
func (c *GroupsGetCall) IfNoneMatch(entityTag string) *GroupsGetCall
```

## GroupsInsertCall

```go
func (c *GroupsInsertCall) Context(ctx context.Context) *GroupsInsertCall
func (c *GroupsInsertCall) Do(opts ...googleapi.CallOption) (*Group, error)
func (c *GroupsInsertCall) Fields(s ...googleapi.Field) *GroupsInsertCall
```

## GroupsListCall

```go
func (c *GroupsListCall) Context(ctx context.Context) *GroupsListCall
func (c *GroupsListCall) Customer(customer string) *GroupsListCall     // "my_customer" or ID
func (c *GroupsListCall) Do(opts ...googleapi.CallOption) (*Groups, error)
func (c *GroupsListCall) Domain(domain string) *GroupsListCall
func (c *GroupsListCall) Fields(s ...googleapi.Field) *GroupsListCall
func (c *GroupsListCall) IfNoneMatch(entityTag string) *GroupsListCall
func (c *GroupsListCall) MaxResults(maxResults int64) *GroupsListCall  // 1-200, default 200
func (c *GroupsListCall) OrderBy(orderBy string) *GroupsListCall       // "email"
func (c *GroupsListCall) PageToken(pageToken string) *GroupsListCall
func (c *GroupsListCall) Pages(ctx context.Context, f func(*Groups) error) error
func (c *GroupsListCall) Query(query string) *GroupsListCall
func (c *GroupsListCall) SortOrder(sortOrder string) *GroupsListCall   // "ASCENDING", "DESCENDING"
func (c *GroupsListCall) UserKey(userKey string) *GroupsListCall       // Filter groups for a user
```

## GroupsUpdateCall / GroupsPatchCall

```go
func (c *GroupsUpdateCall) Context(ctx context.Context) *GroupsUpdateCall
func (c *GroupsUpdateCall) Do(opts ...googleapi.CallOption) (*Group, error)
func (c *GroupsUpdateCall) Fields(s ...googleapi.Field) *GroupsUpdateCall

func (c *GroupsPatchCall) Context(ctx context.Context) *GroupsPatchCall
func (c *GroupsPatchCall) Do(opts ...googleapi.CallOption) (*Group, error)
func (c *GroupsPatchCall) Fields(s ...googleapi.Field) *GroupsPatchCall
```

## GroupsDeleteCall

```go
func (c *GroupsDeleteCall) Context(ctx context.Context) *GroupsDeleteCall
func (c *GroupsDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *GroupsDeleteCall) Fields(s ...googleapi.Field) *GroupsDeleteCall
```

---

## GroupsAliasesService

```go
type GroupsAliasesService struct{}

func NewGroupsAliasesService(s *Service) *GroupsAliasesService
func (r *GroupsAliasesService) Delete(groupKey string, alias string) *GroupsAliasesDeleteCall
func (r *GroupsAliasesService) Insert(groupKey string, alias *Alias) *GroupsAliasesInsertCall
func (r *GroupsAliasesService) List(groupKey string) *GroupsAliasesListCall
```

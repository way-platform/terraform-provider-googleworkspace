# UsersService

```go
type UsersService struct {
    Aliases *UsersAliasesService
    Photos  *UsersPhotosService
}

func NewUsersService(s *Service) *UsersService
func (r *UsersService) Delete(userKey string) *UsersDeleteCall
func (r *UsersService) Get(userKey string) *UsersGetCall
func (r *UsersService) Insert(user *User) *UsersInsertCall
func (r *UsersService) List() *UsersListCall
func (r *UsersService) MakeAdmin(userKey string, usermakeadmin *UserMakeAdmin) *UsersMakeAdminCall
func (r *UsersService) Patch(userKey string, user *User) *UsersPatchCall
func (r *UsersService) SignOut(userKey string) *UsersSignOutCall
func (r *UsersService) Undelete(userKey string, userundelete *UserUndelete) *UsersUndeleteCall
func (r *UsersService) Update(userKey string, user *User) *UsersUpdateCall
func (r *UsersService) Watch(channel *Channel) *UsersWatchCall
```

## UsersGetCall

```go
func (c *UsersGetCall) Context(ctx context.Context) *UsersGetCall
func (c *UsersGetCall) CustomFieldMask(customFieldMask string) *UsersGetCall
func (c *UsersGetCall) Do(opts ...googleapi.CallOption) (*User, error)
func (c *UsersGetCall) Fields(s ...googleapi.Field) *UsersGetCall
func (c *UsersGetCall) IfNoneMatch(entityTag string) *UsersGetCall
func (c *UsersGetCall) Projection(projection string) *UsersGetCall   // "basic", "custom", "full"
func (c *UsersGetCall) ViewType(viewType string) *UsersGetCall       // "admin_view", "domain_public"
```

## UsersInsertCall

```go
func (c *UsersInsertCall) Context(ctx context.Context) *UsersInsertCall
func (c *UsersInsertCall) Do(opts ...googleapi.CallOption) (*User, error)
func (c *UsersInsertCall) Fields(s ...googleapi.Field) *UsersInsertCall
func (c *UsersInsertCall) ResolveConflictAccount(v bool) *UsersInsertCall
```

## UsersListCall

```go
func (c *UsersListCall) Context(ctx context.Context) *UsersListCall
func (c *UsersListCall) CustomFieldMask(customFieldMask string) *UsersListCall
func (c *UsersListCall) Customer(customer string) *UsersListCall     // "my_customer" or ID
func (c *UsersListCall) Do(opts ...googleapi.CallOption) (*Users, error)
func (c *UsersListCall) Domain(domain string) *UsersListCall
func (c *UsersListCall) Event(event string) *UsersListCall
func (c *UsersListCall) Fields(s ...googleapi.Field) *UsersListCall
func (c *UsersListCall) IfNoneMatch(entityTag string) *UsersListCall
func (c *UsersListCall) MaxResults(maxResults int64) *UsersListCall  // 1-500, default 100
func (c *UsersListCall) OrderBy(orderBy string) *UsersListCall       // "email", "familyName", "givenName"
func (c *UsersListCall) PageToken(pageToken string) *UsersListCall
func (c *UsersListCall) Pages(ctx context.Context, f func(*Users) error) error
func (c *UsersListCall) Projection(projection string) *UsersListCall // "basic", "custom", "full"
func (c *UsersListCall) Query(query string) *UsersListCall
func (c *UsersListCall) ShowDeleted(showDeleted string) *UsersListCall
func (c *UsersListCall) SortOrder(sortOrder string) *UsersListCall   // "ASCENDING", "DESCENDING"
func (c *UsersListCall) ViewType(viewType string) *UsersListCall
```

## UsersUpdateCall / UsersPatchCall

```go
// Update replaces entire resource; Patch modifies only provided fields
func (c *UsersUpdateCall) Context(ctx context.Context) *UsersUpdateCall
func (c *UsersUpdateCall) Do(opts ...googleapi.CallOption) (*User, error)
func (c *UsersUpdateCall) Fields(s ...googleapi.Field) *UsersUpdateCall

func (c *UsersPatchCall) Context(ctx context.Context) *UsersPatchCall
func (c *UsersPatchCall) Do(opts ...googleapi.CallOption) (*User, error)
func (c *UsersPatchCall) Fields(s ...googleapi.Field) *UsersPatchCall
```

## UsersDeleteCall

```go
func (c *UsersDeleteCall) Context(ctx context.Context) *UsersDeleteCall
func (c *UsersDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *UsersDeleteCall) Fields(s ...googleapi.Field) *UsersDeleteCall
```

## UsersMakeAdminCall

```go
func (c *UsersMakeAdminCall) Context(ctx context.Context) *UsersMakeAdminCall
func (c *UsersMakeAdminCall) Do(opts ...googleapi.CallOption) error
func (c *UsersMakeAdminCall) Fields(s ...googleapi.Field) *UsersMakeAdminCall
```

## UsersSignOutCall

```go
func (c *UsersSignOutCall) Context(ctx context.Context) *UsersSignOutCall
func (c *UsersSignOutCall) Do(opts ...googleapi.CallOption) error
func (c *UsersSignOutCall) Fields(s ...googleapi.Field) *UsersSignOutCall
```

## UsersUndeleteCall

```go
func (c *UsersUndeleteCall) Context(ctx context.Context) *UsersUndeleteCall
func (c *UsersUndeleteCall) Do(opts ...googleapi.CallOption) error
func (c *UsersUndeleteCall) Fields(s ...googleapi.Field) *UsersUndeleteCall
```

---

## UsersAliasesService

```go
type UsersAliasesService struct{}

func NewUsersAliasesService(s *Service) *UsersAliasesService
func (r *UsersAliasesService) Delete(userKey string, alias string) *UsersAliasesDeleteCall
func (r *UsersAliasesService) Insert(userKey string, alias *Alias) *UsersAliasesInsertCall
func (r *UsersAliasesService) List(userKey string) *UsersAliasesListCall
func (r *UsersAliasesService) Watch(userKey string, channel *Channel) *UsersAliasesWatchCall
```

## UsersPhotosService

```go
type UsersPhotosService struct{}

func NewUsersPhotosService(s *Service) *UsersPhotosService
func (r *UsersPhotosService) Delete(userKey string) *UsersPhotosDeleteCall
func (r *UsersPhotosService) Get(userKey string) *UsersPhotosGetCall
func (r *UsersPhotosService) Patch(userKey string, userphoto *UserPhoto) *UsersPhotosPatchCall
func (r *UsersPhotosService) Update(userKey string, userphoto *UserPhoto) *UsersPhotosUpdateCall
```

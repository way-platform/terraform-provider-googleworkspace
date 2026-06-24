# Role & RoleAssignment Structs + Services

## Role Struct

```go
type Role struct {
    Etag             string               `json:"etag,omitempty"`
    IsSuperAdminRole bool                 `json:"isSuperAdminRole,omitempty"`
    IsSystemRole     bool                 `json:"isSystemRole,omitempty"`
    Kind             string               `json:"kind,omitempty"`
    RoleDescription  string               `json:"roleDescription,omitempty"`
    RoleId           int64                `json:"roleId,omitempty,string"`  // Note: int64, serialized as string
    RoleName         string               `json:"roleName,omitempty"`
    RolePrivileges   []*RoleRolePrivileges `json:"rolePrivileges,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## RoleRolePrivileges

```go
type RoleRolePrivileges struct {
    PrivilegeName string `json:"privilegeName,omitempty"`
    ServiceId     string `json:"serviceId,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## RoleAssignment Struct

```go
type RoleAssignment struct {
    AssignedTo       string `json:"assignedTo,omitempty"`       // user_id, group_id, or service account uniqueId
    AssigneeType     string `json:"assigneeType,omitempty"`     // "user" or "group" (read-only)
    Etag             string `json:"etag,omitempty"`
    Kind             string `json:"kind,omitempty"`
    OrgUnitId        string `json:"orgUnitId,omitempty"`        // Restricts role to this org unit
    RoleAssignmentId int64  `json:"roleAssignmentId,omitempty,string"` // Note: int64
    RoleId           int64  `json:"roleId,omitempty,string"`           // Note: int64
    ScopeType        string `json:"scopeType,omitempty"`        // "CUSTOMER" or "ORG_UNIT"

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

### ScopeType Values

| Value        | Meaning                                     |
| ------------ | ------------------------------------------- |
| `"CUSTOMER"` | Role applies to the entire customer account |
| `"ORG_UNIT"` | Role restricted to the specified OrgUnitId  |

## Roles (List Response)

```go
type Roles struct {
    Items         []*Role `json:"items,omitempty"`
    NextPageToken string  `json:"nextPageToken,omitempty"`
    Etag          string  `json:"etag,omitempty"`
    Kind          string  `json:"kind,omitempty"`
}
```

## RoleAssignments (List Response)

```go
type RoleAssignments struct {
    Items         []*RoleAssignment `json:"items,omitempty"`
    NextPageToken string            `json:"nextPageToken,omitempty"`
    Etag          string            `json:"etag,omitempty"`
    Kind          string            `json:"kind,omitempty"`
}
```

---

## RolesService

All methods require `customer` as the first parameter. Use `"my_customer"` as alias.

```go
type RolesService struct{}

func NewRolesService(s *Service) *RolesService
func (r *RolesService) Delete(customer string, roleId string) *RolesDeleteCall
func (r *RolesService) Get(customer string, roleId string) *RolesGetCall
func (r *RolesService) Insert(customer string, role *Role) *RolesInsertCall
func (r *RolesService) List(customer string) *RolesListCall
func (r *RolesService) Patch(customer string, roleId string, role *Role) *RolesPatchCall
func (r *RolesService) Update(customer string, roleId string, role *Role) *RolesUpdateCall
```

Note: `roleId` parameter is a string in method signatures but represents the int64 RoleId. Use `strconv.FormatInt(role.RoleId, 10)` to convert.

## RolesListCall

```go
func (c *RolesListCall) Context(ctx context.Context) *RolesListCall
func (c *RolesListCall) Do(opts ...googleapi.CallOption) (*Roles, error)
func (c *RolesListCall) Fields(s ...googleapi.Field) *RolesListCall
func (c *RolesListCall) IfNoneMatch(entityTag string) *RolesListCall
func (c *RolesListCall) MaxResults(maxResults int64) *RolesListCall
func (c *RolesListCall) PageToken(pageToken string) *RolesListCall
func (c *RolesListCall) Pages(ctx context.Context, f func(*Roles) error) error
```

---

## RoleAssignmentsService

```go
type RoleAssignmentsService struct{}

func NewRoleAssignmentsService(s *Service) *RoleAssignmentsService
func (r *RoleAssignmentsService) Delete(customer string, roleAssignmentId string) *RoleAssignmentsDeleteCall
func (r *RoleAssignmentsService) Get(customer string, roleAssignmentId string) *RoleAssignmentsGetCall
func (r *RoleAssignmentsService) Insert(customer string, roleassignment *RoleAssignment) *RoleAssignmentsInsertCall
func (r *RoleAssignmentsService) List(customer string) *RoleAssignmentsListCall
```

Note: `roleAssignmentId` parameter is a string but represents the int64 RoleAssignmentId.

## RoleAssignmentsListCall

```go
func (c *RoleAssignmentsListCall) Context(ctx context.Context) *RoleAssignmentsListCall
func (c *RoleAssignmentsListCall) Do(opts ...googleapi.CallOption) (*RoleAssignments, error)
func (c *RoleAssignmentsListCall) Fields(s ...googleapi.Field) *RoleAssignmentsListCall
func (c *RoleAssignmentsListCall) IfNoneMatch(entityTag string) *RoleAssignmentsListCall
func (c *RoleAssignmentsListCall) IncludeIndirectRoleAssignments(v bool) *RoleAssignmentsListCall
func (c *RoleAssignmentsListCall) MaxResults(maxResults int64) *RoleAssignmentsListCall
func (c *RoleAssignmentsListCall) PageToken(pageToken string) *RoleAssignmentsListCall
func (c *RoleAssignmentsListCall) Pages(ctx context.Context, f func(*RoleAssignments) error) error
func (c *RoleAssignmentsListCall) RoleId(roleId string) *RoleAssignmentsListCall     // Filter by role
func (c *RoleAssignmentsListCall) UserKey(userKey string) *RoleAssignmentsListCall   // Filter by assignee
```

## int64 ID Handling

RoleId and RoleAssignmentId are `int64` in Go structs but passed as strings in method signatures:

```go
import "strconv"

// Struct → method parameter
roleIdStr := strconv.FormatInt(role.RoleId, 10)
svc.Roles.Get("my_customer", roleIdStr).Do()

// String → struct field
roleId, err := strconv.ParseInt(roleIdStr, 10, 64)
```

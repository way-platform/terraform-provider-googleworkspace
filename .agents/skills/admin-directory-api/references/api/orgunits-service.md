# OrgunitsService

All methods require `customerId` as the first parameter. Use `"my_customer"` as alias.

```go
type OrgunitsService struct{}

func NewOrgunitsService(s *Service) *OrgunitsService
func (r *OrgunitsService) Delete(customerId string, orgUnitPath string) *OrgunitsDeleteCall
func (r *OrgunitsService) Get(customerId string, orgUnitPath string) *OrgunitsGetCall
func (r *OrgunitsService) Insert(customerId string, orgunit *OrgUnit) *OrgunitsInsertCall
func (r *OrgunitsService) List(customerId string) *OrgunitsListCall
func (r *OrgunitsService) Patch(customerId string, orgUnitPath string, orgunit *OrgUnit) *OrgunitsPatchCall
func (r *OrgunitsService) Update(customerId string, orgUnitPath string, orgunit *OrgUnit) *OrgunitsUpdateCall
```

## OrgunitsGetCall

```go
func (c *OrgunitsGetCall) Context(ctx context.Context) *OrgunitsGetCall
func (c *OrgunitsGetCall) Do(opts ...googleapi.CallOption) (*OrgUnit, error)
func (c *OrgunitsGetCall) Fields(s ...googleapi.Field) *OrgunitsGetCall
func (c *OrgunitsGetCall) IfNoneMatch(entityTag string) *OrgunitsGetCall
```

## OrgunitsInsertCall

```go
func (c *OrgunitsInsertCall) Context(ctx context.Context) *OrgunitsInsertCall
func (c *OrgunitsInsertCall) Do(opts ...googleapi.CallOption) (*OrgUnit, error)
func (c *OrgunitsInsertCall) Fields(s ...googleapi.Field) *OrgunitsInsertCall
```

## OrgunitsListCall

```go
func (c *OrgunitsListCall) Context(ctx context.Context) *OrgunitsListCall
func (c *OrgunitsListCall) Do(opts ...googleapi.CallOption) (*OrgUnits, error)
func (c *OrgunitsListCall) Fields(s ...googleapi.Field) *OrgunitsListCall
func (c *OrgunitsListCall) IfNoneMatch(entityTag string) *OrgunitsListCall
func (c *OrgunitsListCall) OrgUnitPath(orgUnitPath string) *OrgunitsListCall  // Subtree root
func (c *OrgunitsListCall) Type(type_ string) *OrgunitsListCall               // "all" or "children"
```

### Type Values

| Value        | Returns                                            |
| ------------ | -------------------------------------------------- |
| `"all"`      | All org units below the specified path (recursive) |
| `"children"` | Only immediate children of the specified path      |

## OrgunitsUpdateCall / OrgunitsPatchCall

```go
func (c *OrgunitsUpdateCall) Context(ctx context.Context) *OrgunitsUpdateCall
func (c *OrgunitsUpdateCall) Do(opts ...googleapi.CallOption) (*OrgUnit, error)
func (c *OrgunitsUpdateCall) Fields(s ...googleapi.Field) *OrgunitsUpdateCall

func (c *OrgunitsPatchCall) Context(ctx context.Context) *OrgunitsPatchCall
func (c *OrgunitsPatchCall) Do(opts ...googleapi.CallOption) (*OrgUnit, error)
func (c *OrgunitsPatchCall) Fields(s ...googleapi.Field) *OrgunitsPatchCall
```

## OrgunitsDeleteCall

```go
func (c *OrgunitsDeleteCall) Context(ctx context.Context) *OrgunitsDeleteCall
func (c *OrgunitsDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *OrgunitsDeleteCall) Fields(s ...googleapi.Field) *OrgunitsDeleteCall
```

## orgUnitPath Parameter

The `orgUnitPath` parameter in Get/Update/Patch/Delete accepts:

- **By path**: Full path without leading `/`: `"corp/sales"` (for OrgUnitPath `/corp/sales`)
- **By ID**: `"id:03ph8a2z1enr9sn"` (prefix `"id:"` + OrgUnitId)

The path-based form requires URL-encoding of special characters in path segments.

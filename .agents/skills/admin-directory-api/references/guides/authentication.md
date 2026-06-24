# Authentication & Service Creation

## Creating a Service

```go
import admin "google.golang.org/api/admin/directory/v1"

svc, err := admin.NewService(ctx, option.WithHTTPClient(httpClient))
```

Or with Application Default Credentials:

```go
svc, err := admin.NewService(ctx)
```

## Restricting Scopes

By default, all available scopes are used. To restrict:

```go
svc, err := admin.NewService(ctx,
    option.WithScopes(
        admin.AdminDirectoryUserScope,
        admin.AdminDirectoryGroupScope,
    ),
)
```

## Service Structure

After creation, access sub-services via fields:

```go
svc.Users          // *UsersService
svc.Users.Aliases  // *UsersAliasesService (nested)
svc.Users.Photos   // *UsersPhotosService (nested)
svc.Groups         // *GroupsService
svc.Groups.Aliases // *GroupsAliasesService (nested)
svc.Members        // *MembersService
svc.Orgunits       // *OrgunitsService
svc.Roles          // *RolesService
svc.RoleAssignments // *RoleAssignmentsService
svc.Schemas        // *SchemasService
svc.Domains        // *DomainsService
svc.DomainAliases  // *DomainAliasesService
svc.Customers      // *CustomersService
svc.Chromeosdevices // *ChromeosdevicesService
svc.Mobiledevices  // *MobiledevicesService
svc.Privileges     // *PrivilegesService
```

## Scope Selection

| Scope Constant                                | Grants                            |
| --------------------------------------------- | --------------------------------- |
| `AdminDirectoryUserScope`                     | User read/write                   |
| `AdminDirectoryUserReadonlyScope`             | User read-only                    |
| `AdminDirectoryUserAliasScope`                | User alias read/write             |
| `AdminDirectoryUserAliasReadonlyScope`        | User alias read-only              |
| `AdminDirectoryUserSecurityScope`             | User security (ASPs, tokens, 2SV) |
| `AdminDirectoryGroupScope`                    | Group read/write                  |
| `AdminDirectoryGroupReadonlyScope`            | Group read-only                   |
| `AdminDirectoryGroupMemberScope`              | Group member read/write           |
| `AdminDirectoryGroupMemberReadonlyScope`      | Group member read-only            |
| `AdminDirectoryOrgunitScope`                  | Org unit read/write               |
| `AdminDirectoryOrgunitReadonlyScope`          | Org unit read-only                |
| `AdminDirectoryRolemanagementScope`           | Role management read/write        |
| `AdminDirectoryRolemanagementReadonlyScope`   | Role management read-only         |
| `AdminDirectoryCustomerScope`                 | Customer read/write               |
| `AdminDirectoryCustomerReadonlyScope`         | Customer read-only                |
| `AdminDirectoryDomainScope`                   | Domain read/write                 |
| `AdminDirectoryDomainReadonlyScope`           | Domain read-only                  |
| `AdminDirectoryDeviceChromeosScope`           | Chrome OS device read/write       |
| `AdminDirectoryDeviceChromeosReadonlyScope`   | Chrome OS device read-only        |
| `AdminDirectoryDeviceMobileScope`             | Mobile device read/write          |
| `AdminDirectoryDeviceMobileActionScope`       | Mobile device action              |
| `AdminDirectoryDeviceMobileReadonlyScope`     | Mobile device read-only           |
| `AdminDirectoryUserschemaScope`               | User schema read/write            |
| `AdminDirectoryUserschemaReadonlyScope`       | User schema read-only             |
| `AdminDirectoryResourceCalendarScope`         | Calendar resource read/write      |
| `AdminDirectoryResourceCalendarReadonlyScope` | Calendar resource read-only       |
| `AdminChromePrintersScope`                    | Chrome printers read/write        |
| `AdminChromePrintersReadonlyScope`            | Chrome printers read-only         |
| `CloudPlatformScope`                          | Full Cloud Platform access        |

## "my_customer" Convention

Many service methods require a `customerId` or `customer` parameter. The string `"my_customer"` is a valid alias that refers to the authenticated account's customer:

```go
svc.Orgunits.List("my_customer").Do()
svc.Roles.List("my_customer").Do()
svc.Schemas.List("my_customer").Do()
svc.Domains.List("my_customer").Do()
```

## userKey Addressing

Any `userKey` parameter accepts three forms:

- Primary email: `"user@example.com"`
- Alias email: `"alias@example.com"`
- Unique ID: `"118234567890123456789"` (from `User.Id`)

Same pattern applies to `groupKey` and `memberKey`.

---
name: admin-directory-api
description: >
  Admin Directory API v1 Go SDK reference (google.golang.org/api/admin/directory/v1).
  Use this skill whenever implementing or modifying code that calls the Admin
  Directory API: users, groups, members, org units, roles, role assignments,
  schemas, domains, or devices. Also use when the user asks about Directory API
  call patterns, struct fields, or ForceSendFields. Always use this skill
  alongside terraform-provider-dev when adding new Directory-based Terraform
  resources.
---

# Admin Directory API v1 (Go SDK)

Package: `google.golang.org/api/admin/directory/v1`

## Service Overview

Access services via `svc.Users`, `svc.Groups`, etc. after creating a service:

```go
svc, err := admin.NewService(ctx, option.WithHTTPClient(httpClient))
```

| Service                   | Handles                                         |
| ------------------------- | ----------------------------------------------- |
| `svc.Users`               | User CRUD, list, aliases, photos, admin ops     |
| `svc.Groups`              | Group CRUD, list, aliases                       |
| `svc.Members`             | Group membership CRUD, list, hasMember          |
| `svc.Orgunits`            | Organizational unit CRUD, list                  |
| `svc.Roles`               | Custom admin role CRUD, list                    |
| `svc.RoleAssignments`     | Role assignment CRUD, list                      |
| `svc.Schemas`             | Custom user schema CRUD, list                   |
| `svc.Domains`             | Domain CRUD, list                               |
| `svc.DomainAliases`       | Domain alias CRUD, list                         |
| `svc.Customers`           | Customer account get, patch, update             |
| `svc.Chromeosdevices`     | Chrome OS device management                     |
| `svc.Mobiledevices`       | Mobile device management                        |
| `svc.Privileges`          | Admin privilege listing                         |
| `svc.Resources`           | Calendar resources (buildings, rooms, features) |
| `svc.Tokens`              | OAuth token management                          |
| `svc.Asps`                | Application-specific passwords                  |
| `svc.TwoStepVerification` | 2SV management                                  |
| `svc.VerificationCodes`   | Backup verification codes                       |
| `svc.Channels`            | Stop push notification channels                 |

## Universal Call Pattern

Every operation is a builder chain:

```go
user, err := svc.Users.Get(userKey).
    Projection("full").
    Fields("primaryEmail,name,suspended").
    Context(ctx).
    Do()
```

## Critical Conventions

### ForceSendFields (boolean/zero-value pitfall)

All struct fields use `omitempty`. Boolean `false`, int `0`, and empty strings are silently dropped unless listed:

```go
user := &admin.User{
    Suspended:       false,
    ForceSendFields: []string{"Suspended"},
}
```

### "my_customer" Alias

Services requiring `customerId` (Orgunits, Roles, Schemas, Domains) accept `"my_customer"` as an alias for the authenticated account's customer ID.

### userKey / groupKey Flexibility

Any parameter named `userKey` or `groupKey` accepts: primary email, alias email, or unique ID.

### Projection("full") for Users

User reads default to "basic" projection (limited fields). Use `Projection("full")` to get all fields including custom schemas.

### Insert, Not Create

This API uses `Insert` for creation methods (not `Create` like Drive API).

## Routing Table

| Working on...                                 | Read                                  |
| --------------------------------------------- | ------------------------------------- |
| ForceSendFields, pagination, errors, Fields() | `references/guides/conventions.md`    |
| Service creation, scopes, auth                | `references/guides/authentication.md` |
| User CRUD, aliases, photos, admin ops         | `references/guides/users.md`          |
| Group CRUD, membership, roles                 | `references/guides/groups-members.md` |
| Org unit hierarchy, CRUD                      | `references/guides/org-units.md`      |
| Roles, role assignments, privileges           | `references/guides/roles.md`          |
| Domain and domain alias management            | `references/guides/domains.md`        |
| Chrome OS and mobile device management        | `references/guides/devices.md`        |
| Custom user schemas                           | `references/guides/schemas.md`        |

## Detailed API Reference

For complete struct fields and all call options:

| Need                                     | Read                                  |
| ---------------------------------------- | ------------------------------------- |
| User struct (50+ fields)                 | `references/api/user-struct.md`       |
| UsersService methods + call options      | `references/api/users-service.md`     |
| Group struct                             | `references/api/group-struct.md`      |
| GroupsService methods + call options     | `references/api/groups-service.md`    |
| Member struct + MembersService           | `references/api/members-service.md`   |
| OrgUnit struct                           | `references/api/orgunit-struct.md`    |
| OrgunitsService methods + call options   | `references/api/orgunits-service.md`  |
| Role + RoleAssignment + both services    | `references/api/roles-service.md`     |
| Schema + SchemaFieldSpec + service       | `references/api/schemas-service.md`   |
| Domain + DomainAlias + both services     | `references/api/domains-service.md`   |
| Customer struct + CustomersService       | `references/api/customers-service.md` |
| ChromeOsDevice + MobileDevice + services | `references/api/devices-service.md`   |
| OAuth2 scope constants                   | `references/api/scopes.md`            |

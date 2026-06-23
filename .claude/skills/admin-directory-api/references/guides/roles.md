# Roles & Role Assignments

## Listing Roles

```go
var allRoles []*admin.Role
err := svc.Roles.List("my_customer").
    Pages(ctx, func(resp *admin.Roles) error {
        allRoles = append(allRoles, resp.Items...)
        return nil
    })
```

Note: List response uses `Items` (not `Roles`).

## Getting a Role

```go
role, err := svc.Roles.Get("my_customer", roleIdStr).Context(ctx).Do()
// role.RoleId is int64
// role.RoleName, role.RoleDescription
// role.IsSystemRole, role.IsSuperAdminRole
// role.RolePrivileges
```

The `roleId` parameter is a string representation of the int64 RoleId:

```go
roleIdStr := strconv.FormatInt(role.RoleId, 10)
```

## Creating a Custom Role

```go
role := &admin.Role{
    RoleName:        "HelpDeskAdmin",
    RoleDescription: "Help desk support role",
    RolePrivileges: []*admin.RoleRolePrivileges{
        {PrivilegeName: "USERS_RETRIEVE", ServiceId: "00haapch16h1ysv"},
        {PrivilegeName: "USERS_UPDATE", ServiceId: "00haapch16h1ysv"},
    },
}
created, err := svc.Roles.Insert("my_customer", role).Context(ctx).Do()
```

## Updating a Role

```go
role := &admin.Role{
    RoleDescription: "Updated description",
    RolePrivileges: []*admin.RoleRolePrivileges{
        {PrivilegeName: "USERS_RETRIEVE", ServiceId: "00haapch16h1ysv"},
    },
}
updated, err := svc.Roles.Update("my_customer", roleIdStr, role).Context(ctx).Do()
```

Use `Patch` for partial updates (only modifies provided fields).

## Deleting a Role

```go
err := svc.Roles.Delete("my_customer", roleIdStr).Context(ctx).Do()
```

Only custom roles can be deleted. System roles return an error.

---

## Creating a Role Assignment

```go
assignment := &admin.RoleAssignment{
    RoleId:    roleId,          // int64
    AssignedTo: userId,         // user_id, group_id, or service account uniqueId
    ScopeType: "CUSTOMER",     // or "ORG_UNIT"
    // OrgUnitId: "03ph8a2z",  // required if ScopeType == "ORG_UNIT"
}
created, err := svc.RoleAssignments.Insert("my_customer", assignment).Context(ctx).Do()
```

## Getting a Role Assignment

```go
assignment, err := svc.RoleAssignments.Get("my_customer", assignmentIdStr).Context(ctx).Do()
```

## Listing Role Assignments

```go
// All assignments for a specific role
var assignments []*admin.RoleAssignment
err := svc.RoleAssignments.List("my_customer").
    RoleId(roleIdStr).
    Pages(ctx, func(resp *admin.RoleAssignments) error {
        assignments = append(assignments, resp.Items...)
        return nil
    })

// All assignments for a specific user
svc.RoleAssignments.List("my_customer").UserKey("user@example.com").Do()
```

## Deleting a Role Assignment

```go
err := svc.RoleAssignments.Delete("my_customer", assignmentIdStr).Context(ctx).Do()
```

---

## Listing Privileges

Get all available privileges (needed when building custom roles):

```go
privs, err := svc.Privileges.List("my_customer").Context(ctx).Do()
for _, p := range privs.Items {
    fmt.Printf("%s (service: %s)\n", p.PrivilegeName, p.ServiceId)
}
```

## int64 ID Conversion Pattern

RoleId and RoleAssignmentId are int64 in Go but passed as strings in API methods:

```go
import "strconv"

// int64 → string (for method parameters)
roleIdStr := strconv.FormatInt(role.RoleId, 10)
assignmentIdStr := strconv.FormatInt(assignment.RoleAssignmentId, 10)

// string → int64 (e.g., from Terraform state)
roleId, err := strconv.ParseInt(stateValue, 10, 64)
```

## Scope Types

| ScopeType    | OrgUnitId | Effect                                     |
| ------------ | --------- | ------------------------------------------ |
| `"CUSTOMER"` | empty     | Role applies across entire customer        |
| `"ORG_UNIT"` | required  | Role restricted to specified org unit tree |

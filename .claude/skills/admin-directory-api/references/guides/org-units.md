# Org Units

## Hierarchy Concept

Organizational units form a tree rooted at `/`. Users, devices, and policies are assigned to org units to control access and configuration. Maximum depth: 35 levels.

## Creating an Org Unit

```go
ou := &admin.OrgUnit{
    Name:              "Engineering",
    ParentOrgUnitPath: "/",
    Description:       "Engineering department",
}
created, err := svc.Orgunits.Insert("my_customer", ou).Context(ctx).Do()
// created.OrgUnitPath == "/Engineering"
// created.OrgUnitId == "03ph8a2z1enr9sn"
```

Required: `Name` and one of `ParentOrgUnitPath` or `ParentOrgUnitId`.

## Reading an Org Unit

By path (strip leading `/`):

```go
ou, err := svc.Orgunits.Get("my_customer", "Engineering").Context(ctx).Do()
// For nested: "Engineering/Backend"
```

By ID (prefix with `"id:"`):

```go
ou, err := svc.Orgunits.Get("my_customer", "id:03ph8a2z1enr9sn").Context(ctx).Do()
```

## Updating an Org Unit

```go
ou := &admin.OrgUnit{
    Description: "Updated description",
}
updated, err := svc.Orgunits.Update("my_customer", "id:03ph8a2z1enr9sn", ou).
    Context(ctx).Do()
```

To move an org unit (reparent):

```go
ou := &admin.OrgUnit{
    ParentOrgUnitPath: "/NewParent",
}
svc.Orgunits.Update("my_customer", "id:03ph8a2z1enr9sn", ou).Context(ctx).Do()
```

## Deleting an Org Unit

```go
err := svc.Orgunits.Delete("my_customer", "id:03ph8a2z1enr9sn").Context(ctx).Do()
```

An org unit must be empty (no users/devices assigned) before deletion.

## Listing Org Units

All org units in the account:

```go
result, err := svc.Orgunits.List("my_customer").Type("all").Context(ctx).Do()
for _, ou := range result.OrganizationUnits {
    fmt.Println(ou.OrgUnitPath)
}
```

Only immediate children of a path:

```go
result, err := svc.Orgunits.List("my_customer").
    OrgUnitPath("Engineering").
    Type("children").
    Context(ctx).Do()
```

Note: The List response is NOT paginated. All matching org units are returned in a single response.

## Path vs ID Addressing

| Form | Example                 | When to use                                   |
| ---- | ----------------------- | --------------------------------------------- |
| Path | `"Engineering/Backend"` | Human-readable; breaks if org unit is renamed |
| ID   | `"id:03ph8a2z1enr9sn"`  | Stable; survives renames and moves            |

The path form strips the leading `/` from `OrgUnitPath`. For the root org unit path `/Engineering`, pass `"Engineering"`.

## Type Filter

| Value        | Returns                                           |
| ------------ | ------------------------------------------------- |
| `"all"`      | All descendants recursively                       |
| `"children"` | Only immediate children of the OrgUnitPath filter |

If no `OrgUnitPath` is set, listing starts from root (`/`).

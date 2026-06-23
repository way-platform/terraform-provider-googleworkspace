# Groups & Members

## Creating a Group

```go
group := &admin.Group{
    Email:       "engineering@example.com",
    Name:        "Engineering",
    Description: "Engineering team distribution list",
}
created, err := svc.Groups.Insert(group).Context(ctx).Do()
```

Required field: `Email`.

## Reading a Group

```go
group, err := svc.Groups.Get("engineering@example.com").Context(ctx).Do()
```

The `groupKey` accepts group email, alias email, or unique group ID.

## Updating a Group

```go
group := &admin.Group{
    Name:        "Engineering Team",
    Description: "Updated description",
}
updated, err := svc.Groups.Update("engineering@example.com", group).Context(ctx).Do()
```

## Deleting a Group

```go
err := svc.Groups.Delete("engineering@example.com").Context(ctx).Do()
```

## Listing Groups

```go
var allGroups []*admin.Group
err := svc.Groups.List().
    Customer("my_customer").
    MaxResults(200).
    Pages(ctx, func(resp *admin.Groups) error {
        allGroups = append(allGroups, resp.Groups...)
        return nil
    })
```

Filter groups for a specific user:

```go
svc.Groups.List().UserKey("user@example.com").Do()
```

## Group Aliases

```go
// Add alias
alias := &admin.Alias{Alias: "eng@example.com"}
_, err := svc.Groups.Aliases.Insert("engineering@example.com", alias).Context(ctx).Do()

// List aliases
aliases, err := svc.Groups.Aliases.List("engineering@example.com").Context(ctx).Do()

// Delete alias
err := svc.Groups.Aliases.Delete("engineering@example.com", "eng@example.com").Context(ctx).Do()
```

---

## Adding a Member

```go
member := &admin.Member{
    Email: "user@example.com",
    Role:  "MEMBER",  // "OWNER", "MANAGER", or "MEMBER"
}
created, err := svc.Members.Insert("engineering@example.com", member).Context(ctx).Do()
```

Members can be users or groups (nested membership).

## Reading a Member

```go
member, err := svc.Members.Get("engineering@example.com", "user@example.com").Context(ctx).Do()
```

## Updating a Member

```go
member := &admin.Member{
    Role: "MANAGER",
}
updated, err := svc.Members.Update("engineering@example.com", "user@example.com", member).
    Context(ctx).Do()
```

## Removing a Member

```go
err := svc.Members.Delete("engineering@example.com", "user@example.com").Context(ctx).Do()
```

## Listing Members

```go
var allMembers []*admin.Member
err := svc.Members.List("engineering@example.com").
    MaxResults(200).
    Pages(ctx, func(resp *admin.Members) error {
        allMembers = append(allMembers, resp.Members...)
        return nil
    })
```

Filter by role:

```go
svc.Members.List(groupKey).Roles("OWNER").Do()
svc.Members.List(groupKey).Roles("OWNER,MANAGER").Do()
```

Include indirect (nested group) members:

```go
svc.Members.List(groupKey).IncludeDerivedMembership(true).Do()
```

## Checking Membership

```go
result, err := svc.Members.HasMember("engineering@example.com", "user@example.com").
    Context(ctx).Do()
if result.IsMember {
    // user is a member
}
```

Returns 200 with `IsMember: true/false`. Does not check nested group membership by default.

## Member Roles

| Role        | Can manage members | Can manage settings | Full control |
| ----------- | ------------------ | ------------------- | ------------ |
| `"OWNER"`   | Yes                | Yes                 | Yes          |
| `"MANAGER"` | Yes                | Yes                 | No           |
| `"MEMBER"`  | No                 | No                  | No           |

## Member Types

| Type         | Description                      |
| ------------ | -------------------------------- |
| `"USER"`     | Individual user account          |
| `"GROUP"`    | Nested group (group-in-group)    |
| `"EXTERNAL"` | External user outside the domain |

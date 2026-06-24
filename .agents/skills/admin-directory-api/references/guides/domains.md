# Domains & Domain Aliases

## Listing Domains

```go
result, err := svc.Domains.List("my_customer").Context(ctx).Do()
for _, d := range result.Domains {
    fmt.Printf("%s (primary: %v, verified: %v)\n", d.DomainName, d.IsPrimary, d.Verified)
}
```

Not paginated; all domains returned in one response.

## Getting a Domain

```go
domain, err := svc.Domains.Get("my_customer", "example.com").Context(ctx).Do()
// domain.DomainName, domain.IsPrimary, domain.Verified
// domain.DomainAliases (embedded aliases for this domain)
```

## Inserting a Domain

```go
domain := &admin.Domains{
    DomainName: "newdomain.com",
}
created, err := svc.Domains.Insert("my_customer", domain).Context(ctx).Do()
```

New domains are unverified. DNS verification must be completed separately.

## Deleting a Domain

```go
err := svc.Domains.Delete("my_customer", "newdomain.com").Context(ctx).Do()
```

Cannot delete the primary domain.

---

## Listing Domain Aliases

```go
result, err := svc.DomainAliases.List("my_customer").Context(ctx).Do()
for _, a := range result.DomainAliases {
    fmt.Printf("%s -> %s (verified: %v)\n", a.DomainAliasName, a.ParentDomainName, a.Verified)
}
```

Filter by parent domain:

```go
svc.DomainAliases.List("my_customer").ParentDomainName("example.com").Do()
```

## Getting a Domain Alias

```go
alias, err := svc.DomainAliases.Get("my_customer", "alias.example.com").Context(ctx).Do()
```

## Inserting a Domain Alias

```go
alias := &admin.DomainAlias{
    DomainAliasName:  "alias.example.com",
    ParentDomainName: "example.com",
}
created, err := svc.DomainAliases.Insert("my_customer", alias).Context(ctx).Do()
```

## Deleting a Domain Alias

```go
err := svc.DomainAliases.Delete("my_customer", "alias.example.com").Context(ctx).Do()
```

## Primary vs Secondary Domains

| Type      | `IsPrimary` | Notes                               |
| --------- | ----------- | ----------------------------------- |
| Primary   | `true`      | One per customer; cannot be deleted |
| Secondary | `false`     | Additional domains; can be deleted  |

Domain aliases are aliases for either primary or secondary domains. They share the same user namespace.

## Verification

- `Verified: true` means DNS verification is complete
- New domains and aliases start as unverified
- Verification is performed outside the API (DNS TXT/CNAME records)
- Unverified domains cannot be used for user email addresses

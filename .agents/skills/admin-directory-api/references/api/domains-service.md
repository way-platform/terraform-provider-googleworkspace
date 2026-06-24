# Domain & DomainAlias Structs + Services

## Domain Struct

```go
type Domains struct {
    CreationTime int64            `json:"creationTime,omitempty,string"` // Unix millis
    DomainAliases []*DomainAlias  `json:"domainAliases,omitempty"`
    DomainName   string           `json:"domainName,omitempty"`
    Etag         string           `json:"etag,omitempty"`
    IsPrimary    bool             `json:"isPrimary,omitempty"`
    Kind         string           `json:"kind,omitempty"`
    Verified     bool             `json:"verified,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

Note: The type is named `Domains` (plural) for the individual domain resource in this SDK. This is a quirk of the generated code.

## DomainAlias Struct

```go
type DomainAlias struct {
    CreationTime     int64  `json:"creationTime,omitempty,string"` // Unix millis
    DomainAliasName  string `json:"domainAliasName,omitempty"`
    Etag             string `json:"etag,omitempty"`
    Kind             string `json:"kind,omitempty"`
    ParentDomainName string `json:"parentDomainName,omitempty"` // Associated primary/secondary domain
    Verified         bool   `json:"verified,omitempty"`         // Read-only

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## Domains2 (List Response)

```go
type Domains2 struct {
    Domains []*Domains `json:"domains,omitempty"`
    Etag    string     `json:"etag,omitempty"`
    Kind    string     `json:"kind,omitempty"`
}
```

## DomainAliases (List Response)

```go
type DomainAliases struct {
    DomainAliases []*DomainAlias `json:"domainAliases,omitempty"`
    Etag          string         `json:"etag,omitempty"`
    Kind          string         `json:"kind,omitempty"`
}
```

Neither list response uses pagination.

---

## DomainsService

All methods require `customer` as the first parameter.

```go
type DomainsService struct{}

func NewDomainsService(s *Service) *DomainsService
func (r *DomainsService) Delete(customer string, domainName string) *DomainsDeleteCall
func (r *DomainsService) Get(customer string, domainName string) *DomainsGetCall
func (r *DomainsService) Insert(customer string, domains *Domains) *DomainsInsertCall
func (r *DomainsService) List(customer string) *DomainsListCall
```

## DomainsGetCall

```go
func (c *DomainsGetCall) Context(ctx context.Context) *DomainsGetCall
func (c *DomainsGetCall) Do(opts ...googleapi.CallOption) (*Domains, error)
func (c *DomainsGetCall) Fields(s ...googleapi.Field) *DomainsGetCall
func (c *DomainsGetCall) IfNoneMatch(entityTag string) *DomainsGetCall
```

## DomainsInsertCall

```go
func (c *DomainsInsertCall) Context(ctx context.Context) *DomainsInsertCall
func (c *DomainsInsertCall) Do(opts ...googleapi.CallOption) (*Domains, error)
func (c *DomainsInsertCall) Fields(s ...googleapi.Field) *DomainsInsertCall
```

## DomainsListCall

```go
func (c *DomainsListCall) Context(ctx context.Context) *DomainsListCall
func (c *DomainsListCall) Do(opts ...googleapi.CallOption) (*Domains2, error)
func (c *DomainsListCall) Fields(s ...googleapi.Field) *DomainsListCall
func (c *DomainsListCall) IfNoneMatch(entityTag string) *DomainsListCall
```

## DomainsDeleteCall

```go
func (c *DomainsDeleteCall) Context(ctx context.Context) *DomainsDeleteCall
func (c *DomainsDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *DomainsDeleteCall) Fields(s ...googleapi.Field) *DomainsDeleteCall
```

---

## DomainAliasesService

```go
type DomainAliasesService struct{}

func NewDomainAliasesService(s *Service) *DomainAliasesService
func (r *DomainAliasesService) Delete(customer string, domainAliasName string) *DomainAliasesDeleteCall
func (r *DomainAliasesService) Get(customer string, domainAliasName string) *DomainAliasesGetCall
func (r *DomainAliasesService) Insert(customer string, domainalias *DomainAlias) *DomainAliasesInsertCall
func (r *DomainAliasesService) List(customer string) *DomainAliasesListCall
```

## DomainAliasesGetCall

```go
func (c *DomainAliasesGetCall) Context(ctx context.Context) *DomainAliasesGetCall
func (c *DomainAliasesGetCall) Do(opts ...googleapi.CallOption) (*DomainAlias, error)
func (c *DomainAliasesGetCall) Fields(s ...googleapi.Field) *DomainAliasesGetCall
func (c *DomainAliasesGetCall) IfNoneMatch(entityTag string) *DomainAliasesGetCall
```

## DomainAliasesInsertCall

```go
func (c *DomainAliasesInsertCall) Context(ctx context.Context) *DomainAliasesInsertCall
func (c *DomainAliasesInsertCall) Do(opts ...googleapi.CallOption) (*DomainAlias, error)
func (c *DomainAliasesInsertCall) Fields(s ...googleapi.Field) *DomainAliasesInsertCall
```

## DomainAliasesListCall

```go
func (c *DomainAliasesListCall) Context(ctx context.Context) *DomainAliasesListCall
func (c *DomainAliasesListCall) Do(opts ...googleapi.CallOption) (*DomainAliases, error)
func (c *DomainAliasesListCall) Fields(s ...googleapi.Field) *DomainAliasesListCall
func (c *DomainAliasesListCall) IfNoneMatch(entityTag string) *DomainAliasesListCall
func (c *DomainAliasesListCall) ParentDomainName(parentDomainName string) *DomainAliasesListCall
```

## DomainAliasesDeleteCall

```go
func (c *DomainAliasesDeleteCall) Context(ctx context.Context) *DomainAliasesDeleteCall
func (c *DomainAliasesDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *DomainAliasesDeleteCall) Fields(s ...googleapi.Field) *DomainAliasesDeleteCall
```

# Customer Struct & CustomersService

## Customer Struct

```go
type Customer struct {
    AlternateEmail       string                 `json:"alternateEmail,omitempty"`       // Secondary contact email (cannot be on customerDomain)
    CustomerCreationTime string                 `json:"customerCreationTime,omitempty"` // Read-only
    CustomerDomain       string                 `json:"customerDomain,omitempty"`       // Primary domain (no "www" prefix)
    Etag                 string                 `json:"etag,omitempty"`
    Id                   string                 `json:"id,omitempty"`                   // Read-only unique ID
    Kind                 string                 `json:"kind,omitempty"`
    Language             string                 `json:"language,omitempty"`             // ISO 639-2; default "en"
    PhoneNumber          string                 `json:"phoneNumber,omitempty"`          // E.164 format
    PostalAddress        *CustomerPostalAddress `json:"postalAddress,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## CustomerPostalAddress Struct

```go
type CustomerPostalAddress struct {
    AddressLine1     string `json:"addressLine1,omitempty"`
    AddressLine2     string `json:"addressLine2,omitempty"`
    AddressLine3     string `json:"addressLine3,omitempty"`
    ContactName      string `json:"contactName,omitempty"`
    CountryCode      string `json:"countryCode,omitempty"`      // ISO 3166; required
    Locality         string `json:"locality,omitempty"`         // City
    OrganizationName string `json:"organizationName,omitempty"`
    PostalCode       string `json:"postalCode,omitempty"`
    Region           string `json:"region,omitempty"`           // State/province (e.g. "NY")

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

---

## CustomersService

```go
type CustomersService struct {
    Chrome *CustomersChromeService
}

func NewCustomersService(s *Service) *CustomersService
func (r *CustomersService) Get(customerKey string) *CustomersGetCall
func (r *CustomersService) Patch(customerKey string, customer *Customer) *CustomersPatchCall
func (r *CustomersService) Update(customerKey string, customer *Customer) *CustomersUpdateCall
```

The `customerKey` parameter accepts:

- `"my_customer"` alias
- The customer's primary domain
- The unique customer ID

Note: No Insert or Delete methods. Customers are created/deleted through other means.

## CustomersGetCall

```go
func (c *CustomersGetCall) Context(ctx context.Context) *CustomersGetCall
func (c *CustomersGetCall) Do(opts ...googleapi.CallOption) (*Customer, error)
func (c *CustomersGetCall) Fields(s ...googleapi.Field) *CustomersGetCall
func (c *CustomersGetCall) IfNoneMatch(entityTag string) *CustomersGetCall
```

## CustomersPatchCall

```go
func (c *CustomersPatchCall) Context(ctx context.Context) *CustomersPatchCall
func (c *CustomersPatchCall) Do(opts ...googleapi.CallOption) (*Customer, error)
func (c *CustomersPatchCall) Fields(s ...googleapi.Field) *CustomersPatchCall
```

## CustomersUpdateCall

```go
func (c *CustomersUpdateCall) Context(ctx context.Context) *CustomersUpdateCall
func (c *CustomersUpdateCall) Do(opts ...googleapi.CallOption) (*Customer, error)
func (c *CustomersUpdateCall) Fields(s ...googleapi.Field) *CustomersUpdateCall
```

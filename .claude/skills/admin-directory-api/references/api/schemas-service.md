# Schema & SchemaFieldSpec Structs + SchemasService

## Schema Struct

```go
type Schema struct {
    DisplayName string             `json:"displayName,omitempty"`
    Etag        string             `json:"etag,omitempty"`
    Fields      []*SchemaFieldSpec `json:"fields,omitempty"`
    Kind        string             `json:"kind,omitempty"`
    SchemaId    string             `json:"schemaId,omitempty"`   // Read-only unique ID
    SchemaName  string             `json:"schemaName,omitempty"` // Must be unique per customer

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## SchemaFieldSpec Struct

```go
type SchemaFieldSpec struct {
    DisplayName        string                            `json:"displayName,omitempty"`
    Etag               string                            `json:"etag,omitempty"`
    FieldId            string                            `json:"fieldId,omitempty"`    // Read-only
    FieldName          string                            `json:"fieldName,omitempty"`
    FieldType          string                            `json:"fieldType,omitempty"`  // See table below
    Indexed            *bool                             `json:"indexed,omitempty"`    // Default: true
    Kind               string                            `json:"kind,omitempty"`
    MultiValued        bool                              `json:"multiValued,omitempty"` // Default: false
    NumericIndexingSpec *SchemaFieldSpecNumericIndexingSpec `json:"numericIndexingSpec,omitempty"`
    ReadAccessType     string                            `json:"readAccessType,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

### FieldType Values

| Value      | Go equivalent | Notes                   |
| ---------- | ------------- | ----------------------- |
| `"STRING"` | string        | General text            |
| `"INT64"`  | int64         | Integer value           |
| `"BOOL"`   | bool          | Boolean value           |
| `"DOUBLE"` | float64       | Floating point          |
| `"EMAIL"`  | string        | Email address           |
| `"PHONE"`  | string        | Phone number            |
| `"DATE"`   | string        | Date in ISO 8601 format |

### ReadAccessType Values

| Value                | Visibility                                     |
| -------------------- | ---------------------------------------------- |
| `"ALL_DOMAIN_USERS"` | Visible to all users in the domain             |
| `"ADMINS_AND_SELF"`  | Visible only to admins and the user themselves |

## SchemaFieldSpecNumericIndexingSpec

```go
type SchemaFieldSpecNumericIndexingSpec struct {
    MaxValue double `json:"maxValue,omitempty"`
    MinValue double `json:"minValue,omitempty"`

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## Schemas (List Response)

```go
type Schemas struct {
    Schemas []*Schema `json:"schemas,omitempty"`
    Etag    string    `json:"etag,omitempty"`
    Kind    string    `json:"kind,omitempty"`
}
```

Note: Schemas list is NOT paginated. All schemas returned in a single response.

---

## SchemasService

All methods require `customerId` as the first parameter.

```go
type SchemasService struct{}

func NewSchemasService(s *Service) *SchemasService
func (r *SchemasService) Delete(customerId string, schemaKey string) *SchemasDeleteCall
func (r *SchemasService) Get(customerId string, schemaKey string) *SchemasGetCall
func (r *SchemasService) Insert(customerId string, schema *Schema) *SchemasInsertCall
func (r *SchemasService) List(customerId string) *SchemasListCall
func (r *SchemasService) Patch(customerId string, schemaKey string, schema *Schema) *SchemasPatchCall
func (r *SchemasService) Update(customerId string, schemaKey string, schema *Schema) *SchemasUpdateCall
```

The `schemaKey` parameter accepts either `SchemaName` or `SchemaId`.

## SchemasListCall

```go
func (c *SchemasListCall) Context(ctx context.Context) *SchemasListCall
func (c *SchemasListCall) Do(opts ...googleapi.CallOption) (*Schemas, error)
func (c *SchemasListCall) Fields(s ...googleapi.Field) *SchemasListCall
func (c *SchemasListCall) IfNoneMatch(entityTag string) *SchemasListCall
```

## SchemasGetCall

```go
func (c *SchemasGetCall) Context(ctx context.Context) *SchemasGetCall
func (c *SchemasGetCall) Do(opts ...googleapi.CallOption) (*Schema, error)
func (c *SchemasGetCall) Fields(s ...googleapi.Field) *SchemasGetCall
func (c *SchemasGetCall) IfNoneMatch(entityTag string) *SchemasGetCall
```

## SchemasInsertCall

```go
func (c *SchemasInsertCall) Context(ctx context.Context) *SchemasInsertCall
func (c *SchemasInsertCall) Do(opts ...googleapi.CallOption) (*Schema, error)
func (c *SchemasInsertCall) Fields(s ...googleapi.Field) *SchemasInsertCall
```

## SchemasUpdateCall / SchemasPatchCall

```go
func (c *SchemasUpdateCall) Context(ctx context.Context) *SchemasUpdateCall
func (c *SchemasUpdateCall) Do(opts ...googleapi.CallOption) (*Schema, error)
func (c *SchemasUpdateCall) Fields(s ...googleapi.Field) *SchemasUpdateCall

func (c *SchemasPatchCall) Context(ctx context.Context) *SchemasPatchCall
func (c *SchemasPatchCall) Do(opts ...googleapi.CallOption) (*Schema, error)
func (c *SchemasPatchCall) Fields(s ...googleapi.Field) *SchemasPatchCall
```

## SchemasDeleteCall

```go
func (c *SchemasDeleteCall) Context(ctx context.Context) *SchemasDeleteCall
func (c *SchemasDeleteCall) Do(opts ...googleapi.CallOption) error
func (c *SchemasDeleteCall) Fields(s ...googleapi.Field) *SchemasDeleteCall
```

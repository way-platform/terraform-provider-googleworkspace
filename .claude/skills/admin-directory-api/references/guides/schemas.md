# Custom User Schemas

## Concept

Schemas define custom fields that can be attached to user profiles. Each schema has a unique name and contains one or more field specs. Schema data on users lives in `User.CustomSchemas`.

## Creating a Schema

```go
schema := &admin.Schema{
    SchemaName:  "EmployeeInfo",
    DisplayName: "Employee Information",
    Fields: []*admin.SchemaFieldSpec{
        {
            FieldName:   "department",
            DisplayName: "Department",
            FieldType:   "STRING",
        },
        {
            FieldName:   "startDate",
            DisplayName: "Start Date",
            FieldType:   "DATE",
        },
        {
            FieldName:      "skills",
            DisplayName:    "Skills",
            FieldType:      "STRING",
            MultiValued:    true,
            ReadAccessType: "ALL_DOMAIN_USERS",
        },
    },
}
created, err := svc.Schemas.Insert("my_customer", schema).Context(ctx).Do()
```

## Reading a Schema

```go
schema, err := svc.Schemas.Get("my_customer", "EmployeeInfo").Context(ctx).Do()
// or by ID:
schema, err := svc.Schemas.Get("my_customer", schemaId).Context(ctx).Do()
```

The `schemaKey` parameter accepts either `SchemaName` or `SchemaId`.

## Updating a Schema

```go
schema := &admin.Schema{
    Fields: []*admin.SchemaFieldSpec{
        {
            FieldName:   "department",
            DisplayName: "Department",
            FieldType:   "STRING",
        },
        {
            FieldName:   "level",
            DisplayName: "Level",
            FieldType:   "INT64",
        },
    },
}
updated, err := svc.Schemas.Update("my_customer", "EmployeeInfo", schema).Context(ctx).Do()
```

Note: Update replaces all fields. Use Patch for partial updates.

## Deleting a Schema

```go
err := svc.Schemas.Delete("my_customer", "EmployeeInfo").Context(ctx).Do()
```

## Listing All Schemas

```go
result, err := svc.Schemas.List("my_customer").Context(ctx).Do()
for _, schema := range result.Schemas {
    fmt.Printf("%s (%d fields)\n", schema.SchemaName, len(schema.Fields))
}
```

Not paginated; all schemas returned in one response.

## Relation to User.CustomSchemas

Schema data on users is accessed via `User.CustomSchemas`:

```go
// Reading custom schema values from a user
user, err := svc.Users.Get(userKey).Projection("full").Do()
if raw, ok := user.CustomSchemas["EmployeeInfo"]; ok {
    var data map[string]interface{}
    json.Unmarshal(raw, &data)
    // data["department"] == "Engineering"
}

// Writing custom schema values on a user
import "encoding/json"

schemaData, _ := json.Marshal(map[string]interface{}{
    "department": "Engineering",
    "startDate":  "2024-01-15",
    "skills": []map[string]interface{}{
        {"value": "Go"},
        {"value": "Terraform"},
    },
})
user := &admin.User{
    CustomSchemas: map[string]googleapi.RawMessage{
        "EmployeeInfo": schemaData,
    },
}
svc.Users.Patch(userKey, user).Context(ctx).Do()
```

## Multi-Valued Fields

Multi-valued fields are represented as arrays of objects with a `value` key:

```go
// Single-valued: {"fieldName": "value"}
// Multi-valued:  {"fieldName": [{"value": "v1"}, {"value": "v2"}]}
```

## Field Type Reference

| FieldType  | JSON representation | Notes             |
| ---------- | ------------------- | ----------------- |
| `"STRING"` | `"text"`            | General text      |
| `"INT64"`  | `"123"`             | Integer as string |
| `"BOOL"`   | `"true"`            | Boolean as string |
| `"DOUBLE"` | `"1.5"`             | Float as string   |
| `"EMAIL"`  | `"a@b.com"`         | Email address     |
| `"PHONE"`  | `"+1234"`           | Phone number      |
| `"DATE"`   | `"2024-01-15"`      | ISO 8601 date     |

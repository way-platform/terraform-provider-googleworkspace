# OrgUnit Struct

```go
type OrgUnit struct {
    // --- Writable fields ---

    Name              string `json:"name,omitempty"`              // Required; path segment name (e.g. "sales_support")
    Description       string `json:"description,omitempty"`
    ParentOrgUnitId   string `json:"parentOrgUnitId,omitempty"`   // One of parent fields required for creation
    ParentOrgUnitPath string `json:"parentOrgUnitPath,omitempty"` // One of parent fields required for creation
    BlockInheritance  bool   `json:"blockInheritance,omitempty"`  // Deprecated; default false

    // --- Read-only fields ---

    OrgUnitId   string `json:"orgUnitId,omitempty"`   // Unique ID (without "id:" prefix)
    OrgUnitPath string `json:"orgUnitPath,omitempty"` // Full path (e.g. "/corp/sales/sales_support")
    Etag        string `json:"etag,omitempty"`
    Kind        string `json:"kind,omitempty"`

    // --- JSON control ---
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## OrgUnits (List Response)

```go
type OrgUnits struct {
    OrganizationUnits []*OrgUnit `json:"organizationUnits,omitempty"`
    Etag              string     `json:"etag,omitempty"`
    Kind              string     `json:"kind,omitempty"`
}
```

Note: OrgUnits list response does NOT use pagination (no NextPageToken). All matching units are returned in a single response.

## Field Behavior Notes

| Field               | Behavior                                                       |
| ------------------- | -------------------------------------------------------------- |
| `Name`              | Required on Insert; becomes the last segment of OrgUnitPath    |
| `ParentOrgUnitPath` | One of ParentOrgUnitPath or ParentOrgUnitId required on Insert |
| `OrgUnitPath`       | Read-only; full path derived from parent path + name           |
| `OrgUnitId`         | Stable unique ID; use with "id:" prefix for Get/Update/Delete  |
| `BlockInheritance`  | Deprecated; setting to true is not recommended                 |

## Path vs ID Addressing

The `orgUnitPath` parameter in service methods accepts two forms:

- **By path**: The full path without leading slash segment: `"corp/sales/sales_support"` (note: strip the leading `/` from OrgUnitPath)
- **By ID**: The prefix `"id:"` followed by the OrgUnitId: `"id:03ph8a2z1enr9sn"`

## Hierarchy Constraints

- Maximum 35 levels of depth
- Full path is derived from `ParentOrgUnitPath + "/" + Name`
- Root org unit path is `/`
- Deleting an org unit moves its users to the parent org unit

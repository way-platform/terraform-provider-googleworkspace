# terraform-provider-googleworkspace

Always use `/claude-md-improver` when updating this file.

Terraform provider for managing Google Workspace resources via the Admin Directory, Drive, and Groups Settings APIs.

- **Module**: `github.com/way-platform/terraform-provider-googleworkspace`
- **Package**: `internal/provider/` (single flat package, all resources here)

## Commands

```bash
mise run test          # go test -count=1 -cover ./...
mise run lint          # golangci-lint run --fix ./...
mise run build         # full CI pipeline (download, tidy, lint, test, diff)
```

Single test:

```bash
go test ./internal/provider/ -v -run TestAccUser
```

Environment variables for acceptance tests:

- `SUBJECT` — impersonated user email
- `GOOGLEWORKSPACE_CUSTOMER_ID` — Workspace customer ID

## Skills

Always load `terraform-provider-dev`. Load the relevant API skill alongside it when working on a resource.

| Files                                                                                                                                         | Skills                   |
| --------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| Any resource or data source                                                                                                                   | `terraform-provider-dev` |
| resource_user, resource_group, resource_group_members, resource_group_settings, resource_org_unit, resource_role_assignment, data_source_role | + `admin-directory-api`  |
| resource_drive, resource_drive_permission                                                                                                     | + `google-drive-api`     |

## Architecture

- File naming: `resource_<name>.go`, `resource_<name>_test.go`, `data_source_<name>.go`
- Provider client: `*apiClient` struct (retryable HTTP client + customerID + basePath overrides)
- Auth: service account with Domain-Wide Delegation, impersonating a target user
- New resources must be registered in `provider.go` `Resources()` / `DataSources()`

## Conventions

### Google API zero-value fields

Always set `ForceSendFields` when a Google API struct field can legitimately be `false`, `0`, or `""`. Without it, Go's `omitempty` silently drops the field from the JSON request.

### 404 handling

- **Read**: call `resp.State.RemoveResource(ctx)` and return (resource was deleted externally)
- **Delete**: return without error (idempotent)

### Partial response Fields() selector

Every Google API call that uses `.Fields(...)` only returns the listed fields.
When adding a new attribute to a resource:

1. Add the field name to **every** `.Fields()` call in Create, Read, and Update
2. In tests, the mock handler must **only return fields present in the
   request's `fields` query parameter** (parse `r.URL.Query().Get("fields")`).
   Never return fields unconditionally; this catches a missing `Fields()`
   entry at test time rather than in production.

This prevents the "inconsistent result after apply" class of bugs where
Terraform plans a value but the Read returns null because the field wasn't
requested.

### Testing

Tests use a mock HTTP server, never real API calls:

- `setupTestServer()` — creates `httptest.Server` with a route handler
- `setupTestClient()` — injects mock client into provider, bypasses auth
- `jsonResponse()` — helper to write JSON responses in handlers
- Mock handlers must respect the `fields` query parameter (see above)

### Retry

Automatic retry on 429, 403 quota errors, and 5xx (except 501). Configurable via provider `retry_on` attribute.

## Benchmarking

When designing resources or solving implementation questions, reference these providers for patterns and prior art:

- [`hashicorp/terraform-provider-googleworkspace`](https://github.com/hashicorp/terraform-provider-googleworkspace) — HashiCorp's archived provider (SDK v2); good reference for schema design and API coverage
- [`hanneshayashi/terraform-provider-gdrive`](https://github.com/hanneshayashi/terraform-provider-gdrive) — Community Drive provider; good reference for Drive API patterns and permission handling

Clone or browse these repos to compare schema choices, CRUD implementations, and error handling before building new resources. Use them for inspiration and reference only; do not copy code verbatim and respect each project's license.

## CI

- PR: lint + test + build (GitHub Actions)
- Release: semantic-release + goreleaser on push to main

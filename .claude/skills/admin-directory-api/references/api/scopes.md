# OAuth2 Scope Constants

```go
package admin // google.golang.org/api/admin/directory/v1

const (
    // User management
    AdminDirectoryUserScope         = "https://www.googleapis.com/auth/admin.directory.user"
    AdminDirectoryUserReadonlyScope = "https://www.googleapis.com/auth/admin.directory.user.readonly"

    // User aliases
    AdminDirectoryUserAliasScope         = "https://www.googleapis.com/auth/admin.directory.user.alias"
    AdminDirectoryUserAliasReadonlyScope = "https://www.googleapis.com/auth/admin.directory.user.alias.readonly"

    // User security (ASPs, tokens, 2SV)
    AdminDirectoryUserSecurityScope = "https://www.googleapis.com/auth/admin.directory.user.security"

    // Group management
    AdminDirectoryGroupScope         = "https://www.googleapis.com/auth/admin.directory.group"
    AdminDirectoryGroupReadonlyScope = "https://www.googleapis.com/auth/admin.directory.group.readonly"

    // Group membership
    AdminDirectoryGroupMemberScope         = "https://www.googleapis.com/auth/admin.directory.group.member"
    AdminDirectoryGroupMemberReadonlyScope = "https://www.googleapis.com/auth/admin.directory.group.member.readonly"

    // Organizational units
    AdminDirectoryOrgunitScope         = "https://www.googleapis.com/auth/admin.directory.orgunit"
    AdminDirectoryOrgunitReadonlyScope = "https://www.googleapis.com/auth/admin.directory.orgunit.readonly"

    // Role management
    AdminDirectoryRolemanagementScope         = "https://www.googleapis.com/auth/admin.directory.rolemanagement"
    AdminDirectoryRolemanagementReadonlyScope = "https://www.googleapis.com/auth/admin.directory.rolemanagement.readonly"

    // Customer management
    AdminDirectoryCustomerScope         = "https://www.googleapis.com/auth/admin.directory.customer"
    AdminDirectoryCustomerReadonlyScope = "https://www.googleapis.com/auth/admin.directory.customer.readonly"

    // Domain management
    AdminDirectoryDomainScope         = "https://www.googleapis.com/auth/admin.directory.domain"
    AdminDirectoryDomainReadonlyScope = "https://www.googleapis.com/auth/admin.directory.domain.readonly"

    // Chrome OS devices
    AdminDirectoryDeviceChromeosScope         = "https://www.googleapis.com/auth/admin.directory.device.chromeos"
    AdminDirectoryDeviceChromeosReadonlyScope = "https://www.googleapis.com/auth/admin.directory.device.chromeos.readonly"

    // Mobile devices
    AdminDirectoryDeviceMobileScope         = "https://www.googleapis.com/auth/admin.directory.device.mobile"
    AdminDirectoryDeviceMobileActionScope   = "https://www.googleapis.com/auth/admin.directory.device.mobile.action"
    AdminDirectoryDeviceMobileReadonlyScope = "https://www.googleapis.com/auth/admin.directory.device.mobile.readonly"

    // User custom schemas
    AdminDirectoryUserschemaScope         = "https://www.googleapis.com/auth/admin.directory.userschema"
    AdminDirectoryUserschemaReadonlyScope = "https://www.googleapis.com/auth/admin.directory.userschema.readonly"

    // Calendar resources (buildings, rooms, features)
    AdminDirectoryResourceCalendarScope         = "https://www.googleapis.com/auth/admin.directory.resource.calendar"
    AdminDirectoryResourceCalendarReadonlyScope = "https://www.googleapis.com/auth/admin.directory.resource.calendar.readonly"

    // Chrome printers
    AdminChromePrintersScope         = "https://www.googleapis.com/auth/admin.chrome.printers"
    AdminChromePrintersReadonlyScope = "https://www.googleapis.com/auth/admin.chrome.printers.readonly"

    // Full Cloud Platform access
    CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"
)
```

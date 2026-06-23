# Device Structs & Services

## ChromeOsDevice Struct (Key Fields)

```go
type ChromeOsDevice struct {
    AnnotatedAssetId  string `json:"annotatedAssetId,omitempty"`
    AnnotatedLocation string `json:"annotatedLocation,omitempty"`
    AnnotatedUser     string `json:"annotatedUser,omitempty"`
    BootMode          string `json:"bootMode,omitempty"`          // "Verified", "Dev"
    DeviceId          string `json:"deviceId,omitempty"`          // Unique ID
    Etag              string `json:"etag,omitempty"`
    FirmwareVersion   string `json:"firmwareVersion,omitempty"`
    Kind              string `json:"kind,omitempty"`
    LastSync          string `json:"lastSync,omitempty"`          // RFC 3339
    MacAddress        string `json:"macAddress,omitempty"`
    Meid              string `json:"meid,omitempty"`
    Model             string `json:"model,omitempty"`
    Notes             string `json:"notes,omitempty"`
    OrgUnitPath       string `json:"orgUnitPath,omitempty"`
    OsVersion         string `json:"osVersion,omitempty"`
    PlatformVersion   string `json:"platformVersion,omitempty"`
    SerialNumber      string `json:"serialNumber,omitempty"`
    Status            string `json:"status,omitempty"`            // "ACTIVE", "DEPROVISIONED", etc.
    SupportEndDate    string `json:"supportEndDate,omitempty"`

    // Many additional fields for CPU, disk, network info...
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## ChromeOsDeviceAction Struct

```go
type ChromeOsDeviceAction struct {
    Action            string `json:"action,omitempty"`            // See table below
    DeprovisionReason string `json:"deprovisionReason,omitempty"` // Required for "deprovision"

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

### Chrome OS Action Values

| Action                       | Effect                                                     |
| ---------------------------- | ---------------------------------------------------------- |
| `"deprovision"`              | Remove device from management (requires DeprovisionReason) |
| `"disable"`                  | Disable the device                                         |
| `"reenable"`                 | Re-enable a disabled device                                |
| `"pre_provisioned_disable"`  | Disable before provisioning                                |
| `"pre_provisioned_reenable"` | Re-enable before provisioning                              |

### DeprovisionReason Values

| Value                           | Meaning                       |
| ------------------------------- | ----------------------------- |
| `"different_model_replacement"` | Replaced with different model |
| `"retiring_device"`             | Device being retired          |
| `"same_model_replacement"`      | Replaced with same model      |
| `"upgrade_transfer"`            | Transferred on upgrade        |

## ChromeOsMoveDevicesToOu Struct

```go
type ChromeOsMoveDevicesToOu struct {
    DeviceIds       []string `json:"deviceIds,omitempty"` // Up to 50 device IDs
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## ChromeOsDevices (List Response)

```go
type ChromeOsDevices struct {
    Chromeosdevices []*ChromeOsDevice `json:"chromeosdevices,omitempty"`
    NextPageToken   string            `json:"nextPageToken,omitempty"`
    Etag            string            `json:"etag,omitempty"`
    Kind            string            `json:"kind,omitempty"`
}
```

---

## ChromeosdevicesService

All methods require `customerId` as the first parameter.

```go
type ChromeosdevicesService struct{}

func NewChromeosdevicesService(s *Service) *ChromeosdevicesService
func (r *ChromeosdevicesService) Action(customerId string, resourceId string, chromeosdeviceaction *ChromeOsDeviceAction) *ChromeosdevicesActionCall
func (r *ChromeosdevicesService) Get(customerId string, deviceId string) *ChromeosdevicesGetCall
func (r *ChromeosdevicesService) List(customerId string) *ChromeosdevicesListCall
func (r *ChromeosdevicesService) MoveDevicesToOu(customerId string, orgUnitPath string, chromeosmovedevicestoou *ChromeOsMoveDevicesToOu) *ChromeosdevicesMoveDevicesToOuCall
func (r *ChromeosdevicesService) Patch(customerId string, deviceId string, chromeosdevice *ChromeOsDevice) *ChromeosdevicesPatchCall
func (r *ChromeosdevicesService) Update(customerId string, deviceId string, chromeosdevice *ChromeOsDevice) *ChromeosdevicesUpdateCall
```

## ChromeosdevicesListCall

```go
func (c *ChromeosdevicesListCall) Context(ctx context.Context) *ChromeosdevicesListCall
func (c *ChromeosdevicesListCall) Do(opts ...googleapi.CallOption) (*ChromeOsDevices, error)
func (c *ChromeosdevicesListCall) Fields(s ...googleapi.Field) *ChromeosdevicesListCall
func (c *ChromeosdevicesListCall) IfNoneMatch(entityTag string) *ChromeosdevicesListCall
func (c *ChromeosdevicesListCall) IncludeChildOrgunits(v bool) *ChromeosdevicesListCall
func (c *ChromeosdevicesListCall) MaxResults(maxResults int64) *ChromeosdevicesListCall // 1-300
func (c *ChromeosdevicesListCall) OrderBy(orderBy string) *ChromeosdevicesListCall
func (c *ChromeosdevicesListCall) OrgUnitPath(orgUnitPath string) *ChromeosdevicesListCall
func (c *ChromeosdevicesListCall) PageToken(pageToken string) *ChromeosdevicesListCall
func (c *ChromeosdevicesListCall) Pages(ctx context.Context, f func(*ChromeOsDevices) error) error
func (c *ChromeosdevicesListCall) Projection(projection string) *ChromeosdevicesListCall // "BASIC", "FULL"
func (c *ChromeosdevicesListCall) Query(query string) *ChromeosdevicesListCall
func (c *ChromeosdevicesListCall) SortOrder(sortOrder string) *ChromeosdevicesListCall
```

### OrderBy Values (Chrome OS)

- `"annotatedLocation"`
- `"annotatedUser"`
- `"lastSync"`
- `"notes"`
- `"serialNumber"`
- `"status"`

---

## MobileDevice Struct (Key Fields)

```go
type MobileDevice struct {
    DeviceId         string   `json:"deviceId,omitempty"`
    Email            []string `json:"email,omitempty"`
    Etag             string   `json:"etag,omitempty"`
    HardwareId       string   `json:"hardwareId,omitempty"`
    Kind             string   `json:"kind,omitempty"`
    Model            string   `json:"model,omitempty"`
    Name             []string `json:"name,omitempty"`
    Os               string   `json:"os,omitempty"`
    ResourceId       string   `json:"resourceId,omitempty"` // Used as resourceId in methods
    SerialNumber     string   `json:"serialNumber,omitempty"`
    Status           string   `json:"status,omitempty"`
    Type             string   `json:"type,omitempty"` // "android", "ios", etc.

    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

## MobileDeviceAction Struct

```go
type MobileDeviceAction struct {
    Action          string `json:"action,omitempty"` // See table below
    ForceSendFields []string `json:"-"`
    NullFields      []string `json:"-"`
}
```

### Mobile Action Values

| Action                               | Effect                      |
| ------------------------------------ | --------------------------- |
| `"admin_remote_wipe"`                | Remote wipe the device      |
| `"admin_account_wipe"`               | Remove managed account only |
| `"approve"`                          | Approve device for sync     |
| `"block"`                            | Block device from sync      |
| `"cancel_remote_wipe_then_activate"` | Cancel wipe, reactivate     |
| `"cancel_remote_wipe_then_block"`    | Cancel wipe, block          |

## MobileDevices (List Response)

```go
type MobileDevices struct {
    Mobiledevices []*MobileDevice `json:"mobiledevices,omitempty"`
    NextPageToken string          `json:"nextPageToken,omitempty"`
    Etag          string          `json:"etag,omitempty"`
    Kind          string          `json:"kind,omitempty"`
}
```

---

## MobiledevicesService

```go
type MobiledevicesService struct{}

func NewMobiledevicesService(s *Service) *MobiledevicesService
func (r *MobiledevicesService) Action(customerId string, resourceId string, mobiledeviceaction *MobileDeviceAction) *MobiledevicesActionCall
func (r *MobiledevicesService) Delete(customerId string, resourceId string) *MobiledevicesDeleteCall
func (r *MobiledevicesService) Get(customerId string, resourceId string) *MobiledevicesGetCall
func (r *MobiledevicesService) List(customerId string) *MobiledevicesListCall
```

## MobiledevicesListCall

```go
func (c *MobiledevicesListCall) Context(ctx context.Context) *MobiledevicesListCall
func (c *MobiledevicesListCall) Do(opts ...googleapi.CallOption) (*MobileDevices, error)
func (c *MobiledevicesListCall) Fields(s ...googleapi.Field) *MobiledevicesListCall
func (c *MobiledevicesListCall) MaxResults(maxResults int64) *MobiledevicesListCall
func (c *MobiledevicesListCall) OrderBy(orderBy string) *MobiledevicesListCall
func (c *MobiledevicesListCall) PageToken(pageToken string) *MobiledevicesListCall
func (c *MobiledevicesListCall) Pages(ctx context.Context, f func(*MobileDevices) error) error
func (c *MobiledevicesListCall) Projection(projection string) *MobiledevicesListCall // "BASIC", "FULL"
func (c *MobiledevicesListCall) Query(query string) *MobiledevicesListCall
func (c *MobiledevicesListCall) SortOrder(sortOrder string) *MobiledevicesListCall
```

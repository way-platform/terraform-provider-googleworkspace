# Chrome OS & Mobile Devices

## Chrome OS Devices

### Listing Devices

```go
var allDevices []*admin.ChromeOsDevice
err := svc.Chromeosdevices.List("my_customer").
    MaxResults(300).
    Projection("FULL").
    Pages(ctx, func(resp *admin.ChromeOsDevices) error {
        allDevices = append(allDevices, resp.Chromeosdevices...)
        return nil
    })
```

Filter by org unit:

```go
svc.Chromeosdevices.List("my_customer").
    OrgUnitPath("/Engineering").
    IncludeChildOrgunits(true).
    Do()
```

Filter by query:

```go
svc.Chromeosdevices.List("my_customer").
    Query("status:ACTIVE").
    Do()
```

### Getting a Device

```go
device, err := svc.Chromeosdevices.Get("my_customer", deviceId).
    Projection("FULL").
    Context(ctx).Do()
```

### Updating Device Annotations

```go
device := &admin.ChromeOsDevice{
    AnnotatedUser:     "jane@example.com",
    AnnotatedLocation: "Building A, Floor 2",
    AnnotatedAssetId:  "ASSET-001",
    Notes:             "Assigned 2024-01",
}
updated, err := svc.Chromeosdevices.Patch("my_customer", deviceId, device).
    Context(ctx).Do()
```

### Device Actions

```go
// Disable a device
action := &admin.ChromeOsDeviceAction{Action: "disable"}
err := svc.Chromeosdevices.Action("my_customer", deviceId, action).Context(ctx).Do()

// Deprovision a device (requires reason)
action := &admin.ChromeOsDeviceAction{
    Action:            "deprovision",
    DeprovisionReason: "retiring_device",
}
err := svc.Chromeosdevices.Action("my_customer", deviceId, action).Context(ctx).Do()

// Re-enable a disabled device
action := &admin.ChromeOsDeviceAction{Action: "reenable"}
err := svc.Chromeosdevices.Action("my_customer", deviceId, action).Context(ctx).Do()
```

### Moving Devices to an Org Unit

```go
move := &admin.ChromeOsMoveDevicesToOu{
    DeviceIds: []string{deviceId1, deviceId2}, // Max 50 per call
}
err := svc.Chromeosdevices.MoveDevicesToOu("my_customer", "/NewOrgUnit", move).
    Context(ctx).Do()
```

### Projection Values

| Value     | Returns                                    |
| --------- | ------------------------------------------ |
| `"BASIC"` | Core fields only                           |
| `"FULL"`  | All fields including hardware/network info |

### Query Syntax (Chrome OS)

| Field           | Example                           |
| --------------- | --------------------------------- |
| `status`        | `status:ACTIVE`                   |
| `user`          | `user:jane`                       |
| `location`      | `location:Building`               |
| `notes`         | `notes:replacement`               |
| `serial_number` | `serial_number:ABC123`            |
| `asset_id`      | `asset_id:ASSET-001`              |
| `register`      | `register:2024-01-01..2024-12-31` |

---

## Mobile Devices

### Listing Devices

```go
var allDevices []*admin.MobileDevice
err := svc.Mobiledevices.List("my_customer").
    Projection("FULL").
    Pages(ctx, func(resp *admin.MobileDevices) error {
        allDevices = append(allDevices, resp.Mobiledevices...)
        return nil
    })
```

Filter by query:

```go
svc.Mobiledevices.List("my_customer").
    Query("email:jane@example.com").
    Do()
```

### Getting a Device

```go
device, err := svc.Mobiledevices.Get("my_customer", resourceId).
    Projection("FULL").
    Context(ctx).Do()
```

Note: Mobile devices use `resourceId` (not `deviceId`) in method parameters.

### Device Actions

```go
// Approve a device
action := &admin.MobileDeviceAction{Action: "approve"}
err := svc.Mobiledevices.Action("my_customer", resourceId, action).Context(ctx).Do()

// Block a device
action := &admin.MobileDeviceAction{Action: "block"}
err := svc.Mobiledevices.Action("my_customer", resourceId, action).Context(ctx).Do()

// Remote wipe
action := &admin.MobileDeviceAction{Action: "admin_remote_wipe"}
err := svc.Mobiledevices.Action("my_customer", resourceId, action).Context(ctx).Do()

// Account wipe (removes managed account only)
action := &admin.MobileDeviceAction{Action: "admin_account_wipe"}
err := svc.Mobiledevices.Action("my_customer", resourceId, action).Context(ctx).Do()
```

### Deleting a Device Record

```go
err := svc.Mobiledevices.Delete("my_customer", resourceId).Context(ctx).Do()
```

### Query Syntax (Mobile)

| Field    | Example                  |
| -------- | ------------------------ |
| `email`  | `email:jane@example.com` |
| `name`   | `name:Jane`              |
| `status` | `status:APPROVED`        |
| `os`     | `os:android`             |
| `model`  | `model:Pixel`            |
| `serial` | `serial:ABC123`          |
| `type`   | `type:android`           |

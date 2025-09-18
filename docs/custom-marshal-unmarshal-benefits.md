# Custom Marshal/Unmarshal Benefits in DataCrunch Go SDK

This document explains why the DataCrunch Go SDK implements custom JSON marshaling and unmarshaling, with concrete examples from the service layer code.

## Overview

The DataCrunch Go SDK uses custom marshal/unmarshal implementations instead of Go's standard `encoding/json` package to handle specific API requirements that standard JSON processing cannot address. These custom implementations are found in the `internal/protocol/json/jsonutil` package and provide enhanced functionality for API communication.

## Key Benefits with Real Examples

### 1. **Custom Field Name Mapping with `locationName` Tags**

**Problem**: APIs often use field names that don't match Go naming conventions or require different mappings for different contexts (URL parameters vs JSON body).

**Solution**: The SDK supports `locationName` tags for flexible field mapping.

#### Example from Instance Availability Service

```go
// From service/instanceavailability/api.go
type InstanceAvailabilityResponse struct {
    LocationCode   string   `json:"location_code" locationName:"location_code"`
    Availabilities []string `json:"availabilities" locationName:"availabilities"`
}
```

#### Example from Volumes Service

```go
// From service/volumes/api.go
type GetVolumeInput struct {
    ID string `location:"uri" locationName:"id"`
}

type DeleteVolumeInput struct {
    ID string `location:"uri" locationName:"id"`
}
```

**Benefit**: The same struct field can be mapped to different names in different contexts (URI parameters, query strings, JSON body) without requiring separate struct definitions.

### 2. **Mixed Response Type Handling**

**Problem**: Some APIs return different response types for the same endpoint - sometimes JSON objects, sometimes plain strings.

**Solution**: Dynamic unmarshaler switching based on response content type.

#### Example from Instance Creation

```go
// From service/instance/api.go
func (c *Instance) CreateInstance(input *CreateInstanceInput) (string, error) {
    // ... operation setup ...
    
    var instanceID string
    req := c.newRequest(op, input, &instanceID)

    // Replace JSON unmarshaler with string unmarshaler for plain text response
    req.Handlers.Unmarshal.RemoveByName("datacrunchsdk.restjson.Unmarshal")
    req.Handlers.Unmarshal.PushBackNamed(restjson.StringUnmarshalHandler)

    return instanceID, req.Send()
}
```

#### Example from SSH Key Creation

```go
// From service/sshkeys/api.go
func (c *SSHKey) CreateSSHKey(input *CreateSSHKeyInput) (string, error) {
    // ... operation setup ...
    
    var sshKey string
    req := c.newRequest(op, input, &sshKey)

    // API returns plain string, not JSON
    req.Handlers.Unmarshal.RemoveByName("datacrunchsdk.restjson.Unmarshal")
    req.Handlers.Unmarshal.PushBackNamed(restjson.StringUnmarshalHandler)

    return sshKey, req.Send()
}
```

**Benefit**: Single SDK can handle APIs that return both JSON objects and plain text responses without requiring different client implementations.

### 3. **Null vs Missing Field Distinction**

**Problem**: APIs often distinguish between fields that are `null`, missing entirely, or have zero values. Standard JSON unmarshaling cannot reliably detect this difference.

**Solution**: Pointer types combined with custom field presence detection.

#### Example from Instance Types

```go
// From service/instancetypes/api.go
type CPU struct {
    Description   string `json:"description"`
    NumberOfCores *int64 `json:"number_of_cores"` // Pointer to detect null vs missing
}

type GPU struct {
    Description  string `json:"description"`
    NumberOfGPUs *int64 `json:"number_of_gpus"` // Pointer to detect null vs missing
}

type InstanceTypeResponse struct {
    // ... other fields ...
    P2P         *string `json:"p2p"`         // Can be null
    DisplayName *string `json:"display_name"` // Can be null
}
```

#### Example from Instance Service

```go
// From service/instance/api.go
type ListInstancesResponse struct {
    // ... other fields ...
    StartupScriptID *string `json:"startup_script_id"` // Can be null
    JupyterToken    *string `json:"jupyter_token"`     // Can be null
}
```

**Benefit**: Applications can distinguish between:
- Field not provided in API response
- Field explicitly set to `null`
- Field set to empty string `""`

### 4. **Enhanced Error Handling with Context**

**Problem**: Standard JSON errors provide limited context about what failed to parse.

**Solution**: Custom error types that include the raw data that failed to unmarshal.

#### Example Error Handling

```go
// From internal/protocol/json/jsonutil/unmarshal.go
func UnmarshalJSONError(v interface{}, stream io.Reader) error {
    var errBuf bytes.Buffer
    body := io.TeeReader(stream, &errBuf)

    err := json.NewDecoder(body).Decode(v)
    if err != nil {
        msg := "failed decoding error message"
        if err == io.EOF {
            msg = "error message missing"
            err = nil
        }
        return dcerr.NewUnmarshalError(err, msg, errBuf.Bytes())
    }
    return nil
}
```

**Benefit**: When JSON parsing fails, developers get:
- The original error
- A descriptive message
- The raw bytes that failed to parse (for debugging)

### 5. **Enum Type Safety with String Constants**

**Problem**: APIs use string enums, but Go doesn't have built-in enum support.

**Solution**: Custom string types with predefined constants.

#### Example from Instance Service

```go
// From service/instance/api.go
type InstanceActionType string

const (
    InstanceActionBoot          InstanceActionType = "boot"
    InstanceActionStart         InstanceActionType = "start"
    InstanceActionShutdown      InstanceActionType = "shutdown"
    InstanceActionDelete        InstanceActionType = "delete"
    InstanceActionDiscontinue   InstanceActionType = "discontinue"
    InstanceActionHibernate     InstanceActionType = "hibernate"
    InstanceActionConfigureSpot InstanceActionType = "configure_spot"
    InstanceActionForceShutdown InstanceActionType = "force_shutdown"
)

type InstanceStatus string

const (
    InstanceStatusRunning      InstanceStatus = "running"
    InstanceStatusProvisioning InstanceStatus = "provisioning"
    InstanceStatusOffline      InstanceStatus = "offline"
    // ... more statuses
)
```

#### Example from Volumes Service

```go
// From service/volumes/api.go
type VolumeStatus string

const (
    VolumeStatusOrdered   VolumeStatus = "ordered"
    VolumeStatusAttached  VolumeStatus = "attached"
    VolumeStatusAttaching VolumeStatus = "attaching"
    VolumeStatusDetached  VolumeStatus = "detached"
    VolumeStatusDeleted   VolumeStatus = "deleted"
)
```

**Benefit**: 
- Type safety at compile time
- IDE autocompletion for valid values
- Clear documentation of valid enum values
- Runtime validation possible

### 6. **Flexible Query Parameter Handling**

**Problem**: Different endpoints require different ways of encoding parameters (query string, URI path, JSON body).

**Solution**: Location-aware parameter encoding using struct tags.

#### Example from Instance Service

```go
// From service/instance/api.go
type ListInstancesInput struct {
    Status string `location:"querystring" locationName:"status"`
}
```

**Benefit**: Single struct can specify where each field should be encoded in the HTTP request.

## Real-World API Challenges Addressed

### Challenge 1: Inconsistent API Response Formats

The Instance Availability service demonstrates a real challenge:

```go
// From service/instanceavailability/api.go
// CheckInstanceAvailability is DISABLED due to inconsistent API behavior:
// - With insufficient params: returns array
// - With instance type: returns boolean as string ("true"/"false")
```

**Solution**: Custom unmarshalers can handle multiple response formats for the same endpoint.

### Challenge 2: Mixed Data Types in Responses

Instance types return prices as strings (not numbers) to avoid floating-point precision issues:

```go
// From service/instancetypes/api.go
type InstanceTypeResponse struct {
    PricePerHour        string `json:"price_per_hour"`        // String, not float64
    SpotPrice           string `json:"spot_price"`
    DynamicPrice        string `json:"dynamic_price"`
    MaxDynamicPrice     string `json:"max_dynamic_price"`
}
```

But price history uses actual numbers:

```go
type PriceHistoryEntry struct {
    FixedPricePerHour   float64 `json:"fixed_price_per_hour"`   // float64, not string
    DynamicPricePerHour float64 `json:"dynamic_price_per_hour"`
}
```

**Benefit**: Custom unmarshaling can handle these inconsistencies transparently.

## Performance Benefits

1. **Single Pass Processing**: Custom unmarshalers can perform validation, type conversion, and field mapping in a single pass
2. **Reduced Memory Allocation**: Avoid creating intermediate representations
3. **Optimized for Use Case**: Tailored for the specific needs of the DataCrunch API rather than general-purpose JSON

## Conclusion

The custom marshal/unmarshal implementation in the DataCrunch Go SDK provides essential functionality that standard JSON processing cannot handle:

- **API Compatibility**: Handles real-world API quirks and inconsistencies
- **Type Safety**: Maintains Go's strong typing while working with flexible JSON
- **Developer Experience**: Better error messages and IDE support
- **Maintainability**: Single codebase handles multiple API response formats
- **Robustness**: Graceful handling of API changes and edge cases

This approach allows the SDK to provide a clean, type-safe Go interface while handling the complexities of real-world API communication behind the scenes.

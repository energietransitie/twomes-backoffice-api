package devicetype

// A DeviceType contains information about a group of devices with the same functionality.
type DeviceType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Create a new DeviceType.
func MakeDeviceType(name string) DeviceType {
	return DeviceType{
		Name: name,
	}
}

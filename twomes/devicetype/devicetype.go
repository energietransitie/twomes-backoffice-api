package devicetype

// A DeviceType contains information about a group of devices with the same functionality.
type DeviceType struct {
	ID                    uint   `json:"id"`
	Name                  string `json:"name"`
	InstallationManualURL string `json:"installation_manual_url"`
	InfoURL               string `json:"info_url"`
}

// Create a new DeviceType.
func MakeDeviceType(name, installationManualURL, infoURL string) DeviceType {
	return DeviceType{
		Name:                  name,
		InstallationManualURL: installationManualURL,
		InfoURL:               infoURL,
	}
}

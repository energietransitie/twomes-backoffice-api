package building

import "github.com/energietransitie/twomes-backoffice-api/twomes/device"

// A Building belongs to a research subject where we take regular measurements with devices.
type Building struct {
	ID        uint             `json:"id"`
	AccountID uint             `json:"account_id"`
	Longitude float32          `json:"longitude"`
	Latitude  float32          `json:"latitude"`
	TZName    string           `json:"tz_name"`
	Devices   []*device.Device `json:"devices,omitempty"`
}

// Create a new Building.
func MakeBuilding(accountID uint, long, lat float32, tzName string) Building {
	return Building{
		AccountID: accountID,
		Longitude: long,
		Latitude:  lat,
		TZName:    tzName,
	}
}

// Add a [Device] to a Building.
func (b *Building) AddDevice(device *device.Device) {
	b.Devices = append(b.Devices, device)
}

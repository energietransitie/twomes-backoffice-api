package device

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/measurement"
	"github.com/energietransitie/twomes-backoffice-api/twomes/property"
)

// A DeviceRepository can load, store and delete devices.
type DeviceRepository interface {
	Find(device Device) (Device, error)
	FindCloudFeedAuthCreationTimeFromDeviceID(deviceID uint) (*time.Time, error)
	GetProperties(device Device) ([]property.Property, error)
	GetMeasurements(device Device, filters map[string]string) ([]measurement.Measurement, error)
	GetAll() ([]Device, error)
	Create(Device) (Device, error)
	Update(Device) (Device, error)
	Delete(Device) error
}

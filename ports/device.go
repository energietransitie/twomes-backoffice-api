package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/device"
	"github.com/energietransitie/twomes-backoffice-api/twomes/measurement"
	"github.com/energietransitie/twomes-backoffice-api/twomes/property"
	"github.com/energietransitie/twomes-backoffice-api/twomes/upload"
)

// A DeviceRepository can load, store and delete devices.
type DeviceRepository interface {
	Find(device device.Device) (device.Device, error)
	FindCloudFeedAuthCreationTimeFromDeviceID(deviceID uint) (*time.Time, error)
	GetProperties(device device.Device) ([]property.Property, error)
	GetMeasurements(device device.Device, filters map[string]string) ([]measurement.Measurement, error)
	GetAll() ([]device.Device, error)
	Create(device.Device) (device.Device, error)
	Update(device.Device) (device.Device, error)
	Delete(device.Device) error
}

// DeviceService exposes all operations that can be performed on a [device.Device].
type DeviceService interface {
	Create(name string, buildingID, accountID uint, activationSecret string) (device.Device, error)
	GetByID(id uint) (device.Device, error)
	GetByName(name string) (device.Device, error)
	Activate(name, activationSecret string) (device.Device, error)
	AddUpload(id uint, upload upload.Upload) (device.Device, error)
	GetAccountByDeviceID(id uint) (uint, error)
	GetMeasurementsByDeviceID(id uint, filters map[string]string) ([]measurement.Measurement, error)
	GetPropertiesByDeviceID(id uint) ([]property.Property, error)
}

package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes"

// A DeviceRepository can load, store and delete devices.
type DeviceRepository interface {
	Find(device twomes.Device) (twomes.Device, error)
	GetProperties(device twomes.Device) ([]twomes.Property, error)
	GetMeasurements(device twomes.Device, filters map[string]string) ([]twomes.Measurement, error)
	GetAll() ([]twomes.Device, error)
	Create(twomes.Device) (twomes.Device, error)
	Update(twomes.Device) (twomes.Device, error)
	Delete(twomes.Device) error
}

// DeviceService exposes all operations that can be performed on a [twomes.Device].
type DeviceService interface {
	Create(name string, buildingID, accountID uint, activationSecret string) (twomes.Device, error)
	GetByID(id uint) (twomes.Device, error)
	GetByName(name string) (twomes.Device, error)
	Activate(name, activationSecret string) (twomes.Device, error)
	AddUpload(id uint, upload twomes.Upload) (twomes.Device, error)
	GetAccountByDeviceID(id uint) (uint, error)
	GetMeasurementsByDeviceID(id uint, filters map[string]string) ([]twomes.Measurement, error)
	GetPropertiesByDeviceID(id uint) ([]twomes.Property, error)
}

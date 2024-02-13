package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes/devicetype"

// A DeviceTypeRepository can load, store and delete device types.
type DeviceTypeRepository interface {
	Find(deviceType devicetype.DeviceType) (devicetype.DeviceType, error)
	GetAll() ([]devicetype.DeviceType, error)
	Create(devicetype.DeviceType) (devicetype.DeviceType, error)
	Delete(devicetype.DeviceType) error
}

// DeviceTypeService exposes all operations that can be performed on a [devicetype.DeviceType].
type DeviceTypeService interface {
	Create(name, installationManualURL, infoURL string) (devicetype.DeviceType, error)
	Find(deviceType devicetype.DeviceType) (devicetype.DeviceType, error)
	GetByHash(deviceTypeHash string) (devicetype.DeviceType, error)
	GetByID(id uint) (devicetype.DeviceType, error)
	GetByName(name string) (devicetype.DeviceType, error)
}

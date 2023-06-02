package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes"

// A DeviceTypeRepository can load, store and delete device types.
type DeviceTypeRepository interface {
	Find(deviceType twomes.DeviceType) (twomes.DeviceType, error)
	GetAll() ([]twomes.DeviceType, error)
	Create(twomes.DeviceType) (twomes.DeviceType, error)
	Delete(twomes.DeviceType) error
}

// DeviceTypeService exposes all operations that can be performed on a [twomes.DeviceType].
type DeviceTypeService interface {
	Create(name, installationManualURL, infoURL string) (twomes.DeviceType, error)
	Find(deviceType twomes.DeviceType) (twomes.DeviceType, error)
	GetByHash(deviceTypeHash string) (twomes.DeviceType, error)
	GetByID(id uint) (twomes.DeviceType, error)
	GetByName(name string) (twomes.DeviceType, error)
}

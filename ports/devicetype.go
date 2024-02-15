package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes/devicetype"

// DeviceTypeService exposes all operations that can be performed on a [devicetype.DeviceType].
type DeviceTypeService interface {
	Create(name, installationManualURL, infoURL string) (devicetype.DeviceType, error)
	Find(deviceType devicetype.DeviceType) (devicetype.DeviceType, error)
	GetByHash(deviceTypeHash string) (devicetype.DeviceType, error)
	GetByID(id uint) (devicetype.DeviceType, error)
	GetByName(name string) (devicetype.DeviceType, error)
}

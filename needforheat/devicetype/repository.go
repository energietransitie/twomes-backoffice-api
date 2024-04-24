package devicetype

// A DeviceTypeRepository can load, store and delete device types.
type DeviceTypeRepository interface {
	Find(deviceType DeviceType) (DeviceType, error)
	GetAll() ([]DeviceType, error)
	Create(DeviceType) (DeviceType, error)
	Delete(DeviceType) error
}

package property

// A DeviceTypeRepository can load, store and delete properties.
type PropertyRepository interface {
	Find(property Property) (Property, error)
	GetAll() ([]Property, error)
	Create(Property) (Property, error)
	Delete(Property) error
}

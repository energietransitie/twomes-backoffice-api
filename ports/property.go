package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes/property"

// A DeviceTypeRepository can load, store and delete properties.
type PropertyRepository interface {
	Find(property property.Property) (property.Property, error)
	GetAll() ([]property.Property, error)
	Create(property.Property) (property.Property, error)
	Delete(property.Property) error
}

// PropertyService exposes all operations that can be performed on a [property.Property].
type PropertyService interface {
	Create(name string) (property.Property, error)
	Find(property property.Property) (property.Property, error)
	GetByID(id uint) (property.Property, error)
	GetByName(name string) (property.Property, error)
}

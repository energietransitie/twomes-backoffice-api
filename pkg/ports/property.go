package ports

import "github.com/energietransitie/twomes-backoffice-api/pkg/twomes"

// A DeviceTypeRepository can load, store and delete properties.
type PropertyRepository interface {
	Find(property twomes.Property) (twomes.Property, error)
	GetAll() ([]twomes.Property, error)
	Create(twomes.Property) (twomes.Property, error)
	Delete(twomes.Property) error
}

// PropertyService exposes all operations that can be performed on a [twomes.Property].
type PropertyService interface {
	Create(name string) (twomes.Property, error)
	Find(property twomes.Property) (twomes.Property, error)
	GetByID(id uint) (twomes.Property, error)
	GetByName(name string) (twomes.Property, error)
}

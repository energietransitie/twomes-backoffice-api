package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes/property"

// PropertyService exposes all operations that can be performed on a [property.Property].
type PropertyService interface {
	Create(name string) (property.Property, error)
	Find(property property.Property) (property.Property, error)
	GetByID(id uint) (property.Property, error)
	GetByName(name string) (property.Property, error)
}

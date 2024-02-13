package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes/building"

// A BuildingRepository can load, store and delete buildings.
type BuildingRepository interface {
	Find(building building.Building) (building.Building, error)
	GetAll() ([]building.Building, error)
	Create(building.Building) (building.Building, error)
	Delete(building.Building) error
}

// BuildingService exposes all operations that can be performed on a [building.Building].
type BuildingService interface {
	Create(accountID uint, long, lat float32, tzName string) (building.Building, error)
	GetByID(id uint) (building.Building, error)
}

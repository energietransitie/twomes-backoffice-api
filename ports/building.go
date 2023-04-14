package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes"

// A BuildingRepository can load, store and delete buildings.
type BuildingRepository interface {
	Find(building twomes.Building) (twomes.Building, error)
	GetAll() ([]twomes.Building, error)
	Create(twomes.Building) (twomes.Building, error)
	Delete(twomes.Building) error
}

// BuildingService exposes all operations that can be performed on a [twomes.Building].
type BuildingService interface {
	Create(accountID uint, long, lat float32, tzName string) (twomes.Building, error)
	GetByID(id uint) (twomes.Building, error)
}

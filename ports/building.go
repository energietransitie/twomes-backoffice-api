package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes/building"

// BuildingService exposes all operations that can be performed on a [building.Building].
type BuildingService interface {
	Create(accountID uint, long, lat float32, tzName string) (building.Building, error)
	GetByID(id uint) (building.Building, error)
}

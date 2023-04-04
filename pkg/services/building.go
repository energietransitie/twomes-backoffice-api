package services

import (
	"github.com/energietransitie/twomes-backoffice-api/pkg/ports"
	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
)

type BuildingService struct {
	repository ports.BuildingRepository
}

// Create a new BuildingService.
func NewBuildingService(repository ports.BuildingRepository) *BuildingService {
	return &BuildingService{
		repository: repository,
	}
}

func (s *BuildingService) Create(accountID uint, long float32, lat float32, tzName string) (twomes.Building, error) {
	building := twomes.MakeBuilding(accountID, long, lat, tzName)
	return s.repository.Create(building)
}

func (s *BuildingService) GetByID(id uint) (twomes.Building, error) {
	return s.repository.Find(twomes.Building{ID: id})
}

package services

import (
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

type BuildingService struct {
	repository ports.BuildingRepository

	// Services used when getting device info.
	uploadService ports.UploadService
}

// Create a new BuildingService.
func NewBuildingService(repository ports.BuildingRepository, uploadService ports.UploadService) *BuildingService {
	return &BuildingService{
		repository:    repository,
		uploadService: uploadService,
	}
}

func (s *BuildingService) Create(accountID uint, long float32, lat float32, tzName string) (twomes.Building, error) {
	building := twomes.MakeBuilding(accountID, long, lat, tzName)
	return s.repository.Create(building)
}

func (s *BuildingService) GetByID(id uint) (twomes.Building, error) {
	building, err := s.repository.Find(twomes.Building{ID: id})
	if err != nil {
		return twomes.Building{}, err
	}

	for _, device := range building.Devices {
		device.LatestUpload, _, err = s.uploadService.GetLatestUploadTimeForDeviceWithID(device.ID)
		if err != nil {
			return twomes.Building{}, err
		}
	}

	return building, nil
}

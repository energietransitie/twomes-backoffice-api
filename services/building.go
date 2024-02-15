package services

import (
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes/building"
)

type BuildingService struct {
	repository building.BuildingRepository

	// Services used when getting device info.
	uploadService ports.UploadService
}

// Create a new BuildingService.
func NewBuildingService(repository building.BuildingRepository, uploadService ports.UploadService) *BuildingService {
	return &BuildingService{
		repository:    repository,
		uploadService: uploadService,
	}
}

func (s *BuildingService) Create(accountID uint, long float32, lat float32, tzName string) (building.Building, error) {
	b := building.MakeBuilding(accountID, long, lat, tzName)
	return s.repository.Create(b)
}

func (s *BuildingService) GetByID(id uint) (building.Building, error) {
	b, err := s.repository.Find(building.Building{ID: id})
	if err != nil {
		return building.Building{}, err
	}

	for _, device := range b.Devices {
		device.LatestUpload, _, err = s.uploadService.GetLatestUploadTimeForDeviceWithID(device.ID)
		if err != nil {
			return building.Building{}, err
		}
	}

	return b, nil
}

package services

import (
	"github.com/energietransitie/twomes-backoffice-api/pkg/ports"
	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
)

type DeviceTypeService struct {
	repository ports.DeviceTypeRepository

	// Service used when creating a device type.
	propertyService ports.PropertyService
}

// Create a new DeviceTypeService.
func NewDeviceTypeService(repository ports.DeviceTypeRepository, propertyService ports.PropertyService) *DeviceTypeService {
	return &DeviceTypeService{
		repository:      repository,
		propertyService: propertyService,
	}
}

func (s *DeviceTypeService) Create(name, installationManualURL, infoURL string, properties []twomes.Property, uploadInterval twomes.Duration) (twomes.DeviceType, error) {
	for i := range properties {
		var err error
		properties[i], err = s.propertyService.Find(properties[i])
		if err != nil {
			return twomes.DeviceType{}, err
		}
	}

	deviceType := twomes.MakeDeviceType(name, installationManualURL, infoURL, properties, uploadInterval)
	return s.repository.Create(deviceType)
}

func (s *DeviceTypeService) Find(deviceType twomes.DeviceType) (twomes.DeviceType, error) {
	return s.repository.Find(deviceType)
}

func (s *DeviceTypeService) GetByID(id uint) (twomes.DeviceType, error) {
	return s.repository.Find(twomes.DeviceType{ID: id})
}

func (s *DeviceTypeService) GetByName(name string) (twomes.DeviceType, error) {
	return s.repository.Find(twomes.DeviceType{Name: name})
}

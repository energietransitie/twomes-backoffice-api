package services

import (
	"errors"

	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/twomes"
)

var (
	ErrDeviceDoesNotBelongToAccount   = errors.New("device does not belong to this account")
	ErrBuildingDoesNotBelongToAccount = errors.New("building does not belong to this account")
)

type DeviceService struct {
	repository ports.DeviceRepository

	// Services used when activating a device.
	authService ports.AuthorizationService

	// Services used when creating a device.
	deviceTypeService ports.DeviceTypeService
	BuildingService   ports.BuildingService
}

// Create a new DeviceService.
func NewDeviceService(repository ports.DeviceRepository, authService ports.AuthorizationService, deviceTypeService ports.DeviceTypeService, BuildingService ports.BuildingService) *DeviceService {
	return &DeviceService{
		repository:        repository,
		authService:       authService,
		deviceTypeService: deviceTypeService,
		BuildingService:   BuildingService,
	}
}

func (s *DeviceService) Create(name string, deviceType twomes.DeviceType, buildingID, accountID uint, activationSecret string) (twomes.Device, error) {
	building, err := s.BuildingService.GetByID(buildingID)
	if err != nil {
		return twomes.Device{}, err
	}

	if building.AccountID != accountID {
		return twomes.Device{}, ErrBuildingDoesNotBelongToAccount
	}

	deviceType, err = s.deviceTypeService.Find(deviceType)
	if err != nil {
		return twomes.Device{}, err
	}

	device := twomes.MakeDevice(name, deviceType, buildingID, activationSecret)
	return s.repository.Create(device)
}

func (s *DeviceService) GetByID(id uint) (twomes.Device, error) {
	return s.repository.Find(twomes.Device{ID: id})
}

func (s *DeviceService) GetByName(name string) (twomes.Device, error) {
	device, err := s.repository.Find(twomes.Device{Name: name})
	if err != nil {
		return twomes.Device{}, err
	}

	device.UpdateHealth()

	return device, nil
}

func (s *DeviceService) Activate(name, activationSecret string) (twomes.Device, error) {
	device, err := s.repository.Find(twomes.Device{Name: name})
	if err != nil {
		return twomes.Device{}, err
	}

	err = device.Activate(activationSecret)
	if err != nil {
		return device, err
	}

	device, err = s.repository.Update(device)
	if err != nil {
		return device, err
	}

	device.AuthorizationToken, err = s.authService.CreateToken(twomes.DeviceToken, device.ID)
	if err != nil {
		return twomes.Device{}, err
	}

	return device, nil
}

func (s *DeviceService) AddUpload(id uint, upload twomes.Upload) (twomes.Device, error) {
	device, err := s.repository.Find(twomes.Device{ID: id})
	if err != nil {
		return twomes.Device{}, err
	}

	device.AddUpload(upload)

	return s.repository.Update(device)
}

func (s *DeviceService) GetAccountByDeviceID(id uint) (uint, error) {
	device, err := s.repository.Find(twomes.Device{ID: id})
	if err != nil {
		return 0, err
	}

	building, err := s.BuildingService.GetByID(device.BuildingID)
	if err != nil {
		return 0, err
	}

	return building.AccountID, nil
}

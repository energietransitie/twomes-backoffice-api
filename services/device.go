package services

import (
	"errors"
	"strings"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/authorization"
	"github.com/energietransitie/twomes-backoffice-api/twomes/device"
	"github.com/energietransitie/twomes-backoffice-api/twomes/measurement"
	"github.com/energietransitie/twomes-backoffice-api/twomes/property"
	"github.com/energietransitie/twomes-backoffice-api/twomes/upload"
)

var (
	ErrDeviceDoesNotBelongToAccount   = errors.New("device does not belong to this account")
	ErrBuildingDoesNotBelongToAccount = errors.New("building does not belong to this account")
	ErrDeviceTypeNameInvalid          = errors.New("device type name invalid")
)

type DeviceService struct {
	repository device.DeviceRepository

	// Services used when activating a device.
	authService *AuthorizationService

	// Services used when creating a device.
	deviceTypeService *DeviceTypeService
	buildingService   *BuildingService

	// Services used when getting device info.
	uploadService *UploadService
}

// Create a new DeviceService.
func NewDeviceService(repository device.DeviceRepository, authService *AuthorizationService, deviceTypeService *DeviceTypeService, BuildingService *BuildingService, uploadService *UploadService) *DeviceService {
	return &DeviceService{
		repository:        repository,
		authService:       authService,
		deviceTypeService: deviceTypeService,
		buildingService:   BuildingService,
		uploadService:     uploadService,
	}
}

func (s *DeviceService) Create(name string, buildingID, accountID uint, activationSecret string) (device.Device, error) {
	building, err := s.buildingService.GetByID(buildingID)
	if err != nil {
		return device.Device{}, err
	}

	if building.AccountID != accountID {
		return device.Device{}, ErrBuildingDoesNotBelongToAccount
	}

	splitDeviceTypeName := strings.Split(name, "-")
	if len(splitDeviceTypeName) != 2 {
		return device.Device{}, ErrDeviceTypeNameInvalid
	}

	deviceTypeHash := splitDeviceTypeName[0]
	deviceType, err := s.deviceTypeService.GetByHash(deviceTypeHash)
	if err != nil {
		return device.Device{}, err
	}

	d := device.MakeDevice(name, deviceType, buildingID, activationSecret)
	return s.repository.Create(d)
}

func (s *DeviceService) GetByID(id uint) (device.Device, error) {
	return s.repository.Find(device.Device{ID: id})
}

func (s *DeviceService) GetByName(name string) (device.Device, error) {
	d, err := s.repository.Find(device.Device{Name: name})
	if err != nil {
		return device.Device{}, err
	}

	d.LatestUpload, _, err = s.uploadService.GetLatestUploadTimeForDeviceWithID(d.ID)
	if err != nil {
		return device.Device{}, err
	}

	return d, nil
}

func (s *DeviceService) Activate(name, activationSecret string) (device.Device, error) {
	d, err := s.repository.Find(device.Device{Name: name})
	if err != nil {
		return device.Device{}, err
	}

	err = d.Activate(activationSecret)
	if err != nil {
		return d, err
	}

	d, err = s.repository.Update(d)
	if err != nil {
		return d, err
	}

	d.AuthorizationToken, err = s.authService.CreateToken(authorization.DeviceToken, d.ID, time.Time{})
	if err != nil {
		return device.Device{}, err
	}

	return d, nil
}

func (s *DeviceService) AddUpload(id uint, upload upload.Upload) (device.Device, error) {
	d, err := s.repository.Find(device.Device{ID: id})
	if err != nil {
		return device.Device{}, err
	}

	d.AddUpload(upload)

	return s.repository.Update(d)
}

func (s *DeviceService) GetAccountByDeviceID(id uint) (uint, error) {
	device, err := s.repository.Find(device.Device{ID: id})
	if err != nil {
		return 0, err
	}

	building, err := s.buildingService.GetByID(device.BuildingID)
	if err != nil {
		return 0, err
	}

	return building.AccountID, nil
}

func (s *DeviceService) GetMeasurementsByDeviceID(id uint, filters map[string]string) ([]measurement.Measurement, error) {
	measurements, err := s.repository.GetMeasurements(device.Device{ID: id}, filters)
	if err != nil {
		return nil, err
	}

	return measurements, nil
}

func (s *DeviceService) GetPropertiesByDeviceID(id uint) ([]property.Property, error) {
	properties, err := s.repository.GetProperties(device.Device{ID: id})
	if err != nil {
		return nil, err
	}

	return properties, nil
}

package services

import (
	"errors"
	"fmt"

	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes/devicetype"
	"github.com/sigurn/crc16"
	"github.com/sirupsen/logrus"
)

var (
	ErrHashDoesNotMatchType = errors.New("hash does not match a device type")
)

type DeviceTypeService struct {
	repository ports.DeviceTypeRepository

	// Service used when creating a device type.
	propertyService ports.PropertyService

	// Hashed device types.
	hashedDeviceTypes map[string]string
}

// Create a new DeviceTypeService.
func NewDeviceTypeService(repository ports.DeviceTypeRepository, propertyService ports.PropertyService) *DeviceTypeService {
	deviceTypeService := &DeviceTypeService{
		repository:      repository,
		propertyService: propertyService,
	}

	deviceTypeService.updateDeviceTypeHashes()

	return deviceTypeService
}

func (s *DeviceTypeService) Create(name, installationManualURL, infoURL string) (devicetype.DeviceType, error) {
	deviceType := devicetype.MakeDeviceType(name, installationManualURL, infoURL)

	deviceType, err := s.repository.Create(deviceType)
	if err != nil {
		return deviceType, err
	}

	s.updateDeviceTypeHashes()

	return deviceType, nil
}

func (s *DeviceTypeService) Find(deviceType devicetype.DeviceType) (devicetype.DeviceType, error) {
	return s.repository.Find(deviceType)
}

func (s *DeviceTypeService) GetByHash(deviceTypeHash string) (devicetype.DeviceType, error) {
	name, ok := s.hashedDeviceTypes[deviceTypeHash]
	if !ok {
		return devicetype.DeviceType{}, ErrHashDoesNotMatchType
	}

	return s.repository.Find(devicetype.DeviceType{Name: name})
}

func (s *DeviceTypeService) GetByID(id uint) (devicetype.DeviceType, error) {
	return s.repository.Find(devicetype.DeviceType{ID: id})
}

func (s *DeviceTypeService) GetByName(name string) (devicetype.DeviceType, error) {
	return s.repository.Find(devicetype.DeviceType{Name: name})
}

// Update the map of hashes to device types.
func (s *DeviceTypeService) updateDeviceTypeHashes() {
	deviceTypes, err := s.repository.GetAll()
	if err != nil {
		logrus.Warn(err)
		return
	}

	s.hashedDeviceTypes = make(map[string]string)

	table := crc16.MakeTable(crc16.CRC16_XMODEM)

	for _, deviceType := range deviceTypes {
		hash := crc16.Checksum([]byte(deviceType.Name), table)
		s.hashedDeviceTypes[fmt.Sprintf("%X", hash)] = deviceType.Name
	}
}

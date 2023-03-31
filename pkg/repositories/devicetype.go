package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
	"gorm.io/gorm"
)

type DeviceTypeRepository struct {
	db *gorm.DB
}

// Create a new DeviceTypeRepository.
func NewDeviceTypeRepository(db *gorm.DB) *DeviceTypeRepository {
	return &DeviceTypeRepository{
		db: db,
	}
}

// Database representation of a [twomes.DeviceType]
type DeviceTypeModel struct {
	gorm.Model
	Name                  string `gorm:"unique;non null"`
	InstallationManualURL string
	InfoURL               string
	Properties            []*PropertyModel `gorm:"many2many:device_type_properties"`
	UploadInterval        time.Duration
}

// Set the name of the table in the database.
func (DeviceTypeModel) TableName() string {
	return "device_types"
}

// Create a DeviceTypeModel from a [twomes.DeviceType].
func MakeDeviceTypeModel(deviceType twomes.DeviceType) DeviceTypeModel {
	var propertyModels []*PropertyModel

	for _, property := range deviceType.Properties {
		propertyModel := MakePropertyModel(property)
		propertyModels = append(propertyModels, &propertyModel)
	}

	return DeviceTypeModel{
		Model:                 gorm.Model{ID: deviceType.ID},
		Name:                  deviceType.Name,
		InstallationManualURL: deviceType.InstallationManualURL,
		InfoURL:               deviceType.InfoURL,
		Properties:            propertyModels,
		UploadInterval:        deviceType.UploadInterval.Duration,
	}
}

// Create a [twomes.DeviceType] from a DeviceTypeModel.
func (m *DeviceTypeModel) fromModel() twomes.DeviceType {
	var properties []twomes.Property

	for _, propertyModel := range m.Properties {
		properties = append(properties, propertyModel.fromModel())
	}

	return twomes.DeviceType{
		ID:                    m.Model.ID,
		Name:                  m.Name,
		InstallationManualURL: m.InstallationManualURL,
		InfoURL:               m.InfoURL,
		Properties:            properties,
		UploadInterval:        twomes.MakeDuration(m.UploadInterval),
	}
}

func (r *DeviceTypeRepository) Find(deviceType twomes.DeviceType) (twomes.DeviceType, error) {
	deviceTypeModel := MakeDeviceTypeModel(deviceType)
	err := r.db.Where(&deviceTypeModel).First(&deviceTypeModel).Error
	return deviceTypeModel.fromModel(), err
}

func (r *DeviceTypeRepository) GetAll() ([]twomes.DeviceType, error) {
	var deviceTypes []twomes.DeviceType

	var deviceTypeModels []DeviceTypeModel
	err := r.db.Find(&deviceTypeModels).Error
	if err != nil {
		return nil, err
	}

	for _, deviceTypeModel := range deviceTypeModels {
		deviceTypes = append(deviceTypes, deviceTypeModel.fromModel())
	}

	return deviceTypes, nil
}

func (r *DeviceTypeRepository) Create(deviceType twomes.DeviceType) (twomes.DeviceType, error) {
	deviceTypeModel := MakeDeviceTypeModel(deviceType)
	err := r.db.Create(&deviceTypeModel).Error
	return deviceTypeModel.fromModel(), err
}

func (r *DeviceTypeRepository) Delete(deviceType twomes.DeviceType) error {
	deviceTypeModel := MakeDeviceTypeModel(deviceType)
	return r.db.Delete(&deviceTypeModel).Error
}

package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

// Create a new DeviceRepository.
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

// Database representation of a [twomes.Device]
type DeviceModel struct {
	gorm.Model
	Name                 string `gorm:"unique;not null"`
	DeviceTypeModelID    uint   `gorm:"column:device_type_id"`
	DeviceType           DeviceTypeModel
	BuildingModelID      uint `gorm:"column:building_id"`
	ActivationSecretHash string
	ActivatedAt          *time.Time
	Uploads              []UploadModel
}

// Set the name of the table in the database.
func (DeviceModel) TableName() string {
	return "device"
}

// Create a DeviceModel from a [twomes.Device].
func MakeDeviceModel(device twomes.Device) DeviceModel {
	var uploadModels []UploadModel

	for _, upload := range device.Uploads {
		uploadModels = append(uploadModels, MakeUploadModel(upload))
	}

	return DeviceModel{
		Model:                gorm.Model{ID: device.ID},
		Name:                 device.Name,
		DeviceTypeModelID:    device.DeviceType.ID,
		DeviceType:           MakeDeviceTypeModel(device.DeviceType),
		BuildingModelID:      device.BuildingID,
		ActivationSecretHash: device.ActivationSecretHash,
		ActivatedAt:          device.ActivatedAt,
		Uploads:              uploadModels,
	}
}

// Create a [twomes.Device] from a DeviceModel.
func (m *DeviceModel) fromModel() twomes.Device {
	var uploads []twomes.Upload

	for _, uploadModel := range m.Uploads {
		uploads = append(uploads, uploadModel.fromModel())
	}

	return twomes.Device{
		ID:                   m.Model.ID,
		Name:                 m.Name,
		DeviceType:           m.DeviceType.fromModel(),
		BuildingID:           m.BuildingModelID,
		ActivationSecretHash: m.ActivationSecretHash,
		ActivatedAt:          m.ActivatedAt,
		Uploads:              uploads,
	}
}

func (r *DeviceRepository) Find(device twomes.Device) (twomes.Device, error) {
	deviceModel := MakeDeviceModel(device)
	err := r.db.Preload("DeviceType").Preload("Uploads").Where(&deviceModel).First(&deviceModel).Error
	return deviceModel.fromModel(), err
}

func (r *DeviceRepository) FindCloudFeedAuthCreationTimeFromDeviceID(deviceID uint) (*time.Time, error) {
	result := struct {
		CreatedAt time.Time
	}{}

	err := r.db.
		Table("device").
		Select("cloud_feed_auth.created_at").
		Joins("JOIN device_type ON device.device_type_id = device_type.id").
		Joins("JOIN cloud_feed ON device_type.name = cloud_feed.name").
		Joins("JOIN building ON device.building_id = building.id").
		Joins("JOIN account ON building.account_id = account.id").
		Joins("JOIN cloud_feed_auth ON account.id = cloud_feed_auth.account_id").
		Where("device.id = ?", deviceID).
		First(&result).
		Error

	if err != nil {
		return nil, err
	}

	return &result.CreatedAt, nil
}

func (r *DeviceRepository) GetMeasurements(device twomes.Device, filters map[string]string) ([]twomes.Measurement, error) {
	// empty array of measurements
	var measurements []twomes.Measurement = make([]twomes.Measurement, 0)

	query := r.db.
		Model(&twomes.Measurement{}).
		Preload("Property").
		Joins("JOIN upload ON measurement.upload_id = upload.id").
		Joins("JOIN device ON upload.device_id = device.id").
		Where("device.id = ?", device.ID)

	// apply filters
	for name, value := range filters {
		switch name {
		case "property":
			name = "property_id"
		case "start":
			name = "measurement.time >= ?"
		case "end":
			name = "measurement.time <= ?"
		}

		query = query.Where(name, value)
	}

	err := query.Find(&measurements).Error

	if err != nil {
		return nil, err
	}

	return measurements, nil
}

func (r *DeviceRepository) GetProperties(device twomes.Device) ([]twomes.Property, error) {
	var properties []twomes.Property = make([]twomes.Property, 0)

	err := r.db.
		Table("device").
		Select("DISTINCT property.id, property.name").
		Joins("JOIN upload ON device.id = upload.device_id").
		Joins("JOIN measurement ON upload.id = measurement.upload_id").
		Joins("JOIN property ON property.id = measurement.property_id").
		Where("device.id = ?", device.ID).
		Scan(&properties).
		Error

	if err != nil {
		return nil, err
	}

	return properties, nil
}

func (r *DeviceRepository) GetAll() ([]twomes.Device, error) {
	var devices []twomes.Device

	var deviceModels []DeviceModel
	err := r.db.Preload("DeviceType").Preload("Uploads").Find(&deviceModels).Error
	if err != nil {
		return nil, err
	}

	for _, deviceModel := range deviceModels {
		devices = append(devices, deviceModel.fromModel())
	}

	return devices, nil
}

func (r *DeviceRepository) Create(device twomes.Device) (twomes.Device, error) {
	deviceModel := MakeDeviceModel(device)
	err := r.db.Preload("").Create(&deviceModel).Error
	return deviceModel.fromModel(), err
}

func (r *DeviceRepository) Update(device twomes.Device) (twomes.Device, error) {
	deviceModel := MakeDeviceModel(device)
	err := r.db.Model(&deviceModel).Updates(deviceModel).Error
	return deviceModel.fromModel(), err
}

func (r *DeviceRepository) Delete(device twomes.Device) error {
	deviceModel := MakeDeviceModel(device)
	return r.db.Delete(&deviceModel).Error
}

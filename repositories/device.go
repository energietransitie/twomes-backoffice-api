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

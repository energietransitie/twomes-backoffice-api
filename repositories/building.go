package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/building"
	"github.com/energietransitie/twomes-backoffice-api/twomes/device"
	"gorm.io/gorm"
)

type BuildingRepository struct {
	db *gorm.DB
}

func NewBuildingRepository(db *gorm.DB) *BuildingRepository {
	return &BuildingRepository{
		db: db,
	}
}

// Database representation of a [building.Building].
type BuildingModel struct {
	gorm.Model
	AccountModelID  uint `gorm:"column:account_id"`
	Longitude       float32
	Latitude        float32
	TZName          string
	Devices         []DeviceModel
}

// Set the name of the table in the database.
func (BuildingModel) TableName() string {
	return "building"
}

// Create a new BuildingModel from a [building.Building]
func MakeBuildingModel(building building.Building) BuildingModel {
	var deviceModels []DeviceModel

	for _, device := range building.Devices {
		deviceModels = append(deviceModels, MakeDeviceModel(*device))
	}

	return BuildingModel{
		Model:          gorm.Model{ID: building.ID},
		AccountModelID: building.AccountID,
		Longitude:      building.Longitude,
		Latitude:       building.Latitude,
		TZName:         building.TZName,
		Devices:        deviceModels,
	}
}

// Create a [building.Building] from a BuildingModel.
func (m *BuildingModel) fromModel() building.Building {
	var devices []*device.Device

	for _, deviceModel := range m.Devices {
		device := deviceModel.fromModel()
		devices = append(devices, &device)
	}

	return building.Building{
		ID:        m.Model.ID,
		AccountID: m.AccountModelID,
		Longitude: m.Longitude,
		Latitude:  m.Latitude,
		TZName:    m.TZName,
		Devices:   devices,
	}
}

func (r *BuildingRepository) Find(building building.Building) (building.Building, error) {
	buildingModel := MakeBuildingModel(building)
	err := r.db.Preload("Devices.DeviceType").Where(&buildingModel).First(&buildingModel).Error
	return buildingModel.fromModel(), err
}

func (r *BuildingRepository) GetAll() ([]building.Building, error) {
	buildings := make([]building.Building, 0)

	var buildingModels []BuildingModel
	err := r.db.Find(&buildingModels).Error
	if err != nil {
		return nil, err
	}

	for _, buildingModel := range buildingModels {
		buildings = append(buildings, buildingModel.fromModel())
	}

	return buildings, nil
}

func (r *BuildingRepository) Create(building building.Building) (building.Building, error) {
	buildingModel := MakeBuildingModel(building)
	err := r.db.Create(&buildingModel).Error
	return buildingModel.fromModel(), err
}

func (r *BuildingRepository) Delete(building building.Building) error {
	buildingModel := MakeBuildingModel(building)
	return r.db.Delete(&buildingModel).Error
}

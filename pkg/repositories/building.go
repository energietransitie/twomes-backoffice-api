package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
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

// Database representation of a [twomes.Building].
type BuildingModel struct {
	gorm.Model
	AccountModelID uint `gorm:"column:account_id"`
	Longtitude     float32
	Latitude       float32
	TZName         string
	Devices        []DeviceModel
}

// Set the name of the table in the database.
func (BuildingModel) TableName() string {
	return "buildings"
}

// Create a new BuildingModel from a [twomes.Building]
func MakeBuildingModel(building twomes.Building) BuildingModel {
	var deviceModels []DeviceModel

	for _, device := range building.Devices {
		deviceModels = append(deviceModels, MakeDeviceModel(*device))
	}

	return BuildingModel{
		Model:          gorm.Model{ID: building.ID},
		AccountModelID: building.AccountID,
		Longtitude:     building.Longtitude,
		Latitude:       building.Latitude,
		TZName:         building.TZName,
		Devices:        deviceModels,
	}
}

// Create a [twomes.Building] from a BuildingModel.
func (m *BuildingModel) fromModel() twomes.Building {
	var devices []*twomes.Device

	for _, deviceModel := range m.Devices {
		device := deviceModel.fromModel()
		devices = append(devices, &device)
	}

	return twomes.Building{
		ID:         m.Model.ID,
		AccountID:  m.AccountModelID,
		Longtitude: m.Longtitude,
		Latitude:   m.Latitude,
		TZName:     m.TZName,
		Devices:    devices,
	}
}

func (r *BuildingRepository) Find(building twomes.Building) (twomes.Building, error) {
	buildingModel := MakeBuildingModel(building)
	err := r.db.Where(&buildingModel).First(&buildingModel).Error
	return buildingModel.fromModel(), err
}

func (r *BuildingRepository) GetAll() ([]twomes.Building, error) {
	buildings := make([]twomes.Building, 0)

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

func (r *BuildingRepository) Create(building twomes.Building) (twomes.Building, error) {
	buildingModel := MakeBuildingModel(building)
	err := r.db.Create(&buildingModel).Error
	return buildingModel.fromModel(), err
}

func (r *BuildingRepository) Delete(building twomes.Building) error {
	buildingModel := MakeBuildingModel(building)
	return r.db.Delete(&buildingModel).Error
}

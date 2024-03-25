package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryproperty"
	"gorm.io/gorm"
)

type EnergyQueryPropertyRepository struct {
	db *gorm.DB
}

// Create a new EnergyQueryPropertyRepository.
func NewEnergyQueryPropertyRepository(db *gorm.DB) *EnergyQueryPropertyRepository {
	return &EnergyQueryPropertyRepository{
		db: db,
	}
}

// Database representation of a [energyqueryproperty.EnergyQueryProperty]
type EnergyQueryPropertyModel struct {
	gorm.Model
	Name string `gorm:"unique;non null"`
	Unit string
}

// Set the name of the table in the database.
func (EnergyQueryPropertyModel) TableName() string {
	return "energy_query_property"
}

// Create a EnergyQueryPropertyModel from a [energyqueryproperty.EnergyQueryProperty].
func MakeEnergyQueryPropertyModel(energyQueryProperty energyqueryproperty.EnergyQueryProperty) EnergyQueryPropertyModel {
	return EnergyQueryPropertyModel{
		Model: gorm.Model{ID: energyQueryProperty.ID},
		Name:  energyQueryProperty.Name,
		Unit:  energyQueryProperty.Unit,
	}
}

// Create a [energyqueryproperty.EnergyQueryProperty] from a EnergyQueryPropertyModel.
func (m *EnergyQueryPropertyModel) fromModel() energyqueryproperty.EnergyQueryProperty {
	return energyqueryproperty.EnergyQueryProperty{
		ID:   m.Model.ID,
		Name: m.Name,
		Unit: m.Unit,
	}
}

func (r *EnergyQueryPropertyRepository) Find(energyQueryProperty energyqueryproperty.EnergyQueryProperty) (energyqueryproperty.EnergyQueryProperty, error) {
	energyQueryPropertyModel := MakeEnergyQueryPropertyModel(energyQueryProperty)
	err := r.db.Where(&energyQueryPropertyModel).First(&energyQueryPropertyModel).Error
	return energyQueryPropertyModel.fromModel(), err
}

func (r *EnergyQueryPropertyRepository) GetAll() ([]energyqueryproperty.EnergyQueryProperty, error) {
	var queryProperties []energyqueryproperty.EnergyQueryProperty

	var energyQueryPropertyModels []EnergyQueryPropertyModel
	err := r.db.Find(&energyQueryPropertyModels).Error
	if err != nil {
		return nil, err
	}

	for _, energyQueryPropertyModel := range energyQueryPropertyModels {
		queryProperties = append(queryProperties, energyQueryPropertyModel.fromModel())
	}

	return queryProperties, nil
}

func (r *EnergyQueryPropertyRepository) Create(energyQueryProperty energyqueryproperty.EnergyQueryProperty) (energyqueryproperty.EnergyQueryProperty, error) {
	energyQueryPropertyModel := MakeEnergyQueryPropertyModel(energyQueryProperty)
	err := r.db.Create(&energyQueryPropertyModel).Error
	return energyQueryPropertyModel.fromModel(), err
}

func (r *EnergyQueryPropertyRepository) Delete(energyQueryProperty energyqueryproperty.EnergyQueryProperty) error {
	energyQueryPropertyModel := MakeEnergyQueryPropertyModel(energyQueryProperty)
	return r.db.Delete(&energyQueryPropertyModel).Error
}

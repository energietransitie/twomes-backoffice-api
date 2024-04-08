package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvariety"
	"gorm.io/gorm"
)

type EnergyQueryVarietyRepository struct {
	db *gorm.DB
}

// Create a new EnergyQueryVarietyRepository.
func NewEnergyQueryVarietyRepository(db *gorm.DB) *EnergyQueryVarietyRepository {
	return &EnergyQueryVarietyRepository{
		db: db,
	}
}

// Database representation of a [energyqueryvariety.EnergyQueryVariety]
type EnergyQueryVarietyModel struct {
	gorm.Model
	Name             string                 `gorm:"unique;non null"`
	EnergyQueryTypes []EnergyQueryTypeModel `gorm:"foreignKey:EnergyQueryVarietyID"`
}

// Set the name of the table in the database.
func (EnergyQueryVarietyModel) TableName() string {
	return "energy_query_variety"
}

// Create a EnergyQueryVarietyModel from a [energyqueryvariety.EnergyQueryVariety].
func MakeEnergyQueryVarietyModel(energyQueryVariety energyqueryvariety.EnergyQueryVariety) EnergyQueryVarietyModel {
	return EnergyQueryVarietyModel{
		Model: gorm.Model{ID: energyQueryVariety.ID},
		Name:  energyQueryVariety.Name,
	}
}

// Create a [energyqueryvariety.EnergyQueryVariety] from a EnergyQueryVarietyModel.
func (m *EnergyQueryVarietyModel) fromModel() energyqueryvariety.EnergyQueryVariety {
	return energyqueryvariety.EnergyQueryVariety{
		ID:   m.Model.ID,
		Name: m.Name,
	}
}

func (r *EnergyQueryVarietyRepository) Find(energyQuery energyqueryvariety.EnergyQueryVariety) (energyqueryvariety.EnergyQueryVariety, error) {
	energyQueryModel := MakeEnergyQueryVarietyModel(energyQuery)
	err := r.db.Where(&energyQueryModel).First(&energyQueryModel).Error
	return energyQueryModel.fromModel(), err
}

func (r *EnergyQueryVarietyRepository) GetAll() ([]energyqueryvariety.EnergyQueryVariety, error) {
	var energyQueries []energyqueryvariety.EnergyQueryVariety

	var energyQueryModels []EnergyQueryVarietyModel
	err := r.db.Find(&energyQueryModels).Error
	if err != nil {
		return nil, err
	}

	for _, energyQueryModel := range energyQueryModels {
		energyQueries = append(energyQueries, energyQueryModel.fromModel())
	}

	return energyQueries, nil
}

func (r *EnergyQueryVarietyRepository) Create(energyQuery energyqueryvariety.EnergyQueryVariety) (energyqueryvariety.EnergyQueryVariety, error) {
	energyQueryModel := MakeEnergyQueryVarietyModel(energyQuery)
	err := r.db.Create(&energyQueryModel).Error
	return energyQueryModel.fromModel(), err
}

func (r *EnergyQueryVarietyRepository) Delete(energyQuery energyqueryvariety.EnergyQueryVariety) error {
	energyQueryModel := MakeEnergyQueryVarietyModel(energyQuery)
	return r.db.Delete(&energyQueryModel).Error
}

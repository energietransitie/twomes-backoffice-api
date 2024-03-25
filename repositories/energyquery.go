package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquery"
	"gorm.io/gorm"
)

type EnergyQueryRepository struct {
	db *gorm.DB
}

// Create a new EnergyQueryRepository.
func NewEnergyQueryRepository(db *gorm.DB) *EnergyQueryRepository {
	return &EnergyQueryRepository{
		db: db,
	}
}

// Database representation of a [energyquery.EnergyQuery]
type EnergyQueryModel struct {
	gorm.Model
	Name    string `gorm:"unique;non null"`
	Formula string
}

// Set the name of the table in the database.
func (EnergyQueryModel) TableName() string {
	return "energy_query"
}

// Create a EnergyQueryModel from a [energyquery.EnergyQuery].
func MakeEnergyQueryModel(energyQuery energyquery.EnergyQuery) EnergyQueryModel {
	return EnergyQueryModel{
		Model:   gorm.Model{ID: energyQuery.ID},
		Name:    energyQuery.Name,
		Formula: energyQuery.Formula,
	}
}

// Create a [energyquery.EnergyQuery] from a EnergyQueryModel.
func (m *EnergyQueryModel) fromModel() energyquery.EnergyQuery {
	return energyquery.EnergyQuery{
		ID:      m.Model.ID,
		Name:    m.Name,
		Formula: m.Formula,
	}
}

func (r *EnergyQueryRepository) Find(energyQuery energyquery.EnergyQuery) (energyquery.EnergyQuery, error) {
	energyQueryModel := MakeEnergyQueryModel(energyQuery)
	err := r.db.Where(&energyQueryModel).First(&energyQueryModel).Error
	return energyQueryModel.fromModel(), err
}

func (r *EnergyQueryRepository) GetAll() ([]energyquery.EnergyQuery, error) {
	var energyQueries []energyquery.EnergyQuery

	var energyQueryModels []EnergyQueryModel
	err := r.db.Find(&energyQueryModels).Error
	if err != nil {
		return nil, err
	}

	for _, energyQueryModel := range energyQueryModels {
		energyQueries = append(energyQueries, energyQueryModel.fromModel())
	}

	return energyQueries, nil
}

func (r *EnergyQueryRepository) Create(energyQuery energyquery.EnergyQuery) (energyquery.EnergyQuery, error) {
	energyQueryModel := MakeEnergyQueryModel(energyQuery)
	err := r.db.Create(&energyQueryModel).Error
	return energyQueryModel.fromModel(), err
}

func (r *EnergyQueryRepository) Delete(energyQuery energyquery.EnergyQuery) error {
	energyQueryModel := MakeEnergyQueryModel(energyQuery)
	return r.db.Delete(&energyQueryModel).Error
}

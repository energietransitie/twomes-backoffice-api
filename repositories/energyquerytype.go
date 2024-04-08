package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquerytype"
	"gorm.io/gorm"
)

type EnergyQueryTypeRepository struct {
	db *gorm.DB
}

// Create a new EnergyQueryTypeRepository.
func NewEnergyQueryTypeRepository(db *gorm.DB) *EnergyQueryTypeRepository {
	return &EnergyQueryTypeRepository{
		db: db,
	}
}

// Database representation of a [energyquerytype.EnergyQueryType]
type EnergyQueryTypeModel struct {
	gorm.Model
	EnergyQueryVarietyID uint
	Formula              string
	EnergyQueries        []EnergyQueryModel `gorm:"foreignKey:EnergyQueryTypeID"`
}

// Set the name of the table in the database.
func (EnergyQueryTypeModel) TableName() string {
	return "energy_query_type"
}

// Create a EnergyQueryTypeModel from a [energyquerytype.EnergyQueryType].
func MakeEnergyQueryTypeModel(energyQueryTypeType energyquerytype.EnergyQueryType) EnergyQueryTypeModel {
	return EnergyQueryTypeModel{
		Model:   gorm.Model{ID: energyQueryTypeType.ID},
		Formula: energyQueryTypeType.Formula,
	}
}

// Create a [energyquerytype.EnergyQueryType] from a EnergyQueryTypeModel.
func (m *EnergyQueryTypeModel) fromModel() energyquerytype.EnergyQueryType {
	return energyquerytype.EnergyQueryType{
		ID:      m.Model.ID,
		Formula: m.Formula,
	}
}

func (r *EnergyQueryTypeRepository) Find(energyQueryType energyquerytype.EnergyQueryType) (energyquerytype.EnergyQueryType, error) {
	energyQueryTypeModel := MakeEnergyQueryTypeModel(energyQueryType)
	err := r.db.Where(&energyQueryTypeModel).Preload("EnergyQueryVariety").First(&energyQueryTypeModel).Error
	return energyQueryTypeModel.fromModel(), err
}

func (r *EnergyQueryTypeRepository) GetAll() ([]energyquerytype.EnergyQueryType, error) {
	var energyQueryTypes []energyquerytype.EnergyQueryType

	var energyQueryTypeModels []EnergyQueryTypeModel
	err := r.db.Find(&energyQueryTypeModels).Error
	if err != nil {
		return nil, err
	}

	for _, energyQueryTypeModel := range energyQueryTypeModels {
		energyQueryTypes = append(energyQueryTypes, energyQueryTypeModel.fromModel())
	}

	return energyQueryTypes, nil
}

func (r *EnergyQueryTypeRepository) Create(energyQueryType energyquerytype.EnergyQueryType) (energyquerytype.EnergyQueryType, error) {
	energyQueryTypeModel := MakeEnergyQueryTypeModel(energyQueryType)
	err := r.db.Create(&energyQueryTypeModel).Error
	return energyQueryTypeModel.fromModel(), err
}

func (r *EnergyQueryTypeRepository) Delete(energyQueryType energyquerytype.EnergyQueryType) error {
	energyQueryTypeModel := MakeEnergyQueryTypeModel(energyQueryType)
	return r.db.Delete(&energyQueryTypeModel).Error
}

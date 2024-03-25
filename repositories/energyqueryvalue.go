package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvalue"
	"gorm.io/gorm"
)

// Database representation of a [EnergyQueryValue.EnergyQueryValue]
type EnergyQueryValueModel struct {
	gorm.Model
	QueryUploadID   uint `gorm:"column:query_upload_id"`
	QueryPropertyID uint `gorm:"column:query_property_id"`
	QueryProperty   EnergyQueryPropertyModel
	Value           string 
}

// Set the name of the table in the database.
func (EnergyQueryValueModel) TableName() string {
	return "energy_query_value"
}

// Create a EnergyQueryValueModel from a [energyqueryvalue.EnergyQueryValue].
func MakeEnergyQueryValueModel(energyQueryValue energyqueryvalue.EnergyQueryValue) EnergyQueryValueModel {
	return EnergyQueryValueModel{
		Model:           gorm.Model{ID: energyQueryValue.ID},
		QueryUploadID:   energyQueryValue.QueryUploadID,
		QueryPropertyID: energyQueryValue.Property.ID,
		QueryProperty:   MakeEnergyQueryPropertyModel(energyQueryValue.Property),
		Value:           energyQueryValue.Value,
	}
}

// Create a [energyqueryvalue.EnergyQueryValue] from a EnergyQueryValueModel.
func (m *EnergyQueryValueModel) fromModel() energyqueryvalue.EnergyQueryValue {
	return energyqueryvalue.EnergyQueryValue{
		ID:            m.Model.ID,
		QueryUploadID: m.QueryUploadID,
		PropertyID:    int(m.QueryProperty.ID),
		Property:      m.QueryProperty.fromModel(),
		Value:         m.Value,
	}
}

package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
	"gorm.io/gorm"
)

// Database representation of a [twomes.Measurement]
type MeasurementModel struct {
	gorm.Model
	PropertyModelID uint `gorm:"column:property_id"`
	Property        PropertyModel
	UploadModelID   uint `gorm:"column:upload_id"`
	Time            time.Time
	Value           string
}

// Set the name of the table in the database.
func (MeasurementModel) TableName() string {
	return "measurements"
}

// Create a MeasurementModel from a [twomes.Measurement].
func MakeMeasurementModel(measurement twomes.Measurement) MeasurementModel {
	return MeasurementModel{
		Model:           gorm.Model{ID: measurement.ID},
		PropertyModelID: measurement.Property.ID,
		Property:        MakePropertyModel(measurement.Property),
		UploadModelID:   measurement.UploadID,
		Time:            measurement.Time,
		Value:           measurement.Value,
	}
}

// Create a [twomes.Measurement] from a MeasurementModel.
func (m *MeasurementModel) fromModel() twomes.Measurement {
	return twomes.Measurement{
		ID:       m.Model.ID,
		UploadID: m.UploadModelID,
		Property: m.Property.fromModel(),
		Time:     m.Time,
		Value:    m.Value,
	}
}

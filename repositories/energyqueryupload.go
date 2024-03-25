package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryupload"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvalue"
	"gorm.io/gorm"
)

type EnergyQueryUploadRepository struct {
	db *gorm.DB
}

// Create a new EnergyQueryUploadRepository.
func NewEnergyQueryUploadRepository(db *gorm.DB) *EnergyQueryUploadRepository {
	return &EnergyQueryUploadRepository{
		db: db,
	}
}

// Database representation of a [energyqueryupload.EnergyQueryUpload]
type EnergyQueryUploadModel struct {
	gorm.Model
	QueryID    uint `gorm:"column:query_id"`
	BuildingID uint `gorm:"column:building_id"`
	Size       int
	Values     []EnergyQueryValueModel `gorm:"foreignKey:query_upload_id"`
}

// Set the name of the table in the database.
func (EnergyQueryUploadModel) TableName() string {
	return "energy_query_upload"
}

// Create an EnergyQueryUploadModel from a [energyqueryupload.EnergyQueryUpload].
func MakeEnergyQueryUploadModel(energyQueryUpload energyqueryupload.EnergyQueryUpload) EnergyQueryUploadModel {
	var valueModels []EnergyQueryValueModel

	for _, value := range energyQueryUpload.EnergyQueryValues {
		valueModels = append(valueModels, MakeEnergyQueryValueModel(value))
	}

	return EnergyQueryUploadModel{
		Model:      gorm.Model{ID: energyQueryUpload.ID},
		QueryID:    energyQueryUpload.QueryID,
		BuildingID: energyQueryUpload.BuildingID,
		Size:       energyQueryUpload.Size,
		Values:     valueModels,
	}
}

// Create a [energyqueryupload.EnergyQueryUpload] from an EnergyQueryUploadModel.
func (m *EnergyQueryUploadModel) fromModel() energyqueryupload.EnergyQueryUpload {
	var values []energyqueryvalue.EnergyQueryValue

	for _, valueModel := range m.Values {
		values = append(values, valueModel.fromModel())
	}

	return energyqueryupload.EnergyQueryUpload{
		ID:                m.Model.ID,
		QueryID:           m.QueryID,
		BuildingID:        m.BuildingID,
		Size:              m.Size,
		EnergyQueryValues: values,
	}
}

func (r *EnergyQueryUploadRepository) Find(energyQueryUpload energyqueryupload.EnergyQueryUpload) (energyqueryupload.EnergyQueryUpload, error) {
	energyQueryUploadModel := MakeEnergyQueryUploadModel(energyQueryUpload)
	err := r.db.Preload("Energy_Query_Values").Where(&energyQueryUploadModel).Find(&energyQueryUploadModel).Error
	return energyQueryUploadModel.fromModel(), err
}

func (r *EnergyQueryUploadRepository) GetAll() ([]energyqueryupload.EnergyQueryUpload, error) {
	var energyQueryUploads []energyqueryupload.EnergyQueryUpload

	var energyQueryUploadModels []EnergyQueryUploadModel
	err := r.db.Preload("Energy_Query_Values").Find(&energyQueryUploadModels).Error
	if err != nil {
		return nil, err
	}

	for _, energyQueryUploadModel := range energyQueryUploadModels {
		energyQueryUploads = append(energyQueryUploads, energyQueryUploadModel.fromModel())
	}

	return energyQueryUploads, nil
}

func (r *EnergyQueryUploadRepository) Create(energyQueryUpload energyqueryupload.EnergyQueryUpload) (energyqueryupload.EnergyQueryUpload, error) {
	energyQueryUploadModel := MakeEnergyQueryUploadModel(energyQueryUpload)
	err := r.db.Create(&energyQueryUploadModel).Error
	return energyQueryUploadModel.fromModel(), err
}

func (r *EnergyQueryUploadRepository) Delete(energyQueryUpload energyqueryupload.EnergyQueryUpload) error {
	energyQueryUploadModel := MakeEnergyQueryUploadModel(energyQueryUpload)
	return r.db.Delete(&energyQueryUploadModel).Error
}

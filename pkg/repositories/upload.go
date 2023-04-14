package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
	"gorm.io/gorm"
)

type UploadRepository struct {
	db *gorm.DB
}

// Create a new UploadRepository.
func NewUploadRepository(db *gorm.DB) *UploadRepository {
	return &UploadRepository{
		db: db,
	}
}

// Database representation of a [twomes.Upload]
type UploadModel struct {
	gorm.Model
	DeviceModelID uint `gorm:"column:device_id"`
	ServerTime    time.Time
	DeviceTime    time.Time
	Size          int
	Measurements  []MeasurementModel
}

// Set the name of the table in the database.
func (UploadModel) TableName() string {
	return "upload"
}

// Create an UploadModel from a [twomes.Upload].
func MakeUploadModel(upload twomes.Upload) UploadModel {
	var measurementModels []MeasurementModel

	for _, measurement := range upload.Measurements {
		measurementModels = append(measurementModels, MakeMeasurementModel(measurement))
	}

	return UploadModel{
		Model:         gorm.Model{ID: upload.ID},
		DeviceModelID: upload.DeviceID,
		ServerTime:    time.Time(upload.ServerTime),
		DeviceTime:    time.Time(upload.DeviceTime),
		Size:          upload.Size,
		Measurements:  measurementModels,
	}
}

// Create a [twomes.Upload] from an UploadModel.
func (m *UploadModel) fromModel() twomes.Upload {
	var measurements []twomes.Measurement

	for _, measurementModel := range m.Measurements {
		measurements = append(measurements, measurementModel.fromModel())
	}

	return twomes.Upload{
		ID:           m.Model.ID,
		DeviceID:     m.DeviceModelID,
		ServerTime:   twomes.Time(m.ServerTime),
		DeviceTime:   twomes.Time(m.DeviceTime),
		Size:         m.Size,
		Measurements: measurements,
	}
}

func (r *UploadRepository) Find(upload twomes.Upload) (twomes.Upload, error) {
	uploadModel := MakeUploadModel(upload)
	err := r.db.Preload("Measurements").Where(&uploadModel).Find(&uploadModel).Error
	return uploadModel.fromModel(), err
}

func (r *UploadRepository) GetAll() ([]twomes.Upload, error) {
	var uploads []twomes.Upload

	var uploadModels []UploadModel
	err := r.db.Preload("Measurements").Find(&uploadModels).Error
	if err != nil {
		return nil, err
	}

	for _, uploadModel := range uploadModels {
		uploads = append(uploads, uploadModel.fromModel())
	}

	return uploads, nil
}

func (r *UploadRepository) Create(upload twomes.Upload) (twomes.Upload, error) {
	uploadModel := MakeUploadModel(upload)
	err := r.db.Create(&uploadModel).Error
	return uploadModel.fromModel(), err
}

func (r *UploadRepository) Delete(upload twomes.Upload) error {
	uploadModel := MakeUploadModel(upload)
	return r.db.Delete(&uploadModel).Error
}

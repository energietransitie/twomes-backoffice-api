package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/energietransitie/twomes-backoffice-api/twomes/measurement"
	"github.com/energietransitie/twomes-backoffice-api/twomes/upload"
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

// Database representation of a [upload.Upload]
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

// Create an UploadModel from a [upload.Upload].
func MakeUploadModel(upload upload.Upload) UploadModel {
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

// Create a [upload.Upload] from an UploadModel.
func (m *UploadModel) fromModel() upload.Upload {
	var measurements []measurement.Measurement

	for _, measurementModel := range m.Measurements {
		measurements = append(measurements, measurementModel.fromModel())
	}

	return upload.Upload{
		ID:           m.Model.ID,
		DeviceID:     m.DeviceModelID,
		ServerTime:   twomes.Time(m.ServerTime),
		DeviceTime:   twomes.Time(m.DeviceTime),
		Size:         m.Size,
		Measurements: measurements,
	}
}

func (r *UploadRepository) Find(upload upload.Upload) (upload.Upload, error) {
	uploadModel := MakeUploadModel(upload)
	err := r.db.Preload("Measurements").Where(&uploadModel).Find(&uploadModel).Error
	return uploadModel.fromModel(), err
}

func (r *UploadRepository) GetAll() ([]upload.Upload, error) {
	var uploads []upload.Upload

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

func (r *UploadRepository) Create(upload upload.Upload) (upload.Upload, error) {
	uploadModel := MakeUploadModel(upload)
	err := r.db.Create(&uploadModel).Error
	return uploadModel.fromModel(), err
}

func (r *UploadRepository) Delete(upload upload.Upload) error {
	uploadModel := MakeUploadModel(upload)
	return r.db.Delete(&uploadModel).Error
}

func (r *UploadRepository) GetLatestUploadForDeviceWithID(id uint) (upload.Upload, error) {
	var uploadModel UploadModel
	err := r.db.Where(UploadModel{DeviceModelID: id}).Order("server_time desc").First(&uploadModel).Error
	return uploadModel.fromModel(), err
}

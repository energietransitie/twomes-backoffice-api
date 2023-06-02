package services

import (
	"errors"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

var (
	ErrEmptyUpload = errors.New("no measurements in upload")
)

type UploadService struct {
	repository ports.UploadRepository

	// Service used when creating an upload.
	propertyService ports.PropertyService
}

// Create a new UploadService.
func NewUploadService(repository ports.UploadRepository, propertyService ports.PropertyService) *UploadService {
	return &UploadService{
		repository:      repository,
		propertyService: propertyService,
	}
}

func (s *UploadService) Create(deviceID uint, deviceTime twomes.Time, measurements []twomes.Measurement) (twomes.Upload, error) {
	if len(measurements) <= 0 {
		return twomes.Upload{}, ErrEmptyUpload
	}

	upload := twomes.MakeUpload(deviceID, deviceTime, measurements)

	upload, err := s.repository.Create(upload)

	return upload, err
}

func (s *UploadService) GetLatestUploadTimeForDeviceWithID(id uint) (*time.Time, error) {
	upload, err := s.repository.GetLatestUploadForDeviceWithID(id)
	if err != nil {
		// If the record is not found, there was no upload. That's not an error.
		if helpers.IsMySQLRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return (*time.Time)(&upload.ServerTime), nil
}

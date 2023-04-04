package services

import (
	"errors"

	"github.com/energietransitie/twomes-backoffice-api/pkg/ports"
	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
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

package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

// An UploadRepository can load, store and delete uploads.
type UploadRepository interface {
	Find(Upload twomes.Upload) (twomes.Upload, error)
	GetAll() ([]twomes.Upload, error)
	Create(twomes.Upload) (twomes.Upload, error)
	Delete(twomes.Upload) error
	GetLatestUploadForDeviceWithID(id uint) (twomes.Upload, error)
}

// UploadService exposes all operations that can be performed on a [twomes.Upload].
type UploadService interface {
	Create(deviceID uint, deviceTime twomes.Time, measurements []twomes.Measurement) (twomes.Upload, error)
	GetLatestUploadTimeForDeviceWithID(id uint) (*time.Time, error)
}

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
	// GetLatestUploadTimeForDeviceWithID returns the latest upload time for a device.
	// If there is no upload, it returns the creation time of the cloud feed auth.
	// The bool is true if the time actually came from the latest upload.
	GetLatestUploadTimeForDeviceWithID(id uint) (*time.Time, bool, error)
}

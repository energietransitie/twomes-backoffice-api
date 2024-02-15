package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/energietransitie/twomes-backoffice-api/twomes/measurement"
	"github.com/energietransitie/twomes-backoffice-api/twomes/upload"
)

// UploadService exposes all operations that can be performed on a [upload.Upload].
type UploadService interface {
	Create(deviceID uint, deviceTime twomes.Time, measurements []measurement.Measurement) (upload.Upload, error)
	// GetLatestUploadTimeForDeviceWithID returns the latest upload time for a device.
	// If there is no upload, it returns the creation time of the cloud feed auth.
	// The bool is true if the time actually came from the latest upload.
	GetLatestUploadTimeForDeviceWithID(id uint) (*time.Time, bool, error)
}

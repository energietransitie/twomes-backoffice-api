package upload

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/energietransitie/twomes-backoffice-api/twomes/measurement"
)

// An Upload is a collection of measurements, with additional information.
type Upload struct {
	ID           uint                      `json:"id"`
	InstanceID   uint                      `json:"instance_id"`
	ServerTime   twomes.Time               `json:"server_time"`
	DeviceTime   twomes.Time               `json:"device_time"`
	Size         int                       `json:"size"`
	Measurements []measurement.Measurement `json:"measurements,omitempty"`
}

// Create a new Upload.
func MakeUpload(instanceID uint, deviceTime twomes.Time, measurements []measurement.Measurement) Upload {
	return Upload{
		InstanceID:   instanceID,
		ServerTime:   twomes.Time(time.Now().UTC()),
		DeviceTime:   deviceTime,
		Size:         len(measurements),
		Measurements: measurements,
	}
}

package upload

import (
	"time"

	"github.com/energietransitie/needforheat-server-api/needforheat"
	"github.com/energietransitie/needforheat-server-api/needforheat/measurement"
)

// An Upload is a collection of measurements, with additional information.
type Upload struct {
	ID           uint                      `json:"id"`
	InstanceID   uint                      `json:"instance_id"`
	ServerTime   needforheat.Time          `json:"server_time"`
	DeviceTime   needforheat.Time          `json:"device_time"`
	Size         int                       `json:"size"`
	Measurements []measurement.Measurement `json:"measurements,omitempty"`
}

// Create a new Upload.
func MakeUpload(instanceID uint, deviceTime needforheat.Time, measurements []measurement.Measurement) Upload {
	return Upload{
		InstanceID:   instanceID,
		ServerTime:   needforheat.Time(time.Now().UTC()),
		DeviceTime:   deviceTime,
		Size:         len(measurements),
		Measurements: measurements,
	}
}

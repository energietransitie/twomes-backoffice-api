package measurement

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/energietransitie/twomes-backoffice-api/twomes/property"
)

// A Measurement is a measured value for a specific property.
type Measurement struct {
	ID         uint              `json:"id"`
	UploadID   uint              `json:"upload_id"`
	PropertyID int               `json:"-"`
	Property   property.Property `json:"property"`
	Time       twomes.Time       `json:"time"`
	Value      string            `json:"value"`
}

package energyqueryvalue

import "github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryproperty"

// An EnergyQueryValue is a value of data we get from the server
type EnergyQueryValue struct {
	ID            uint                                    `json:"id"`
	QueryUploadID uint                                    `json:"query_upload_id"`
	PropertyID    int                                     `json:"-"`
	Property      energyqueryproperty.EnergyQueryProperty `json:"property"`
	Value         string                                  `json:"value"`
}

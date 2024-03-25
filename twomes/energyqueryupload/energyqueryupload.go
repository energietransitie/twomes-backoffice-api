package energyqueryupload

import "github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvalue"

// An EnergyQueryUpload is a collection of extra information we get from the user, tied to a building
type EnergyQueryUpload struct {
	ID                uint                                `json:"id"`
	QueryID           uint                                `json:"query_id"`
	BuildingID        uint                                `json:"building_id"`
	Size              int                                 `json:"size"`
	EnergyQueryValues []energyqueryvalue.EnergyQueryValue `json:"values,omitempty"`
}

// Create a new EnergyQueryUpload
func MakeEnergyQueryUpload(queryID uint, buildingID uint, energyQueryValues []energyqueryvalue.EnergyQueryValue) EnergyQueryUpload {
	return EnergyQueryUpload{
		QueryID:           queryID,
		BuildingID:        buildingID,
		Size:              len(energyQueryValues),
		EnergyQueryValues: energyQueryValues,
	}
}

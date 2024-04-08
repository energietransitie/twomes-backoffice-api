package energyquerytype

import "github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvariety"

// An EnergyQueryType is a datasourceitem in the datasourcelist
type EnergyQueryType struct {
	ID           uint                                  `json:"id"`
	QueryVariety energyqueryvariety.EnergyQueryVariety `json:"energy_query_variety_id"`
	Formula      string                                `json:"formula"` //Can be sent to the app for a calculation?
}

// Create a new EnergyQueryType.
func MakeEnergyQueryType(queryVariety energyqueryvariety.EnergyQueryVariety, formula string) EnergyQueryType {
	return EnergyQueryType{
		QueryVariety: queryVariety,
		Formula:      formula,
	}
}

package energyquerytype

// A EnergyQueryType contains information about a group of energyqueries with the same functionality.
type EnergyQueryType struct {
	ID                 uint   `json:"id"`
	EnergyQueryVariety string `json:"energy_query_variety"`
	Formula            string `json:"formula"`
}

// Create a new EnergyQueryType.
func MakeEnergyQueryType(energyQueryVariety string, formula string) EnergyQueryType {
	return EnergyQueryType{
		EnergyQueryVariety: energyQueryVariety,
		Formula:            formula,
	}
}

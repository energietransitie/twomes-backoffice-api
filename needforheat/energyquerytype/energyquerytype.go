package energyquerytype

// A EnergyQueryType contains information about a group of energyqueries with the same functionality.
type EnergyQueryType struct {
	ID      uint   `json:"id"`
	Name    string `json:"energy_query_variety"`
	Formula string `json:"formula"`
}

// Create a new EnergyQueryType.
func MakeEnergyQueryType(energyQueryVariety string, formula string) EnergyQueryType {
	return EnergyQueryType{
		Name:    energyQueryVariety,
		Formula: formula,
	}
}

package energyqueryvariety

// An EnergyQueryVariety contains only the name of a query that SHOULD BE IN THE APP
type EnergyQueryVariety struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Create a new EnergyQueryVariety.
func MakeEnergyQueryVariety(name string) EnergyQueryVariety {
	return EnergyQueryVariety{
		Name: name,
	}
}

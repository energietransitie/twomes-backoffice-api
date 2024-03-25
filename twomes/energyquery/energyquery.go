package energyquery

// An EnergyQuery contains only the name and (nullable) formula for the app
type EnergyQuery struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Formula string `json:"formula"`
}

// Create a new EnergyQuery.
func MakeEnergyQuery(name string, formula string) EnergyQuery {
	return EnergyQuery{
		Name:    name,
		Formula: formula,
	}
}

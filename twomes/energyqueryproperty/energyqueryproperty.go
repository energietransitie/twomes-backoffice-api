package energyqueryproperty

// A EnergyQueryProperty is the type of data sent by a query.
type EnergyQueryProperty struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

// Create a new EnergyQueryProperty.
func MakeEnergyQueryProperty(name string, unit string) EnergyQueryProperty {
	return EnergyQueryProperty{
		Name: name,
		Unit: unit,
	}
}

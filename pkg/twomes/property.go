package twomes

// A Property is the type of a measurement.
type Property struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

// Create a new Property.
func MakeProperty(name, unit string) Property {
	return Property{
		Name: name,
		Unit: unit,
	}
}

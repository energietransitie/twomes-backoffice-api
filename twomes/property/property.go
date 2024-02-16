package property

// A Property is the type of a measurement.
type Property struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Create a new Property.
func MakeProperty(name string) Property {
	return Property{
		Name: name,
	}
}

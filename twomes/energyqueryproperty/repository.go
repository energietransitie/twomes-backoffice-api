package energyqueryproperty

// A DeviceTypeRepository can load, store and delete properties.
type EnergyQueryPropertyRepository interface {
	Find(energyQueryProperty EnergyQueryProperty) (EnergyQueryProperty, error)
	GetAll() ([]EnergyQueryProperty, error)
	Create(EnergyQueryProperty) (EnergyQueryProperty, error)
	Delete(EnergyQueryProperty) error
}

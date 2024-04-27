package energyquerytype

// A EnergyQueryTypeRepository can load, store and delete device types.
type EnergyQueryTypeRepository interface {
	Find(energyQueryType EnergyQueryType) (EnergyQueryType, error)
	GetAll() ([]EnergyQueryType, error)
	Create(EnergyQueryType) (EnergyQueryType, error)
	Delete(EnergyQueryType) error
}

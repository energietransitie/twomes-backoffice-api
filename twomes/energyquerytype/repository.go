package energyquerytype

// A EnergyQueryTypeRepository can load, store and delete query types.
type EnergyQueryTypeRepository interface {
	Find(energyQueryVariety EnergyQueryType) (EnergyQueryType, error)
	GetAll() ([]EnergyQueryType, error)
	Create(EnergyQueryType) (EnergyQueryType, error)
	Delete(EnergyQueryType) error
}

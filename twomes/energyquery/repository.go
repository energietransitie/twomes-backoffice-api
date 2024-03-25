package energyquery

// A EnergyQueryRepository can load, store and delete device types.
type EnergyQueryRepository interface {
	Find(energyQuery EnergyQuery) (EnergyQuery, error)
	GetAll() ([]EnergyQuery, error)
	Create(EnergyQuery) (EnergyQuery, error)
	Delete(EnergyQuery) error
}

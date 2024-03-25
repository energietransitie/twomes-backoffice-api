package energyquery

// A EnergyQueryRepository can load, store and delete queries.
type EnergyQueryRepository interface {
	Find(energyQuery EnergyQuery) (EnergyQuery, error)
	GetAll() ([]EnergyQuery, error)
	Create(EnergyQuery) (EnergyQuery, error)
	Delete(EnergyQuery) error
}

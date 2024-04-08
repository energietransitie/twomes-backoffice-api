package energyqueryvariety

// A EnergyQueryVarietyRepository can load, store and delete queries.
type EnergyQueryVarietyRepository interface {
	Find(energyQueryVariety EnergyQueryVariety) (EnergyQueryVariety, error)
	GetAll() ([]EnergyQueryVariety, error)
	Create(EnergyQueryVariety) (EnergyQueryVariety, error)
	Delete(EnergyQueryVariety) error
}

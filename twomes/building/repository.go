package building

// A BuildingRepository can load, store and delete buildings.
type BuildingRepository interface {
	Find(building Building) (Building, error)
	GetAll() ([]Building, error)
	Create(Building) (Building, error)
	Delete(Building) error
}

package app

// An AppRepository can load, store and delete apps.
type AppRepository interface {
	Find(app App) (App, error)
	GetAll() ([]App, error)
	Create(App) (App, error)
	Delete(App) error
}

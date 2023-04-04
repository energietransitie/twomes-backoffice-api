package ports

import (
	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
)

// An AppRepository can load, store and delete apps.
type AppRepository interface {
	Find(app twomes.App) (twomes.App, error)
	GetAll() ([]twomes.App, error)
	Create(twomes.App) (twomes.App, error)
	Delete(twomes.App) error
}

// AppService exposes all operations that can be performed on a [twomes.App]
type AppService interface {
	Create(name, provisioningURLTemplate string) (twomes.App, error)
	Find(app twomes.App) (twomes.App, error)
	GetAll() ([]twomes.App, error)
	GetByID(id uint) (twomes.App, error)
}

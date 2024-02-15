package ports

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/app"
)

// AppService exposes all operations that can be performed on a [app.App]
type AppService interface {
	Create(name, provisioningURLTemplate, oauthRedirectURL string) (app.App, error)
	Find(app app.App) (app.App, error)
	GetAll() ([]app.App, error)
	GetByID(id uint) (app.App, error)
}

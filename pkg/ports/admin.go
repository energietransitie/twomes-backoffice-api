package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
)

// An AccountRepository can load, store and delete admins.
type AdminRepository interface {
	Find(admin twomes.Admin) (twomes.Admin, error)
	GetAll() ([]twomes.Admin, error)
	Create(admin twomes.Admin) (twomes.Admin, error)
	Update(admin twomes.Admin) (twomes.Admin, error)
	Delete(admin twomes.Admin) error
}

// AdminService exposes all operations that can be performed on a [twomes.Admin]
type AdminService interface {
	Create(name string, expiry time.Time) (twomes.Admin, error)
	Find(admin twomes.Admin) (twomes.Admin, error)
	GetAll() ([]twomes.Admin, error)
	Delete(admin twomes.Admin) error
	Reactivate(admin twomes.Admin) (twomes.Admin, error)
	SetExpiry(admin twomes.Admin, expiry time.Time) (twomes.Admin, error)
}

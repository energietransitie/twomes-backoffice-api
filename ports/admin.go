package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/admin"
)

// AdminService exposes all operations that can be performed on a [admin.Admin]
type AdminService interface {
	Create(name string, expiry time.Time) (admin.Admin, error)
	Find(admin admin.Admin) (admin.Admin, error)
	GetAll() ([]admin.Admin, error)
	Delete(admin admin.Admin) error
	Reactivate(admin admin.Admin) (admin.Admin, error)
	SetExpiry(admin admin.Admin, expiry time.Time) (admin.Admin, error)
}

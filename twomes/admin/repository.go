package admin

// An AccountRepository can load, store and delete admins.
type AdminRepository interface {
	Find(admin Admin) (Admin, error)
	GetAll() ([]Admin, error)
	Create(admin Admin) (Admin, error)
	Update(admin Admin) (Admin, error)
	Delete(admin Admin) error
}

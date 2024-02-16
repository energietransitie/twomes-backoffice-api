package account

// An AccountRepository can load, store and delete accounts.
type AccountRepository interface {
	Find(account Account) (Account, error)
	GetAll() ([]Account, error)
	Create(Account) (Account, error)
	Update(Account) (Account, error)
	Delete(Account) error
}

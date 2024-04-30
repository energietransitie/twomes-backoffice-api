package admin

import "time"

// An Admin has access to resources that are protected with an AdminToken.
type Admin struct {
	// ID is a short unique identifier.
	ID uint
	// Easily recognisable name for an admin.
	Name string
	// Time the admin was activated.
	// Tokens that are genereted before this time will be invalid.
	// Admins can be reactivated to invalidate old tokens.
	ActivatedAt time.Time
	// Time at which an admin expires.
	// This is by default the expiration date of any token.
	// This can be used to give temporary admin access.
	Expiry time.Time
	// Authorization token that is generated for the admin.
	// This will only be set upon creation.
	AuthorizationToken string
}

// Create a new admin.
func MakeAdmin(name string, expiry time.Time) Admin {
	if expiry.IsZero() {
		// Set to 1 year from now if empty.
		expiry = time.Now().UTC().AddDate(1, 0, 0)
	}

	return Admin{
		Name:        name,
		ActivatedAt: time.Now().UTC().Add(time.Second * -1),
		Expiry:      expiry,
	}
}

// Reactivate admin and invalidate all old tokens.
func (a *Admin) Reactivate() {
	a.ActivatedAt = time.Now().UTC().Add(time.Second * -1)
}

// Change the expiry date.
func (a *Admin) SetExpiry(expiry time.Time) {
	a.Expiry = expiry
}

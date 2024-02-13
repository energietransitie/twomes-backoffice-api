package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/authorization"
)

// An AuthorizationService exposes functionality available for authorization.
type AuthorizationService interface {
	CreateToken(kind authorization.AuthKind, id uint, expiry time.Time) (string, error)
	CreateTokenFromAuthorization(auth authorization.Authorization, expiry time.Time) (string, error)
	ParseToken(tokenString string) (authorization.AuthKind, uint, *authorization.Claims, error)
	ParseTokenToAuthorization(tokenString string) (*authorization.Authorization, error)
}

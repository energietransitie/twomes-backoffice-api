package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

// An AuthorizationService exposes functionality available for authorization.
type AuthorizationService interface {
	CreateToken(kind twomes.AuthKind, id uint, expiry time.Time) (string, error)
	CreateTokenFromAuthorization(auth twomes.Authorization, expiry time.Time) (string, error)
	ParseToken(tokenString string) (twomes.AuthKind, uint, *twomes.Claims, error)
	ParseTokenToAuthorization(tokenString string) (*twomes.Authorization, error)
}

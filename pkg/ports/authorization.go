package ports

import "github.com/energietransitie/twomes-api/pkg/twomes"

// An AuthorizationService exposes functionality available for authorization.
type AuthorizationService interface {
	CreateToken(kind twomes.AuthKind, id uint) (string, error)
	CreateTokenFromAuthorization(auth twomes.Authorization) (string, error)
	ParseToken(tokenString string) (twomes.AuthKind, uint, error)
	ParseTokenToAuthorization(tokenString string) (*twomes.Authorization, error)
}

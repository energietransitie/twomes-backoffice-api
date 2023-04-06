package twomes

import (
	"crypto"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
)

// Kind of authorization.
type AuthKind string

const (
	AdminToken             AuthKind = "adminToken"
	AccountToken           AuthKind = "accountToken"
	DeviceToken            AuthKind = "deviceToken"
	AccountActivationToken AuthKind = "accountActivationToken"
	InvalidToken           AuthKind = "invalidToken"
)

// An Authorization is used to check for permissions.
type Authorization struct {
	Kind   AuthKind
	ID     uint
	Claims *Claims
}

// Claims contained in the JWT.
type Claims struct {
	jwt.RegisteredClaims
	Kind AuthKind `json:"kind"`
}

func ParseTokenToAuthorization(tokenString string, pubkey crypto.PublicKey) (*Authorization, error) {
	kind, id, claims, err := ParseToken(tokenString, pubkey)
	if err != nil {
		return nil, err
	}

	return &Authorization{
		Kind:   kind,
		ID:     id,
		Claims: claims,
	}, nil
}

// Returns if the Authorization is of the specified kind.
func (a *Authorization) IsKind(kind AuthKind) bool {
	return a.Kind == kind
}

// Create a new token of a specified kind, for specified ID.
func NewToken(kind AuthKind, id uint, key crypto.PrivateKey) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "TwomesAPIv2",
			Subject:   strconv.FormatUint(uint64(id), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 24 * 365)),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Kind: kind,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(key)
}

// Create a new token from an Authorization.
func NewTokenFromAuthorization(auth Authorization, key crypto.PrivateKey) (string, error) {
	return NewToken(auth.Kind, auth.ID, key)
}

// Parse a signed token. Check if it is valid and return the kind of token and the corresponding ID.
func ParseToken(tokenString string, pubkey crypto.PublicKey) (AuthKind, uint, *Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return pubkey, nil
	})
	if err != nil {
		return InvalidToken, 0, &Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return InvalidToken, 0, &Claims{}, err
	}

	id, err := (strconv.ParseUint(claims.RegisteredClaims.Subject, 10, 64))
	if err != nil {
		return InvalidToken, 0, &Claims{}, err
	}

	return claims.Kind, uint(id), claims, nil
}

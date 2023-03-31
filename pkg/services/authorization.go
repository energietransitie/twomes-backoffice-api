package services

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"

	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidKeyAlgorithm = errors.New("invalid key algorithm")
)

type AuthorizationService struct {
	key crypto.PrivateKey
}

// Create a new AuthorizationService.
func NewAuthorizationService(key crypto.PrivateKey) *AuthorizationService {
	return &AuthorizationService{
		key: key,
	}
}

// Create a new AuthorizationService with the key from a file.
func NewAuthorizationServiceFromFile(path string) (*AuthorizationService, error) {
	key, err := KeyFromFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		// File did not exist, so generate it.
		err = generateKeyFile(path)
		if err != nil {
			return nil, err
		}

		// Load the key from newly generated file.
		key, err = KeyFromFile(path)
		if err != nil {
			return nil, err
		}
	}

	return &AuthorizationService{key: key}, nil
}

// Open a file and attempt to read the private key from it.
func KeyFromFile(path string) (crypto.PrivateKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pemString, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	defer logrus.Info("key was successfully loaded from file")

	block, _ := pem.Decode(pemString)
	return x509.ParseECPrivateKey(block.Bytes)
}

// Generate a new private key and save it to a file.
func generateKeyFile(path string) error {
	logrus.Info("generating key file")

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	ecder, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = pem.Encode(file, &pem.Block{Type: "EC PRIVATE KEY", Bytes: ecder})
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthorizationService) CreateToken(kind twomes.AuthKind, id uint) (string, error) {
	return twomes.NewToken(kind, id, s.key)
}

func (s *AuthorizationService) CreateTokenFromAuthorization(auth twomes.Authorization) (string, error) {
	return twomes.NewTokenFromAuthorization(auth, s.key)
}

func (s *AuthorizationService) ParseToken(tokenString string) (twomes.AuthKind, uint, error) {
	key, ok := s.key.(*ecdsa.PrivateKey)
	if !ok {
		return twomes.InvalidToken, 0, ErrInvalidKeyAlgorithm
	}

	return twomes.ParseToken(tokenString, key.Public())
}

func (s *AuthorizationService) ParseTokenToAuthorization(tokenString string) (*twomes.Authorization, error) {
	key, ok := s.key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyAlgorithm
	}

	return twomes.ParseTokenToAuthorization(tokenString, key.Public())
}
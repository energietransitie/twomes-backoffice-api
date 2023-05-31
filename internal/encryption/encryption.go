// Package encryption implements types that store in encrypted for in the DB.
// The encrpytion and decryption happens transparantly.
package encryption

import (
	"database/sql/driver"
	"errors"
)

var (
	ErrInvalidTypeInDB = errors.New("invalid type stored in database")
)

// EncryptedString will transparantly encrypt or decrypt data when
// it is saved or loaded from the database.
//
// TODO: Implement actual encryption.
// WARNING actual encryption is not yet implemented.
type EncryptedString string

func (e EncryptedString) Value() (driver.Value, error) {
	return []byte(e), nil
}

func (e *EncryptedString) Scan(src any) error {
	source, ok := src.([]byte)
	if !ok {
		return ErrInvalidTypeInDB
	}

	*e = EncryptedString(source)
	return nil
}

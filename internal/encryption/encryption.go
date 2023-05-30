// Package encryption implements types that store in encrypted for in the DB.
// The encrpytion and decryption happens transparantly.
package encryption

import (
	"database/sql/driver"
	"errors"
)

var (
	ErrInvalidTypeinDB = errors.New("invalid type stored in database")
)

// EncryptedString will transparantly encrypt or decrypt data when
// it is saved or loaded from the database.
//
// TODO: Implement actual encryption.
// WARNING actual encryption is not yet implemented.
type EncrpytedString string

func (e EncrpytedString) Value() (driver.Value, error) {
	return []byte(e), nil
}

func (e *EncrpytedString) Scan(src any) error {
	source, ok := src.([]byte)
	if !ok {
		return ErrInvalidTypeinDB
	}

	*e = EncrpytedString(source)
	return nil
}

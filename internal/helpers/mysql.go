// Package helpers implements common helper functions.
package helpers

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Check if an error is a mySQL duplicate entry error.
func IsMySQLDuplicateError(err error) bool {
	if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1062 {
		return true
	}
	return false
}

func IsMySQLRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

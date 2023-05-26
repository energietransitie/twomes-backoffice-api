// Package helpers implements common helper functions.
package helpers

// string array contains the given string.
func Contains(array []string, s string) bool {
	for _, str := range array {
		if str == s {
			return true
		}
	}

	return false
}

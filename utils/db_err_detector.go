package utils

import "strings"

type DbErrorDetector struct {
	Error error
}

func DbError(err error) *DbErrorDetector {
	return &DbErrorDetector{
		Error: err,
	}
}

// Check is Database error type is not found. 
// Always return false if error is nil
func (d DbErrorDetector) IsNotFound() bool {
	if d.Error == nil {
		return false
	}

	return strings.Contains(d.Error.Error(), "no rows in result set")
}

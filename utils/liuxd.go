package utils

import (
	"errors"
	"strings"
)

func GetRecoverError(err error, recErr any) (resErr error) {
	if err != nil {
		return err
	}
	if recErr != nil {
		switch recErr.(type) {
		case string:
			{
				msg, _ := recErr.(string)
				resErr = errors.New(msg)
			}
		case error:
			{
				if e, ok := recErr.(error); ok {
					resErr = e
				}
			}
		}
	}
	return resErr
}

// IsTruthy returns true if a string is a truthy value.
// Truthy values are "y", "yes", "true", "t", "on", "1" (case-insensitive); everything else is false.
// liuxd
func IsTruthy(val string) bool {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "y", "yes", "true", "t", "on", "1":
		return true
	default:
		return false
	}
}

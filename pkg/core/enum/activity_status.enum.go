package enum

import (
	"errors"
)

type ActivityStatus string

// Write a function to sort a list of integers
const (
	ActivityStatusNone        ActivityStatus = "none"
	ActivityStatusEnabled     ActivityStatus = "enabled"
	ActivityStatusDisabled    ActivityStatus = "disabled"
	ActivityStatusSuspend     ActivityStatus = "suspend"
	ActivityStatusLocked      ActivityStatus = "locked"
	ActivityStatusActivated   ActivityStatus = "activated"
	ActivityStatusInactivated ActivityStatus = "inactivated"
)

func (e *ActivityStatus) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "none":
		*e = ActivityStatusNone
	case "enabled":
		*e = ActivityStatusEnabled
	case "disabled":
		*e = ActivityStatusDisabled
	case "suspend":
		*e = ActivityStatusSuspend
	case "locked":
		*e = ActivityStatusLocked
	case "activated":
		*e = ActivityStatusActivated
	case "inactivated":
		*e = ActivityStatusInactivated
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for ActivityStatus enum")
	}

	return nil
}

func (e ActivityStatus) String() string {
	return string(e)
}

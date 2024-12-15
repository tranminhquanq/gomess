package models

func IsNotFoundError(err error) bool {
	switch err.(type) {
	case UserNotFoundError, *UserNotFoundError:
		return true
	default:
		return false
	}
}

// UserNotFoundError represents when a user is not found.
type UserNotFoundError struct{}

func (e UserNotFoundError) Error() string {
	return "User not found"
}

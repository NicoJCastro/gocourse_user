package user

import (
	"errors"
	"fmt"
)

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")
var ErrEmailRequired = errors.New("email is required")
var ErrPhoneRequired = errors.New("phone is required")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotUpdated = errors.New("user not updated")
var ErrUserNotDeleted = errors.New("user not deleted")
var ErrUserNotCreated = errors.New("user not created")
var ErrUserNotRetrieved = errors.New("user not retrieved")
var ErrUserNotCounted = errors.New("user not counted")
var ErrInvalidRequestType = errors.New("invalid request type")
var ErrInvalidDefaultLimitConfiguration = errors.New("invalid default limit configuration")
var ErrIDRequired = errors.New("id is required")
var ErrAtLeastOneFieldRequired = errors.New("at least one field is required")

// ErrNotFound es un error personalizado que incluye el ID del usuario no encontrado
type ErrNotFound struct {
	UserID string
}

// Error implementa la interfaz error
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("user with ID %s not found", e.UserID)
}

// Unwrap permite usar errors.Is() con este error
func (e *ErrNotFound) Unwrap() error {
	return ErrNotFoundBase
}

// NewErrNotFound crea una nueva instancia de ErrNotFound
func NewErrNotFound(userID string) *ErrNotFound {
	return &ErrNotFound{UserID: userID}
}

// ErrNotFoundBase es un error sentinela para comparaciones con errors.Is()
var ErrNotFoundBase = errors.New("user not found")

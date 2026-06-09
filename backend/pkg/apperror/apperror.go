package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// Kind classifies an application error.
type Kind uint8

const (
	KindValidation   Kind = iota + 1 // Input validation failed
	KindNotFound                     // Resource not found
	KindConflict                     // Duplicate or conflicting state
	KindBusiness                     // Business rule violation
	KindDatabase                     // Database operation failure
	KindInternal                     // Unexpected internal error
	KindUnauthorized                 // Authentication failure
)

// Error is the standard application error.
type Error struct {
	Kind    Kind
	Message string
	Code    string
	Field   string // For validation errors
	Err     error  // Wrapped error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// HTTPStatus maps an error Kind to an HTTP status code.
func (e *Error) HTTPStatus() int {
	switch e.Kind {
	case KindValidation:
		return http.StatusBadRequest
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	case KindBusiness:
		return http.StatusUnprocessableEntity
	case KindUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// --- Constructors ---

// Validation creates a validation error for a specific field.
func Validation(field, message string) *Error {
	return &Error{Kind: KindValidation, Field: field, Message: message}
}

// NotFound creates a not-found error.
func NotFound(resource, id string) *Error {
	return &Error{
		Kind:    KindNotFound,
		Message: fmt.Sprintf("%s with id %q not found", resource, id),
		Code:    "NOT_FOUND",
	}
}

// Conflict creates a conflict/duplicate error.
func Conflict(message string) *Error {
	return &Error{Kind: KindConflict, Message: message, Code: "CONFLICT"}
}

// Business creates a business rule violation error.
func Business(code, message string) *Error {
	return &Error{Kind: KindBusiness, Message: message, Code: code}
}

// Database wraps a database error with context.
func Database(op string, err error) *Error {
	return &Error{
		Kind:    KindDatabase,
		Message: fmt.Sprintf("database error during %s", op),
		Err:     err,
	}
}

// Internal wraps an unexpected error.
func Internal(message string, err error) *Error {
	return &Error{Kind: KindInternal, Message: message, Err: err}
}

// Is checks whether err is an *Error of the given Kind.
func Is(err error, kind Kind) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Kind == kind
	}
	return false
}

// AsError extracts an *Error from the error chain.
func AsError(err error) (*Error, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

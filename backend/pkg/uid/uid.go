package uid

import "github.com/google/uuid"

// New generates a new UUIDv7 (time-ordered).
// UUIDv7 provides better database index locality than v4.
func New() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		// Fallback to v4 if v7 generation fails (should not happen)
		return uuid.New()
	}
	return id
}

// Parse parses a UUID string.
func Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// Nil returns the zero-value UUID.
func Nil() uuid.UUID {
	return uuid.Nil
}

// IsNil checks whether a UUID is the zero value.
func IsNil(id uuid.UUID) bool {
	return id == uuid.Nil
}

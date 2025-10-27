package history

import (
	"strings"

	"github.com/google/uuid"
)

// RecordID is a value object representing a unique record identifier
type RecordID struct {
	value string
}

// NewRecordID generates a new unique record ID
func NewRecordID() (RecordID, error) {
	id := uuid.New().String()
	return RecordID{value: id}, nil
}

// ParseRecordID creates a RecordID from a string
func ParseRecordID(id string) (RecordID, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return RecordID{}, ErrInvalidRecordID
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return RecordID{}, ErrInvalidRecordID
	}

	return RecordID{value: id}, nil
}

// String returns the string representation
func (r RecordID) String() string {
	return r.value
}

// IsValid returns true if the ID is valid
func (r RecordID) IsValid() bool {
	return r.value != ""
}

// Equals checks if two IDs are equal
func (r RecordID) Equals(other RecordID) bool {
	return r.value == other.value
}

// ShortString returns a shortened ID for display (first 8 chars)
func (r RecordID) ShortString() string {
	if len(r.value) >= 8 {
		return r.value[:8]
	}
	return r.value
}

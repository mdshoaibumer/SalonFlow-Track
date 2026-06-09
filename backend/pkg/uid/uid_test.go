package uid

import (
	"testing"

	"github.com/google/uuid"
)

func TestNew_ReturnsValidUUID(t *testing.T) {
	id := New()
	if id == uuid.Nil {
		t.Fatal("New() returned nil UUID")
	}
	// UUIDv7 version byte
	if id.Version() != 7 {
		t.Errorf("expected UUID version 7, got %d", id.Version())
	}
}

func TestNew_IsUnique(t *testing.T) {
	seen := make(map[uuid.UUID]bool, 1000)
	for i := 0; i < 1000; i++ {
		id := New()
		if seen[id] {
			t.Fatalf("duplicate UUID generated at iteration %d", i)
		}
		seen[id] = true
	}
}

func TestNew_IsTimeSorted(t *testing.T) {
	a := New()
	b := New()
	// UUIDv7 embeds timestamp; later UUID should be lexicographically greater
	if a.String() > b.String() {
		t.Error("expected UUIDv7 to be time-sorted (a < b)")
	}
}

func TestParse_Valid(t *testing.T) {
	id := New()
	parsed, err := Parse(id.String())
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if parsed != id {
		t.Error("parsed UUID does not match original")
	}
}

func TestParse_Invalid(t *testing.T) {
	_, err := Parse("not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID string")
	}
}

func TestIsNil(t *testing.T) {
	if !IsNil(uuid.Nil) {
		t.Error("expected IsNil(uuid.Nil) to be true")
	}
	if IsNil(New()) {
		t.Error("expected IsNil(New()) to be false")
	}
}

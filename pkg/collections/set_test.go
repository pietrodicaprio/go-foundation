package collections

import (
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet[int]()

	s.Add(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("Len: got %d, want 3", s.Len())
	}

	if !s.Has(2) {
		t.Error("Has 2: should be true")
	}

	s.Remove(2)
	if s.Has(2) {
		t.Error("Has 2: should be false after removal")
	}

	items := s.Items()
	if len(items) != 2 {
		t.Errorf("Items: got %d, want 2", len(items))
	}

	s.Clear()
	if s.Len() != 0 {
		t.Error("Clear should empty the set")
	}
}

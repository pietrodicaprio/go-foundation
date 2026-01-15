package resiliency

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)

	// State: Closed
	err := cb.Execute(func() error { return nil })
	if err != nil || cb.State() != StateClosed {
		t.Error("Circuit should be closed and return nil")
	}

	// First failure
	_ = cb.Execute(func() error { return errors.New("fail1") })
	if cb.State() != StateClosed {
		t.Error("Circuit should still be closed after 1 failure")
	}

	// Second failure -> State: Open
	_ = cb.Execute(func() error { return errors.New("fail2") })
	if cb.State() != StateOpen {
		t.Error("Circuit should be open after reaching threshold")
	}

	// While Open
	err = cb.Execute(func() error { return nil })
	if err != ErrCircuitOpen {
		t.Errorf("Execute should return ErrCircuitOpen when open, got %v", err)
	}

	// Wait for timeout -> State: Half-Open (on next Execute)
	time.Sleep(60 * time.Millisecond)

	// First success in Half-Open -> State: Closed
	err = cb.Execute(func() error { return nil })
	if err != nil || cb.State() != StateClosed {
		t.Errorf("Circuit should be closed after success in half-open, got state %v", cb.State())
	}
}

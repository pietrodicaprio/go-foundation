package resiliency

import (
	"errors"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker: circuit is open")

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker protects a caller from repeated failures.
//
// It implements a state machine with Closed, Open, and Half-Open states.
//
// Example:
//
//	cb := resiliency.NewCircuitBreaker(3, time.Minute)
//	err := cb.Execute(func() error {
//		return doRiskyOperation()
//	})
type CircuitBreaker struct {
	mu            sync.RWMutex
	state         State
	failures      int
	threshold     int
	timeout       time.Duration
	lastFailure   time.Time
	onStateChange func(from, to State)
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
		state:     StateClosed,
	}
}

// OnStateChange registers a callback for state transitions.
func (cb *CircuitBreaker) OnStateChange(fn func(from, to State)) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.onStateChange = fn
}

// Execute wraps a function call with circuit breaker logic.
//
// Returns:
//
// ErrCircuitOpen if the circuit is open, otherwise the error from fn.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.allow() {
		return ErrCircuitOpen
	}

	err := fn()
	if err != nil {
		cb.onFailure()
		return err
	}

	cb.onSuccess()
	return nil
}

func (cb *CircuitBreaker) allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.changeState(StateHalfOpen)
			return true
		}
		return false
	case StateHalfOpen:
		// In half-open, we only allow one request to probe the system.
		// For simplicity in this foundation version, we'll allow it.
		return true
	}
	return false
}

func (cb *CircuitBreaker) onSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.changeState(StateClosed)
	}
	cb.failures = 0
}

func (cb *CircuitBreaker) onFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.state == StateClosed && cb.failures >= cb.threshold {
		cb.changeState(StateOpen)
	} else if cb.state == StateHalfOpen {
		cb.changeState(StateOpen)
	}
}

func (cb *CircuitBreaker) changeState(to State) {
	from := cb.state
	cb.state = to
	if cb.onStateChange != nil {
		cb.onStateChange(from, to)
	}
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

package skillplans

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitBreakerOpen is returned when the circuit breaker is open
var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	maxFailures  int
	resetTimeout time.Duration

	mu           sync.Mutex
	state        CircuitBreakerState
	failures     int
	lastFailTime time.Time
	successCount int
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Call executes the given function with circuit breaker protection
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if we should try to reset from Open state
	if cb.state == StateOpen {
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.successCount = 0
		} else {
			return ErrCircuitBreakerOpen
		}
	}

	// Execute the function
	err := fn()

	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}

	return err
}

// onFailure handles failure cases
func (cb *CircuitBreaker) onFailure() {
	cb.failures++
	cb.lastFailTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}
	case StateHalfOpen:
		// Single failure in half-open state reopens the circuit
		cb.state = StateOpen
	}
}

// onSuccess handles success cases
func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		cb.failures = 0
	case StateHalfOpen:
		cb.successCount++
		// Require 2 consecutive successes to fully close the circuit
		if cb.successCount >= 2 {
			cb.state = StateClosed
			cb.failures = 0
			cb.successCount = 0
		}
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Reset manually resets the circuit breaker
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
	cb.successCount = 0
}

package skillplans

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker_ClosedState(t *testing.T) {
	cb := NewCircuitBreaker(3, 1*time.Second)

	// Should allow successful calls
	err := cb.Call(func() error {
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestCircuitBreaker_OpensAfterMaxFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 1*time.Second)
	testErr := errors.New("test error")

	// First two failures should keep circuit closed
	for i := 0; i < 2; i++ {
		err := cb.Call(func() error {
			return testErr
		})
		assert.Error(t, err)
		assert.Equal(t, StateClosed, cb.GetState())
	}

	// Third failure should open the circuit
	err := cb.Call(func() error {
		return testErr
	})
	assert.Error(t, err)
	assert.Equal(t, StateOpen, cb.GetState())

	// Further calls should fail immediately
	err = cb.Call(func() error {
		return nil
	})
	assert.Equal(t, ErrCircuitBreakerOpen, err)
}

func TestCircuitBreaker_ResetsToHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(1, 100*time.Millisecond)
	testErr := errors.New("test error")

	// Open the circuit
	err := cb.Call(func() error {
		return testErr
	})
	assert.Error(t, err)
	assert.Equal(t, StateOpen, cb.GetState())

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)

	// Next call should be allowed (half-open state)
	callExecuted := false
	err = cb.Call(func() error {
		callExecuted = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, callExecuted)
}

func TestCircuitBreaker_HalfOpenToClosedAfterSuccess(t *testing.T) {
	cb := NewCircuitBreaker(1, 100*time.Millisecond)
	testErr := errors.New("test error")

	// Open the circuit
	err := cb.Call(func() error {
		return testErr
	})
	assert.Error(t, err)

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)

	// Two successful calls should close the circuit
	for i := 0; i < 2; i++ {
		err = cb.Call(func() error {
			return nil
		})
		assert.NoError(t, err)
	}

	assert.Equal(t, StateClosed, cb.GetState())
}

func TestCircuitBreaker_HalfOpenToOpenOnFailure(t *testing.T) {
	cb := NewCircuitBreaker(1, 100*time.Millisecond)
	testErr := errors.New("test error")

	// Open the circuit
	err := cb.Call(func() error {
		return testErr
	})
	assert.Error(t, err)

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)

	// Failure in half-open state should reopen circuit
	err = cb.Call(func() error {
		return testErr
	})
	assert.Error(t, err)
	assert.Equal(t, StateOpen, cb.GetState())
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(1, 1*time.Second)
	testErr := errors.New("test error")

	// Open the circuit
	err := cb.Call(func() error {
		return testErr
	})
	assert.Error(t, err)
	assert.Equal(t, StateOpen, cb.GetState())

	// Manual reset
	cb.Reset()
	assert.Equal(t, StateClosed, cb.GetState())

	// Should allow calls again
	err = cb.Call(func() error {
		return nil
	})
	assert.NoError(t, err)
}

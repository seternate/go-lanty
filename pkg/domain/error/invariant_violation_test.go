package domainerr

import (
	"errors"
	"testing"
)

func TestInvariantViolationError(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		err := InvariantViolationErr("game", "123")
		if err.Entity != "game" {
			t.Errorf("expected Entity 'game', got %q", err.Entity)
		}
		if err.ID != "123" {
			t.Errorf("expected ID '123', got %q", err.ID)
		}
		if err.Cause != nil {
			t.Error("expected nil Cause")
		}
	})

	t.Run("Error without cause", func(t *testing.T) {
		err := InvariantViolationErr("game", "123")
		expected := "failed to enforce invariant: game (123)"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error with cause", func(t *testing.T) {
		cause := errors.New("invalid state")
		err := InvariantViolationErr("game", "123").WithCause(cause)
		expected := "failed to enforce invariant: game (123): invalid state"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("ErrorCode", func(t *testing.T) {
		err := InvariantViolationErr("game", "123")
		if err.ErrorCode() != string(ErrorCodeInvariantViolation) {
			t.Errorf("expected %q, got %q", ErrorCodeInvariantViolation, err.ErrorCode())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("test cause")
		err := InvariantViolationErr("game", "123").WithCause(cause)
		if err.Unwrap() != cause {
			t.Error("Unwrap() should return the cause")
		}
	})

	t.Run("UserError without cause", func(t *testing.T) {
		err := InvariantViolationErr("game", "123")
		expected := "failed to enforce invariant: game (123)"
		if err.UserError() != expected {
			t.Errorf("expected %q, got %q", expected, err.UserError())
		}
	})

	t.Run("UserError with UserError cause", func(t *testing.T) {
		cause := ValidationErr("field", "invalid value")
		err := InvariantViolationErr("game", "123").WithCause(cause)
		expected := "failed to enforce invariant: game (123): field: invalid value"
		if err.UserError() != expected {
			t.Errorf("expected %q, got %q", expected, err.UserError())
		}
	})

	t.Run("UserError with non-UserError cause", func(t *testing.T) {
		cause := errors.New("some error")
		err := InvariantViolationErr("game", "123").WithCause(cause)
		expected := "failed to enforce invariant: game (123)"
		if err.UserError() != expected {
			t.Errorf("expected %q, got %q", expected, err.UserError())
		}
	})

	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("test cause")
		err := InvariantViolationErr("game", "123")
		result := err.WithCause(cause)
		if result != err {
			t.Error("WithCause should return the same error instance")
		}
		if err.Cause != cause {
			t.Error("WithCause should set the cause")
		}
	})

	t.Run("interface compliance", func(t *testing.T) {
		var _ InvariantViolation = (*InvariantViolationError)(nil)
		var _ UserError = (*InvariantViolationError)(nil)
		var _ DomainError = (*InvariantViolationError)(nil)
	})
}

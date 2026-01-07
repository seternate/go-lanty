package domainerr

import (
	"errors"
	"testing"
)

func TestTrustedInvariantViolationError(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		err := TrustedInvariantViolationErr("game", "123")
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
		err := TrustedInvariantViolationErr("game", "123")
		expected := "failed to enforce trusted invariant: game (123)"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error with cause", func(t *testing.T) {
		cause := errors.New("invalid state")
		err := TrustedInvariantViolationErr("game", "123").WithCause(cause)
		expected := "failed to enforce trusted invariant: game (123): invalid state"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("ErrorCode", func(t *testing.T) {
		err := TrustedInvariantViolationErr("game", "123")
		if err.ErrorCode() != string(ErrorCodeTrustedInvariantViolation) {
			t.Errorf("expected %q, got %q", ErrorCodeTrustedInvariantViolation, err.ErrorCode())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("test cause")
		err := TrustedInvariantViolationErr("game", "123").WithCause(cause)
		if err.Unwrap() != cause {
			t.Error("Unwrap() should return the cause")
		}
	})

	t.Run("Unwrap without cause", func(t *testing.T) {
		err := TrustedInvariantViolationErr("game", "123")
		if err.Unwrap() != nil {
			t.Error("Unwrap() should return nil when no cause")
		}
	})

	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("test cause")
		err := TrustedInvariantViolationErr("game", "123")
		result := err.WithCause(cause)
		if result != err {
			t.Error("WithCause should return the same error instance")
		}
		if err.Cause != cause {
			t.Error("WithCause should set the cause")
		}
	})

	t.Run("interface compliance", func(t *testing.T) {
		var _ TrustedInvariantViolation = (*TrustedInvariantViolationError)(nil)
		var _ Internal = (*TrustedInvariantViolationError)(nil)
		var _ DomainError = (*TrustedInvariantViolationError)(nil)
	})
}

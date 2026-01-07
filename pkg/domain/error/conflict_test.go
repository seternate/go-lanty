package domainerr

import (
	"errors"
	"testing"
)

func TestConflictError(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		err := ConflictErr("game", "123")
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
		err := ConflictErr("game", "123")
		expected := "game (123) already exists"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error with cause", func(t *testing.T) {
		cause := errors.New("database constraint violation")
		err := ConflictErr("game", "123").WithCause(cause)
		expected := "game (123) already exists: database constraint violation"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("ErrorCode", func(t *testing.T) {
		err := ConflictErr("game", "123")
		if err.ErrorCode() != string(ErrorCodeConflict) {
			t.Errorf("expected %q, got %q", ErrorCodeConflict, err.ErrorCode())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("test cause")
		err := ConflictErr("game", "123").WithCause(cause)
		if err.Unwrap() != cause {
			t.Error("Unwrap() should return the cause")
		}
	})

	t.Run("Unwrap without cause", func(t *testing.T) {
		err := ConflictErr("game", "123")
		if err.Unwrap() != nil {
			t.Error("Unwrap() should return nil when no cause")
		}
	})

	t.Run("UserError", func(t *testing.T) {
		err := ConflictErr("game", "123")
		expected := "game (123) already exists"
		if err.UserError() != expected {
			t.Errorf("expected %q, got %q", expected, err.UserError())
		}
	})

	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("test cause")
		err := ConflictErr("game", "123")
		result := err.WithCause(cause)
		if result != err {
			t.Error("WithCause should return the same error instance")
		}
		if err.Cause != cause {
			t.Error("WithCause should set the cause")
		}
	})

	t.Run("interface compliance", func(t *testing.T) {
		var _ Conflict = (*ConflictError)(nil)
		var _ UserError = (*ConflictError)(nil)
		var _ DomainError = (*ConflictError)(nil)
	})
}

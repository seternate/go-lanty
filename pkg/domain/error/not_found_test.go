package domainerr

import (
	"errors"
	"testing"
)

func TestNotFoundError(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		err := NotFoundErr("game", "123")
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
		err := NotFoundErr("game", "123")
		expected := "game (123) not found"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error with cause", func(t *testing.T) {
		cause := errors.New("database query failed")
		err := NotFoundErr("game", "123").WithCause(cause)
		expected := "game (123) not found: database query failed"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("ErrorCode", func(t *testing.T) {
		err := NotFoundErr("game", "123")
		if err.ErrorCode() != string(ErrorCodeNotFound) {
			t.Errorf("expected %q, got %q", ErrorCodeNotFound, err.ErrorCode())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("test cause")
		err := NotFoundErr("game", "123").WithCause(cause)
		if err.Unwrap() != cause {
			t.Error("Unwrap() should return the cause")
		}
	})

	t.Run("UserError", func(t *testing.T) {
		err := NotFoundErr("game", "123")
		expected := "game (123) not found"
		if err.UserError() != expected {
			t.Errorf("expected %q, got %q", expected, err.UserError())
		}
	})

	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("test cause")
		err := NotFoundErr("game", "123")
		result := err.WithCause(cause)
		if result != err {
			t.Error("WithCause should return the same error instance")
		}
		if err.Cause != cause {
			t.Error("WithCause should set the cause")
		}
	})

	t.Run("interface compliance", func(t *testing.T) {
		var _ NotFound = (*NotFoundError)(nil)
		var _ UserError = (*NotFoundError)(nil)
		var _ DomainError = (*NotFoundError)(nil)
	})
}

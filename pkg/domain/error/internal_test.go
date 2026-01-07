package domainerr

import (
	"errors"
	"testing"
)

func TestInternalError(t *testing.T) {
	t.Run("constructor without args", func(t *testing.T) {
		err := InternalErr("something went wrong")
		if err.Reason != "something went wrong" {
			t.Errorf("expected Reason 'something went wrong', got %q", err.Reason)
		}
		if err.Cause != nil {
			t.Error("expected nil Cause")
		}
	})

	t.Run("constructor with args", func(t *testing.T) {
		err := InternalErr("failed to %s: %d", "connect", 500)
		expected := "failed to connect: 500"
		if err.Reason != expected {
			t.Errorf("expected Reason %q, got %q", expected, err.Reason)
		}
	})

	t.Run("Error without cause", func(t *testing.T) {
		err := InternalErr("something went wrong")
		expected := "internal domain error: something went wrong"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := InternalErr("something went wrong").WithCause(cause)
		expected := "internal domain error: something went wrong: underlying error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("ErrorCode", func(t *testing.T) {
		err := InternalErr("test")
		if err.ErrorCode() != string(ErrorCodeInternal) {
			t.Errorf("expected %q, got %q", ErrorCodeInternal, err.ErrorCode())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("test cause")
		err := InternalErr("test").WithCause(cause)
		if err.Unwrap() != cause {
			t.Error("Unwrap() should return the cause")
		}
	})

	t.Run("Unwrap without cause", func(t *testing.T) {
		err := InternalErr("test")
		if err.Unwrap() != nil {
			t.Error("Unwrap() should return nil when no cause")
		}
	})

	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("test cause")
		err := InternalErr("test")
		result := err.WithCause(cause)
		if result != err {
			t.Error("WithCause should return the same error instance")
		}
		if err.Cause != cause {
			t.Error("WithCause should set the cause")
		}
	})

	t.Run("interface compliance", func(t *testing.T) {
		var _ Internal = (*InternalError)(nil)
		var _ DomainError = (*InternalError)(nil)
	})
}

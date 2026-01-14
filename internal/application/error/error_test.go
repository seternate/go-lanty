package apperr

import (
	"errors"
	"testing"

	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

func TestDescribeError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		result := DescribeError(nil)
		if result != nil {
			t.Errorf("expected nil, got %+v", result)
		}
	})

	t.Run("Internal error", func(t *testing.T) {
		err := domainerr.InternalErr("something went wrong")
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != "internal" {
			t.Errorf("expected Code 'internal', got %q", result.Code)
		}
		if result.Message != "Something went wrong. Please try again later." {
			t.Errorf("expected Message 'Something went wrong. Please try again later.', got %q", result.Message)
		}
	})

	t.Run("Internal error with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := domainerr.InternalErr("something went wrong").WithCause(cause)
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != "internal" {
			t.Errorf("expected Code 'internal', got %q", result.Code)
		}
		if result.Message != "Something went wrong. Please try again later." {
			t.Errorf("expected Message 'Something went wrong. Please try again later.', got %q", result.Message)
		}
	})

	t.Run("TrustedInvariantViolationError", func(t *testing.T) {
		err := domainerr.TrustedInvariantViolationErr("game", "123")
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != "internal" {
			t.Errorf("expected Code 'internal' (TrustedInvariantViolation implements Internal), got %q", result.Code)
		}
		if result.Message != "Something went wrong. Please try again later." {
			t.Errorf("expected Message 'Something went wrong. Please try again later.', got %q", result.Message)
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		err := domainerr.ValidationErr("field", "invalid value")
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != string(domainerr.ErrorCodeValidation) {
			t.Errorf("expected Code %q, got %q", domainerr.ErrorCodeValidation, result.Code)
		}
		expectedMessage := "field: invalid value"
		if result.Message != expectedMessage {
			t.Errorf("expected Message %q, got %q", expectedMessage, result.Message)
		}
	})

	t.Run("ValidationError with expected and got", func(t *testing.T) {
		err := domainerr.ValidationErr("field", "invalid value").
			WithExpected("string").
			WithGot("int")
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != string(domainerr.ErrorCodeValidation) {
			t.Errorf("expected Code %q, got %q", domainerr.ErrorCodeValidation, result.Code)
		}
		expectedMessage := "field: invalid value: expected string, got int"
		if result.Message != expectedMessage {
			t.Errorf("expected Message %q, got %q", expectedMessage, result.Message)
		}
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		validationErrs := domainerr.ValidationErrs()
		validationErrs.Errors = append(validationErrs.Errors,
			domainerr.ValidationErr("field1", "error1"),
			domainerr.ValidationErr("field2", "error2"),
		)
		result := DescribeError(validationErrs)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != string(domainerr.ErrorCodeValidation) {
			t.Errorf("expected Code %q, got %q", domainerr.ErrorCodeValidation, result.Code)
		}
		expectedMessage := "field1: error1: field2: error2"
		if result.Message != expectedMessage {
			t.Errorf("expected Message %q, got %q", expectedMessage, result.Message)
		}
	})

	t.Run("NotFoundError", func(t *testing.T) {
		err := domainerr.NotFoundErr("game", "123")
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != string(domainerr.ErrorCodeNotFound) {
			t.Errorf("expected Code %q, got %q", domainerr.ErrorCodeNotFound, result.Code)
		}
		expectedMessage := "game (123) not found"
		if result.Message != expectedMessage {
			t.Errorf("expected Message %q, got %q", expectedMessage, result.Message)
		}
	})

	t.Run("ConflictError", func(t *testing.T) {
		err := domainerr.ConflictErr("game", "123")
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != string(domainerr.ErrorCodeConflict) {
			t.Errorf("expected Code %q, got %q", domainerr.ErrorCodeConflict, result.Code)
		}
		expectedMessage := "game (123) already exists"
		if result.Message != expectedMessage {
			t.Errorf("expected Message %q, got %q", expectedMessage, result.Message)
		}
	})

	t.Run("InvariantViolationError", func(t *testing.T) {
		err := domainerr.InvariantViolationErr("game", "123")
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != string(domainerr.ErrorCodeInvariantViolation) {
			t.Errorf("expected Code %q, got %q", domainerr.ErrorCodeInvariantViolation, result.Code)
		}
		expectedMessage := "failed to enforce invariant: game (123)"
		if result.Message != expectedMessage {
			t.Errorf("expected Message %q, got %q", expectedMessage, result.Message)
		}
	})

	t.Run("InvariantViolationError with UserError cause", func(t *testing.T) {
		cause := domainerr.ValidationErr("field", "invalid")
		err := domainerr.InvariantViolationErr("game", "123").WithCause(cause)
		result := DescribeError(err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Code != string(domainerr.ErrorCodeInvariantViolation) {
			t.Errorf("expected Code %q, got %q", domainerr.ErrorCodeInvariantViolation, result.Code)
		}
		expectedMessage := "failed to enforce invariant: game (123): field: invalid"
		if result.Message != expectedMessage {
			t.Errorf("expected Message %q, got %q", expectedMessage, result.Message)
		}
	})

	t.Run("regular error", func(t *testing.T) {
		err := errors.New("regular error")
		result := DescribeError(err)
		if result != nil {
			t.Errorf("expected nil for regular error, got %+v", result)
		}
	})

	t.Run("wrapped Internal error", func(t *testing.T) {
		internalErr := domainerr.InternalErr("something went wrong")
		// Test with a properly wrapped error using errors.Join
		wrappedErr := errors.Join(errors.New("wrapper"), internalErr)
		result := DescribeError(wrappedErr)
		if result == nil {
			t.Fatal("expected non-nil result for wrapped Internal error")
		}
		if result.Code != "internal" {
			t.Errorf("expected Code 'internal', got %q", result.Code)
		}
	})

	t.Run("wrapped UserError", func(t *testing.T) {
		notFoundErr := domainerr.NotFoundErr("game", "123")
		wrappedErr := errors.Join(errors.New("wrapper"), notFoundErr)
		result := DescribeError(wrappedErr)
		if result == nil {
			t.Fatal("expected non-nil result for wrapped UserError")
		}
		if result.Code != string(domainerr.ErrorCodeNotFound) {
			t.Errorf("expected Code %q, got %q", domainerr.ErrorCodeNotFound, result.Code)
		}
		expectedMessage := "game (123) not found"
		if result.Message != expectedMessage {
			t.Errorf("expected Message %q, got %q", expectedMessage, result.Message)
		}
	})
}

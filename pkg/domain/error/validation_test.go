package domainerr

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	t.Run("constructor without args", func(t *testing.T) {
		err := ValidationErr("name", "cannot be empty")
		if err.Field != "name" {
			t.Errorf("expected Field 'name', got %q", err.Field)
		}
		if err.Message != "cannot be empty" {
			t.Errorf("expected Message 'cannot be empty', got %q", err.Message)
		}
		if err.Cause != nil {
			t.Error("expected nil Cause")
		}
	})

	t.Run("constructor with args", func(t *testing.T) {
		err := ValidationErr("age", "must be at least %d", 18)
		if err.Field != "age" {
			t.Errorf("expected Field 'age', got %q", err.Field)
		}
		expected := "must be at least 18"
		if err.Message != expected {
			t.Errorf("expected Message %q, got %q", expected, err.Message)
		}
	})

	t.Run("Error with only field and message", func(t *testing.T) {
		err := ValidationErr("name", "cannot be empty")
		expected := "name: cannot be empty"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error with expected and got", func(t *testing.T) {
		err := ValidationErr("age", "invalid value").
			WithExpected(">= 18").
			WithGot("15")
		expected := "age: invalid value: expected >= 18, got 15"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error with only got", func(t *testing.T) {
		err := ValidationErr("age", "invalid value").
			WithGot("15")
		expected := "age: invalid value: got 15"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Error with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := ValidationErr("name", "cannot be empty").WithCause(cause)
		expected := "name: cannot be empty: underlying error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("ErrorCode", func(t *testing.T) {
		err := ValidationErr("name", "test")
		if err.ErrorCode() != string(ErrorCodeValidation) {
			t.Errorf("expected %q, got %q", ErrorCodeValidation, err.ErrorCode())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("test cause")
		err := ValidationErr("name", "test").WithCause(cause)
		if err.Unwrap() != cause {
			t.Error("Unwrap() should return the cause")
		}
	})

	t.Run("UserError", func(t *testing.T) {
		err := ValidationErr("name", "cannot be empty")
		expected := "name: cannot be empty"
		if err.UserError() != expected {
			t.Errorf("expected %q, got %q", expected, err.UserError())
		}
	})

	t.Run("WithExpected with args", func(t *testing.T) {
		err := ValidationErr("age", "test")
		result := err.WithExpected(">= %d", 18)
		if result != err {
			t.Error("WithExpected should return the same error instance")
		}
		if err.Expected != ">= 18" {
			t.Errorf("expected Expected '>= 18', got %q", err.Expected)
		}
	})

	t.Run("WithGot with args", func(t *testing.T) {
		err := ValidationErr("age", "test")
		result := err.WithGot("value: %s", "invalid")
		if result != err {
			t.Error("WithGot should return the same error instance")
		}
		if err.Got != "value: invalid" {
			t.Errorf("expected Got 'value: invalid', got %q", err.Got)
		}
	})

	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("test cause")
		err := ValidationErr("name", "test")
		result := err.WithCause(cause)
		if result != err {
			t.Error("WithCause should return the same error instance")
		}
		if err.Cause != cause {
			t.Error("WithCause should set the cause")
		}
	})

	t.Run("interface compliance", func(t *testing.T) {
		var _ Validation = (*ValidationError)(nil)
		var _ UserError = (*ValidationError)(nil)
		var _ DomainError = (*ValidationError)(nil)
	})
}

func TestValidationErrors(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		errs := ValidationErrs()
		if errs.Errors == nil {
			t.Error("expected Errors to be initialized")
		}
		if len(errs.Errors) != 0 {
			t.Errorf("expected empty Errors slice, got %d", len(errs.Errors))
		}
	})

	t.Run("Wrap with ValidationError", func(t *testing.T) {
		errs := ValidationErrs()
		ve := ValidationErr("name", "cannot be empty")
		result := errs.Wrap(ve)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
		if len(errs.Errors) != 1 {
			t.Errorf("expected 1 error, got %d", len(errs.Errors))
		}
		if errs.Errors[0] != ve {
			t.Error("expected wrapped error to be the same instance")
		}
	})

	t.Run("Wrap with ValidationErrors", func(t *testing.T) {
		errs := ValidationErrs()
		ve1 := ValidationErr("name", "cannot be empty")
		ve2 := ValidationErr("age", "must be positive")
		otherErrs := ValidationErrs()
		otherErrs.Wrap(ve1)
		otherErrs.Wrap(ve2)

		if len(otherErrs.Errors) != 2 {
			t.Fatalf("expected otherErrs to have 2 errors, got %d", len(otherErrs.Errors))
		}

		result := errs.Wrap(otherErrs)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
		if len(errs.Errors) != 2 {
			t.Errorf("expected 2 errors, got %d", len(errs.Errors))
		}
		if len(errs.Errors) >= 2 {
			if errs.Errors[0] != ve1 || errs.Errors[1] != ve2 {
				t.Error("expected errors to match ve1 and ve2")
			}
		}
	})

	t.Run("Wrap with nil", func(t *testing.T) {
		errs := ValidationErrs()
		result := errs.Wrap(nil)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
		if len(errs.Errors) != 0 {
			t.Errorf("expected 0 errors, got %d", len(errs.Errors))
		}
	})

	t.Run("Wrap with non-validation error", func(t *testing.T) {
		errs := ValidationErrs()
		otherErr := errors.New("not a validation error")
		result := errs.Wrap(otherErr)
		if result == nil {
			t.Error("expected error when wrapping non-validation error")
		}
		if len(errs.Errors) != 0 {
			t.Errorf("expected 0 errors, got %d", len(errs.Errors))
		}
	})

	t.Run("Error with no errors", func(t *testing.T) {
		errs := ValidationErrs()
		if errs.Error() != "" {
			t.Errorf("expected empty string, got %q", errs.Error())
		}
	})

	t.Run("Error with single error", func(t *testing.T) {
		errs := ValidationErrs()
		ve := ValidationErr("name", "cannot be empty")
		errs.Wrap(ve)
		expected := "name: cannot be empty"
		if errs.Error() != expected {
			t.Errorf("expected %q, got %q", expected, errs.Error())
		}
	})

	t.Run("Error with multiple errors", func(t *testing.T) {
		errs := ValidationErrs()
		ve1 := ValidationErr("name", "cannot be empty")
		ve2 := ValidationErr("age", "must be positive")
		errs.Wrap(ve1)
		errs.Wrap(ve2)
		result := errs.Error()
		// Order may vary, so check that both messages are present
		if result != "name: cannot be empty: age: must be positive" {
			t.Errorf("unexpected error message: %q", result)
		}
	})

	t.Run("Unwrap with no errors", func(t *testing.T) {
		errs := ValidationErrs()
		result := errs.Unwrap()
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("Unwrap with errors", func(t *testing.T) {
		errs := ValidationErrs()
		ve1 := ValidationErr("name", "cannot be empty")
		ve2 := ValidationErr("age", "must be positive")
		errs.Wrap(ve1)
		errs.Wrap(ve2)
		result := errs.Unwrap()
		if len(result) != 2 {
			t.Errorf("expected 2 errors, got %d", len(result))
		}
	})

	t.Run("UserError with no errors", func(t *testing.T) {
		errs := ValidationErrs()
		if errs.UserError() != "" {
			t.Errorf("expected empty string, got %q", errs.UserError())
		}
	})

	t.Run("UserError with errors", func(t *testing.T) {
		errs := ValidationErrs()
		ve1 := ValidationErr("name", "cannot be empty")
		ve2 := ValidationErr("age", "must be positive")
		errs.Wrap(ve1)
		errs.Wrap(ve2)
		result := errs.UserError()
		// Order may vary, so check that both messages are present
		if result != "name: cannot be empty: age: must be positive" {
			t.Errorf("unexpected user error message: %q", result)
		}
	})

	t.Run("HasErrors with no errors", func(t *testing.T) {
		errs := ValidationErrs()
		if errs.HasErrors() {
			t.Error("expected HasErrors() to return false")
		}
	})

	t.Run("HasErrors with errors", func(t *testing.T) {
		errs := ValidationErrs()
		ve := ValidationErr("name", "cannot be empty")
		errs.Wrap(ve)
		if !errs.HasErrors() {
			t.Error("expected HasErrors() to return true")
		}
	})

	t.Run("OrNil with no errors", func(t *testing.T) {
		errs := ValidationErrs()
		result := errs.OrNil()
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("OrNil with errors", func(t *testing.T) {
		errs := ValidationErrs()
		ve := ValidationErr("name", "cannot be empty")
		errs.Wrap(ve)
		result := errs.OrNil()
		if result != errs {
			t.Error("expected OrNil() to return the same instance when there are errors")
		}
	})

	t.Run("ErrorCode", func(t *testing.T) {
		errs := ValidationErrs()
		if errs.ErrorCode() != string(ErrorCodeValidation) {
			t.Errorf("expected %q, got %q", ErrorCodeValidation, errs.ErrorCode())
		}
	})

	t.Run("interface compliance", func(t *testing.T) {
		var _ Validation = (*ValidationErrors)(nil)
		var _ UserError = (*ValidationErrors)(nil)
		var _ DomainError = (*ValidationErrors)(nil)
	})
}

package domainerr

import (
	"errors"
	"fmt"
	"strings"
)

var _ Validation = (*ValidationError)(nil)

type ValidationError struct {
	Field    string
	Message  string
	Expected string
	Got      string
	Cause    error
}

func ValidationErr(field string, message string, args ...any) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf(message, args...),
	}
}

func (e *ValidationError) WithCause(cause error) *ValidationError {
	e.Cause = cause
	return e
}

func (e *ValidationError) WithExpected(expected string, args ...any) *ValidationError {
	e.Expected = fmt.Sprintf(expected, args...)
	return e
}

func (e *ValidationError) WithGot(got string, args ...any) *ValidationError {
	e.Got = fmt.Sprintf(got, args...)
	return e
}

func (e ValidationError) Error() string {
	message := e.msg()

	if e.Cause != nil && len(e.Cause.Error()) > 0 {
		message += ": " + e.Cause.Error()
	}

	return message
}

func (e *ValidationError) Unwrap() error {
	return e.Cause
}

func (e ValidationError) UserError() string {
	return e.msg()
}

func (e ValidationError) msg() string {
	var message string

	if e.Expected != "" || e.Got != "" {
		if e.Expected != "" && e.Got != "" {
			message = e.Field + ": " + e.Message + ": expected " + e.Expected + ", got " + e.Got
		} else if e.Got != "" {
			message = e.Field + ": " + e.Message + ": got " + e.Got
		}
	} else {
		message = e.Field + ": " + e.Message
	}

	return message
}

func (e ValidationError) ErrorCode() string {
	return string(ErrorCodeValidation)
}

func (ValidationError) DomainError() {}
func (ValidationError) Validation()  {}

var _ Validation = (*ValidationErrors)(nil)

type ValidationErrors struct {
	Errors []*ValidationError
}

func ValidationErrs() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]*ValidationError, 0),
	}
}

func (e *ValidationErrors) Wrap(err error) error {
	if err == nil {
		return nil
	}

	var ve *ValidationError
	if errors.As(err, &ve) {
		e.Errors = append(e.Errors, ve)
		return nil
	}

	var ves *ValidationErrors
	if errors.As(err, &ves) {
		e.Errors = append(e.Errors, ves.Errors...)
		return nil
	}

	return fmt.Errorf("expected validation error for wrapping: %w", err)
}

func (e *ValidationErrors) Error() string {
	var errMessages []string
	for _, err := range e.Errors {
		errMsg := err.Error()
		if errMsg != "" {
			errMessages = append(errMessages, errMsg)
		}
	}

	return strings.Join(errMessages, ": ")
}

func (e *ValidationErrors) Unwrap() []error {
	if len(e.Errors) == 0 {
		return nil
	}

	var errs []error
	for _, err := range e.Errors {
		errs = append(errs, err)
	}

	return errs
}

func (e *ValidationErrors) UserError() string {
	var errMessages []string
	for _, err := range e.Errors {
		errMsg := err.UserError()
		if errMsg != "" {
			errMessages = append(errMessages, errMsg)
		}
	}

	return strings.Join(errMessages, ": ")
}

func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *ValidationErrors) OrNil() error {
	if e.HasErrors() {
		return e
	}
	return nil
}

func (e *ValidationErrors) ErrorCode() string {
	return string(ErrorCodeValidation)
}

func (e *ValidationErrors) DomainError() {}
func (e *ValidationErrors) Validation()  {}

package error

import (
	"errors"
	"fmt"
	"strings"
)

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

func (ValidationError) DomainError() {}
func (ValidationError) Validation()  {}

type ValidationErrors struct {
	Errors  []*ValidationError
	Message []string
}

func ValidationErrs() *ValidationErrors {
	return &ValidationErrors{
		Errors:  make([]*ValidationError, 0),
		Message: make([]string, 0),
	}
}

func (e *ValidationErrors) WithMessage(message string, args ...any) *ValidationErrors {
	msg := fmt.Sprintf(message, args...)
	if len(msg) > 0 {
		e.Message = append(e.Message, msg)
	}
	return e
}

func (e *ValidationErrors) Wrap(err error, msg ...string) error {
	if err == nil {
		return nil
	}

	var ve *ValidationError
	if errors.As(err, &ve) {
		e.Errors = append(e.Errors, ve)
		for _, m := range msg {
			if len(m) > 0 {
				e.Message = append(e.Message, m)
			}
		}
		return nil
	}

	var ves *ValidationErrors
	if errors.As(err, &ves) {
		e.Errors = append(e.Errors, ves.Errors...)

		vesMessages := ves.Message
		if len(vesMessages) > 0 && len(e.Message) > 0 && vesMessages[0] == e.Message[len(e.Message)-1] {
			vesMessages = vesMessages[1:]
		}
		e.Message = append(e.Message, vesMessages...)

		for _, m := range msg {
			if len(m) > 0 {
				e.Message = append(e.Message, m)
			}
		}
		return nil
	}

	return err
}

func (e *ValidationErrors) Error() string {
	message := strings.Join(e.Message, ": ")

	var errMessages []string
	for _, err := range e.Errors {
		errMsg := err.Error()
		if errMsg != "" {
			errMessages = append(errMessages, errMsg)
		}
	}

	if message != "" && len(errMessages) > 0 {
		message += ": "
	}

	message += strings.Join(errMessages, ": ")

	return message
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
	message := strings.Join(e.Message, ": ")

	var errMessages []string
	for _, err := range e.Errors {
		errMsg := err.UserError()
		if errMsg != "" {
			errMessages = append(errMessages, errMsg)
		}
	}

	if message != "" && len(errMessages) > 0 {
		message += ": "
	}

	message += strings.Join(errMessages, ": ")

	return message
}

func (e *ValidationErrors) DomainError() {}
func (e *ValidationErrors) Validation()  {}

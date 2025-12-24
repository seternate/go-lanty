package error

import (
	"errors"
	"strings"

	domainErrors "github.com/seternate/go-lanty/pkg/domain/error"
)

type Code string

const (
	CodeValidation Code = "VALIDATION"
	CodeConflict   Code = "CONFLICT"
	CodeNotFound   Code = "NOT_FOUND"
	CodeInternal   Code = "INTERNAL"
)

type Error struct {
	Code    Code
	Message string
	Cause   error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func Err(err error, msg ...string) *Error {
	if err == nil {
		return nil
	}

	message := ""
	if len(msg) > 0 {
		nonEmpty := make([]string, 0, len(msg))
		for _, m := range msg {
			if m != "" {
				nonEmpty = append(nonEmpty, m)
			}
		}
		if len(nonEmpty) > 0 {
			message = strings.Join(nonEmpty, ": ")
		}
	}

	// Validation
	var v domainErrors.Validation
	if errors.As(err, &v) {
		if len(message) == 0 {
			message = v.Error()
		} else {
			message = message + ": " + v.Error()
		}
		return &Error{
			Code:    CodeValidation,
			Message: message,
			Cause:   err,
		}
	}

	// Conflict
	var c domainErrors.Conflict
	if errors.As(err, &c) {
		if len(message) == 0 {
			message = c.Error()
		} else {
			message = message + ": " + c.Error()
		}
		return &Error{
			Code:    CodeConflict,
			Message: message,
			Cause:   err,
		}
	}

	// Not found
	var nf domainErrors.NotFound
	if errors.As(err, &nf) {
		if len(message) == 0 {
			message = nf.Error()
		} else {
			message = message + ": " + nf.Error()
		}
		return &Error{
			Code:    CodeNotFound,
			Message: message,
			Cause:   err,
		}
	}

	// Domain-internal (rare)
	var di domainErrors.Internal
	if errors.As(err, &di) {
		return &Error{
			Code:    CodeInternal,
			Message: "Something went wrong. Please try again later.",
			Cause:   err,
		}
	}

	// Everything else is infrastructure / unexpected
	// Convert each message to its own error and join them with the original error
	// to preserve the complete error chain for logging/debugging while keeping the
	// user-facing message generic to avoid leaking internal details
	allErrors := make([]error, 0)
	if len(msg) > 0 {
		for _, m := range msg {
			if m != "" {
				allErrors = append(allErrors, errors.New(m))
			}
		}
	}

	allErrors = append(allErrors, err)
	wrappedErr := errors.Join(allErrors...)
	return &Error{
		Code:    CodeInternal,
		Message: "Something went wrong. Please try again later.",
		Cause:   wrappedErr,
	}
}

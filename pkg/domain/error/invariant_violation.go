package domainerr

import "errors"

var _ InvariantViolation = (*InvariantViolationError)(nil)

type InvariantViolationError struct {
	Entity string
	ID     string
	Cause  error
}

func InvariantViolationErr(entity string, id string) *InvariantViolationError {
	return &InvariantViolationError{
		Entity: entity,
		ID:     id,
	}
}

func (e *InvariantViolationError) WithCause(cause error) *InvariantViolationError {
	e.Cause = cause
	return e
}

func (e *InvariantViolationError) Error() string {
	if e.Cause != nil && len(e.Cause.Error()) > 0 {
		return "failed to enforce invariant: " + e.Entity + " (" + e.ID + "): " + e.Cause.Error()
	}

	return "failed to enforce invariant: " + e.Entity + " (" + e.ID + ")"
}

func (e *InvariantViolationError) Unwrap() error {
	return e.Cause
}

func (e *InvariantViolationError) ErrorCode() string {
	return string(ErrorCodeInvariantViolation)
}

func (e *InvariantViolationError) UserError() string {
	var userErr UserError
	if errors.As(e.Cause, &userErr) {
		return "failed to enforce invariant: " + e.Entity + " (" + e.ID + "): " + userErr.UserError()
	}

	return "failed to enforce invariant: " + e.Entity + " (" + e.ID + ")"
}

func (InvariantViolationError) DomainError()        {}
func (InvariantViolationError) InvariantViolation() {}

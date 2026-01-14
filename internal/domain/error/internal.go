package domainerr

import "fmt"

var _ Internal = (*InternalError)(nil)

type InternalError struct {
	Reason string
	Cause  error
}

func InternalErr(reason string, args ...any) *InternalError {
	return &InternalError{
		Reason: fmt.Sprintf(reason, args...),
	}
}

func (e *InternalError) WithCause(cause error) *InternalError {
	e.Cause = cause
	return e
}

func (e InternalError) Error() string {
	if e.Cause != nil && len(e.Cause.Error()) > 0 {
		return "internal domain error: " + e.Reason + ": " + e.Cause.Error()
	}

	return "internal domain error: " + e.Reason
}

func (e *InternalError) Unwrap() error {
	return e.Cause
}

func (e *InternalError) ErrorCode() string {
	return string(ErrorCodeInternal)
}

func (InternalError) DomainError() {}
func (InternalError) Internal()    {}

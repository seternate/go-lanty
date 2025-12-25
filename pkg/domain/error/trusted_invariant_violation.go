package domainerr

var _ TrustedInvariantViolation = (*TrustedInvariantViolationError)(nil)

type TrustedInvariantViolationError struct {
	Entity string
	ID     string
	Cause  error
}

func TrustedInvariantViolationErr(entity string, id string) *TrustedInvariantViolationError {
	return &TrustedInvariantViolationError{
		Entity: entity,
		ID:     id,
	}
}

func (e *TrustedInvariantViolationError) WithCause(cause error) *TrustedInvariantViolationError {
	e.Cause = cause
	return e
}

func (e *TrustedInvariantViolationError) Error() string {
	if e.Cause != nil && len(e.Cause.Error()) > 0 {
		return "failed to enforce trusted invariant: " + e.Entity + " (" + e.ID + "): " + e.Cause.Error()
	}
	return "failed to enforce trusted invariant: " + e.Entity + " (" + e.ID + ")"
}

func (e *TrustedInvariantViolationError) Unwrap() error {
	return e.Cause
}

func (e *TrustedInvariantViolationError) ErrorCode() string {
	return string(ErrorCodeTrustedInvariantViolation)
}

func (TrustedInvariantViolationError) DomainError()               {}
func (TrustedInvariantViolationError) Internal()                  {}
func (TrustedInvariantViolationError) TrustedInvariantViolation() {}

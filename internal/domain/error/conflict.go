package domainerr

var _ Conflict = (*ConflictError)(nil)

type ConflictError struct {
	Entity string
	ID     string
	Cause  error
}

func ConflictErr(entity string, id string) *ConflictError {
	return &ConflictError{
		Entity: entity,
		ID:     id,
	}
}

func (e *ConflictError) WithCause(cause error) *ConflictError {
	e.Cause = cause
	return e
}

func (e *ConflictError) Error() string {
	if e.Cause != nil && len(e.Cause.Error()) > 0 {
		return e.Entity + " (" + e.ID + ") already exists: " + e.Cause.Error()
	}
	return e.Entity + " (" + e.ID + ") already exists"
}

func (e *ConflictError) Unwrap() error {
	return e.Cause
}

func (e *ConflictError) ErrorCode() string {
	return string(ErrorCodeConflict)
}

func (e *ConflictError) UserError() string {
	return e.Entity + " (" + e.ID + ") already exists"
}

func (ConflictError) DomainError() {}
func (ConflictError) Conflict()    {}

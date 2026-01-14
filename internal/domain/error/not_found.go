package domainerr

var _ NotFound = (*NotFoundError)(nil)

type NotFoundError struct {
	Entity string
	ID     string
	Cause  error
}

func NotFoundErr(entity string, id string) *NotFoundError {
	return &NotFoundError{
		Entity: entity,
		ID:     id,
	}
}

func (e *NotFoundError) WithCause(cause error) *NotFoundError {
	e.Cause = cause
	return e
}

func (e *NotFoundError) Error() string {
	if e.Cause != nil && len(e.Cause.Error()) > 0 {
		return e.Entity + " (" + e.ID + ") not found: " + e.Cause.Error()
	}
	return e.Entity + " (" + e.ID + ") not found"
}

func (e *NotFoundError) Unwrap() error {
	return e.Cause
}

func (e *NotFoundError) ErrorCode() string {
	return string(ErrorCodeNotFound)
}

func (e *NotFoundError) UserError() string {
	return e.Entity + " (" + e.ID + ") not found"
}

func (NotFoundError) DomainError() {}
func (NotFoundError) NotFound()    {}

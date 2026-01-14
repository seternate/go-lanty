package domainerr

type ErrorCode string

const (
	ErrorCodeValidation                ErrorCode = "validation"
	ErrorCodeInvariantViolation        ErrorCode = "invariant_violation"
	ErrorCodeTrustedInvariantViolation ErrorCode = "trusted_invariant_violation"
	ErrorCodeConflict                  ErrorCode = "conflict"
	ErrorCodeNotFound                  ErrorCode = "not_found"
	ErrorCodeInternal                  ErrorCode = "internal"
)

type DomainError interface {
	error
	ErrorCode() string
	DomainError()
}

type UserError interface {
	DomainError
	UserError() string
}

type Internal interface {
	DomainError
	Internal()
}

type Validation interface {
	UserError
	Validation()
}

type InvariantViolation interface {
	UserError
	InvariantViolation()
}

type TrustedInvariantViolation interface {
	Internal
	TrustedInvariantViolation()
}

type NotFound interface {
	UserError
	NotFound()
}

type Conflict interface {
	UserError
	Conflict()
}

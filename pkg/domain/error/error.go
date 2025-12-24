package error

type DomainError interface {
	error
	DomainError()
}

type UserError interface {
	DomainError
	UserError() string
}

type Validation interface {
	UserError
	Validation()
}

type Internal interface {
	DomainError
	Internal()
}

type NotFound interface {
	UserError
	NotFound()
}

type Conflict interface {
	UserError
	Conflict()
}

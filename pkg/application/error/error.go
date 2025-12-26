package apperr

import (
	"errors"

	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

type ErrorDescriptor struct {
	Code    string
	Message string
}

func DescribeError(err error) *ErrorDescriptor {
	if err == nil {
		return nil
	}

	var internalErr domainerr.Internal
	var userErr domainerr.UserError

	if errors.As(err, &internalErr) {
		return &ErrorDescriptor{
			Code:    "internal",
			Message: "Something went wrong. Please try again later.",
		}
	} else if errors.As(err, &userErr) {
		return &ErrorDescriptor{
			Code:    userErr.ErrorCode(),
			Message: userErr.UserError(),
		}
	}

	return nil
}

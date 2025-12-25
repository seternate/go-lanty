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

	var errDescriptor ErrorDescriptor

	var internalErr domainerr.Internal
	var userErr domainerr.UserError

	if errors.As(err, &internalErr) {
		errDescriptor.Message = "Something went wrong. Please try again later."
		errDescriptor.Code = "internal"
	} else if errors.As(err, &userErr) {
		errDescriptor.Message = userErr.UserError()
		errDescriptor.Code = userErr.ErrorCode()
	} else {
		errDescriptor.Message = "Something went wrong unexpectedly. Please try again later."
		errDescriptor.Code = "internal"
	}

	return &errDescriptor
}

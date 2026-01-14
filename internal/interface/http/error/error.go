package errorx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	errormodel "github.com/seternate/go-lanty/internal/interface/http/model/error"
)

type ErrorCode string

var (
	ErrorBadRequest           = ErrorCode("bad request")
	ErrorUnauthorized         = ErrorCode("unauthorized")
	ErrorForbidden            = ErrorCode("forbidden")
	ErrorNotFound             = ErrorCode("not found")
	ErrorConflict             = ErrorCode("conflict")
	ErrorUnsupportedMediaType = ErrorCode("unsupported media type")
)

type Error struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func ErrBadRequest(message string, args ...any) *Error {
	return &Error{
		Code:    ErrorBadRequest,
		Message: fmt.Sprintf(message, args...),
	}
}

func ErrUnauthorized(message string, args ...any) *Error {
	return &Error{
		Code:    ErrorUnauthorized,
		Message: fmt.Sprintf(message, args...),
	}
}

func ErrForbidden(message string, args ...any) *Error {
	return &Error{
		Code:    ErrorForbidden,
		Message: fmt.Sprintf(message, args...),
	}
}

func ErrNotFound(message string, args ...any) *Error {
	return &Error{
		Code:    ErrorNotFound,
		Message: fmt.Sprintf(message, args...),
	}
}

func ErrConflict(message string, args ...any) *Error {
	return &Error{
		Code:    ErrorConflict,
		Message: fmt.Sprintf(message, args...),
	}
}

func ErrUnsupportedMediaType(message string, args ...any) *Error {
	return &Error{
		Code:    ErrorUnsupportedMediaType,
		Message: fmt.Sprintf(message, args...),
	}
}

func (e *Error) WithCause(cause error) *Error {
	e.Cause = cause
	return e
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func HTTPStatusFromError(err error) int {
	var httpErr *Error
	if errors.As(err, &httpErr) {
		switch httpErr.Code {
		case ErrorBadRequest:
			return http.StatusBadRequest
		case ErrorUnauthorized:
			return http.StatusUnauthorized
		case ErrorForbidden:
			return http.StatusForbidden
		case ErrorNotFound:
			return http.StatusNotFound
		case ErrorConflict:
			return http.StatusConflict
		case ErrorUnsupportedMediaType:
			return http.StatusUnsupportedMediaType
		}
		return http.StatusInternalServerError
	}

	var domainErr domainerr.DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.ErrorCode() {
		case string(domainerr.ErrorCodeValidation):
			return http.StatusBadRequest
		case string(domainerr.ErrorCodeInvariantViolation):
			return http.StatusBadRequest
		case string(domainerr.ErrorCodeTrustedInvariantViolation):
			return http.StatusInternalServerError
		case string(domainerr.ErrorCodeConflict):
			return http.StatusConflict
		case string(domainerr.ErrorCodeNotFound):
			return http.StatusNotFound
		case string(domainerr.ErrorCodeInternal):
			return http.StatusInternalServerError
		}
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}

func AbortWithError(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
	ctx.Abort()
}

func DescribeError(err error) *errormodel.ErrorDescriptor {
	if err == nil {
		return nil
	}

	var httpErr *Error
	if errors.As(err, &httpErr) {
		return &errormodel.ErrorDescriptor{
			Code:    string(httpErr.Code),
			Message: httpErr.Error(),
		}
	}

	return nil
}

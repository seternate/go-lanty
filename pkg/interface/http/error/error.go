package errorx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

func HTTPStatusFromError(err error) int {
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
	}

	return http.StatusInternalServerError
}

func AbortWithError(ctx *gin.Context, err error) {
	ctx.Error(err)
	ctx.Abort()
}

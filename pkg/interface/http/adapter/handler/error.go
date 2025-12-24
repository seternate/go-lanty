package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appErrors "github.com/seternate/go-lanty/pkg/application/error"
)

// HTTPStatusFromError converts an application error to an appropriate HTTP status code.
func HTTPStatusFromError(err error) int {
	var appErr *appErrors.Error
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case appErrors.CodeValidation:
			return http.StatusBadRequest
		case appErrors.CodeConflict:
			return http.StatusConflict
		case appErrors.CodeNotFound:
			return http.StatusNotFound
		case appErrors.CodeInternal:
			return http.StatusInternalServerError
		}
	}

	return http.StatusInternalServerError
}

// AbortWithError sets the HTTP status code from the error, records the error, and aborts the request.
func AbortWithError(ctx *gin.Context, err error) {
	ctx.Status(HTTPStatusFromError(err))
	ctx.Error(err)
	ctx.Abort()
}


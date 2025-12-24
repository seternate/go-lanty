package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	domainErrors "github.com/seternate/go-lanty/pkg/domain/error"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()

	// Only handle errors if the response body hasn't been written yet
	// Check Size() instead of Written() because AbortWithError sets status without writing body
	if ctx.Writer.Size() > 0 {
		return
	}

	// Check if there are any errors
	if len(ctx.Errors) == 0 {
		return
	}

	// Get the last error (most specific one)
	err := ctx.Errors.Last()

	// Extract user-facing message from application error
	var userErr domainErrors.UserError
	var code string
	var message string

	if errors.As(err, &userErr) {
		// Use the user-facing message from the application error
		message = userErr.UserError()
		code = "usererror TO BE IMPLEMENTED"
	} else {
		// For non-application errors, use a generic message
		message = "An error occurred while processing your request"
		code = "UNKNOWN"
	}

	// Write error response using the status code already set by ctx.AbortWithError()
	// Use AbortWithStatusJSON since context is already aborted
	ctx.JSON(-1, gin.H{
		"error": ErrorResponse{
			Message: message,
			Code:    code,
		}},
	)
}

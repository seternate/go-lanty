package middleware

import (
	"github.com/gin-gonic/gin"
	apperr "github.com/seternate/go-lanty/pkg/application/error"
	errorx "github.com/seternate/go-lanty/pkg/interface/http/error"
)

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()

	if ctx.Writer.Size() > 0 {
		return
	}

	if len(ctx.Errors) == 0 {
		return
	}

	err := ctx.Errors.Last()

	errDescriptor := apperr.DescribeError(err)
	if errDescriptor == nil {
		return
	}

	ctx.JSON(errorx.HTTPStatusFromError(err), gin.H{"error": *errDescriptor})
}

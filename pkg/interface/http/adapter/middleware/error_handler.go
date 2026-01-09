package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperr "github.com/seternate/go-lanty/pkg/application/error"
	errorx "github.com/seternate/go-lanty/pkg/interface/http/error"
	"github.com/seternate/go-lanty/pkg/interface/http/model"
)

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()

	if ctx.Writer.Size() > 0 {
		return
	}

	if len(ctx.Errors) == 0 {
		return
	}

	err := ctx.Errors.Last().Err
	if err == nil {
		return
	}

	httpErrDescriptor := errorx.DescribeError(err)
	if httpErrDescriptor != nil {
		ctx.JSON(errorx.HTTPStatusFromError(err), model.ErrorResponse{
			Error: model.ErrorDescriptor{
				Code:    httpErrDescriptor.Code,
				Message: httpErrDescriptor.Message,
			},
		})
		return
	}

	appErrDescriptor := apperr.DescribeError(err)
	if appErrDescriptor != nil {
		ctx.JSON(errorx.HTTPStatusFromError(err), model.ErrorResponse{
			Error: model.ErrorDescriptor{
				Code:    appErrDescriptor.Code,
				Message: appErrDescriptor.Message,
			},
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
		Error: model.ErrorDescriptor{
			Code:    "internal",
			Message: "Something went wrong unexpectedly. Please try again later.",
		},
	})
}

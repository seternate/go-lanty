package middleware

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/hlog"
)

func Logger(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()

	logger := hlog.FromRequest(ctx.Request).With().
		Int("status", ctx.Writer.Status()).
		Str("size", humanize.IBytes(uint64(ctx.Writer.Size()))).
		Stringer("duration", time.Since(start)).Logger()

	if len(ctx.Errors) > 0 && ctx.Writer.Status() >= 400 {
		logger.Error().Str("error", ctx.Errors.Last().Error()).Send()
		return
	} else if ctx.Writer.Status() >= 400 {
		logger.Error().Send()
		return
	}

	logger.Info().Send()
}

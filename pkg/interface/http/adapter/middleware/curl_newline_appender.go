package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// It appends a newline character to the response for 'curl/*' useragent.
// Without this curl request stdout output could be overwritten by some shells (e.g. zsh)
func CurlNewlineAppender(ctx *gin.Context) {
	ctx.Next()
	useragent := ctx.Request.UserAgent()
	if strings.HasPrefix(useragent, "curl") {
		_, _ = ctx.Writer.WriteString("\n")
	}
}

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// MiddlewareNewlineAppender returns a Gin middleware that should be executed just after the handler.
// It appends a newline character to the response for 'curl/*' useragent.
func MiddlewareNewlineAppender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()                           // First call the next handler in the chain.
		useragent := ctx.Request.UserAgent() // Then perform the actual job of appending a newline character.
		if strings.HasPrefix(useragent, "curl") {
			_, _ = ctx.Writer.WriteString("\n")
		}
	}
}

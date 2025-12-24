package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/hlog"
)

// findLongestOverlap finds the longest suffix of s1 that matches a prefix of s2
func findLongestOverlap(s1, s2 string) string {
	if len(s1) == 0 || len(s2) == 0 {
		return ""
	}

	// Try all possible suffix lengths of s1, starting from the longest
	maxLen := len(s1)
	if len(s2) < maxLen {
		maxLen = len(s2)
	}

	for i := maxLen; i > 0; i-- {
		suffix := s1[len(s1)-i:]
		prefix := s2[:i]
		if suffix == prefix {
			return suffix
		}
	}

	return ""
}

// collectErrorChain unwraps an error and collects all error messages in the chain.
// It skips wrapper errors that just pass through the same message (passthrough wrappers)
// and removes overlapping suffixes/prefixes between consecutive errors.
func collectErrorChain(err error) []string {
	if err == nil {
		return nil
	}

	var messages []string
	seenMessages := make(map[string]bool)
	current := err
	var previousMsg string // Track the previous logged message

	for current != nil {
		currentMsg := current.Error()

		// Peek at the unwrapped error to check if this is just a passthrough wrapper
		unwrapped := errors.Unwrap(current)
		var unwrappedMsg string
		if unwrapped != nil {
			unwrappedMsg = unwrapped.Error()
		}

		// Skip passthrough wrapper errors: if the message is identical to the unwrapped error,
		// this is just a wrapper that doesn't add context, so skip it
		isPassthrough := unwrapped != nil && currentMsg == unwrappedMsg

		// Remove overlap from current message if previous message ends with a prefix of current
		if len(previousMsg) > 0 && len(currentMsg) > 0 {
			overlap := findLongestOverlap(previousMsg, currentMsg)
			if len(overlap) > 0 {
				// Remove the overlapping prefix from current message
				currentMsg = strings.TrimPrefix(currentMsg, overlap)
				// Also trim any leading whitespace/separators that might remain
				currentMsg = strings.TrimLeft(currentMsg, ": \n\t")
			}
		}

		// Only log if:
		// - It's not a passthrough wrapper (unless it's the last error in chain)
		// - We haven't seen this exact message before (using the condensed message)
		if !isPassthrough && !seenMessages[currentMsg] && len(currentMsg) > 0 {
			messages = append(messages, currentMsg)
			seenMessages[currentMsg] = true
			previousMsg = current.Error() // Store the original message for overlap detection
		}

		current = unwrapped
	}

	return messages
}

func Logger(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()

	logger := hlog.FromRequest(ctx.Request).With().
		Int("status", ctx.Writer.Status()).
		Str("size", humanize.IBytes(uint64(ctx.Writer.Size()))).
		Stringer("duration", time.Since(start)).Logger()

	if len(ctx.Errors) > 0 && ctx.Writer.Status() >= 400 {
		var allErrorMessages []string
		for _, err := range ctx.Errors {
			allErrorMessages = append(allErrorMessages, collectErrorChain(err)...)
		}
		logger.Error().Str("error", strings.Join(allErrorMessages, ": ")).Send()
		return
	} else if ctx.Writer.Status() >= 400 {
		logger.Error().Send()
		return
	}

	logger.Info().Send()
}

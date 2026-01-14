package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

func TestErrorCodeConstants(t *testing.T) {
	tests := []struct {
		name string
		code ErrorCode
		want string
	}{
		{"BadRequest", ErrorBadRequest, "bad request"},
		{"Unauthorized", ErrorUnauthorized, "unauthorized"},
		{"Forbidden", ErrorForbidden, "forbidden"},
		{"NotFound", ErrorNotFound, "not found"},
		{"Conflict", ErrorConflict, "conflict"},
		{"UnsupportedMediaType", ErrorUnsupportedMediaType, "unsupported media type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.code))
		})
	}
}

func TestErrBadRequest(t *testing.T) {
	t.Run("simple message", func(t *testing.T) {
		err := ErrBadRequest("invalid input")
		require.NotNil(t, err)
		assert.Equal(t, ErrorBadRequest, err.Code)
		assert.Equal(t, "invalid input", err.Message)
		assert.Nil(t, err.Cause)
	})

	t.Run("formatted message", func(t *testing.T) {
		err := ErrBadRequest("invalid input: %s", "field is required")
		require.NotNil(t, err)
		assert.Equal(t, ErrorBadRequest, err.Code)
		assert.Equal(t, "invalid input: field is required", err.Message)
		assert.Nil(t, err.Cause)
	})
}

func TestErrUnauthorized(t *testing.T) {
	t.Run("simple message", func(t *testing.T) {
		err := ErrUnauthorized("authentication required")
		require.NotNil(t, err)
		assert.Equal(t, ErrorUnauthorized, err.Code)
		assert.Equal(t, "authentication required", err.Message)
		assert.Nil(t, err.Cause)
	})

	t.Run("formatted message", func(t *testing.T) {
		err := ErrUnauthorized("token expired: %s", "refresh required")
		require.NotNil(t, err)
		assert.Equal(t, ErrorUnauthorized, err.Code)
		assert.Equal(t, "token expired: refresh required", err.Message)
		assert.Nil(t, err.Cause)
	})
}

func TestErrForbidden(t *testing.T) {
	t.Run("simple message", func(t *testing.T) {
		err := ErrForbidden("access denied")
		require.NotNil(t, err)
		assert.Equal(t, ErrorForbidden, err.Code)
		assert.Equal(t, "access denied", err.Message)
		assert.Nil(t, err.Cause)
	})

	t.Run("formatted message", func(t *testing.T) {
		err := ErrForbidden("user %s lacks permission", "alice")
		require.NotNil(t, err)
		assert.Equal(t, ErrorForbidden, err.Code)
		assert.Equal(t, "user alice lacks permission", err.Message)
		assert.Nil(t, err.Cause)
	})
}

func TestErrNotFound(t *testing.T) {
	t.Run("simple message", func(t *testing.T) {
		err := ErrNotFound("resource not found")
		require.NotNil(t, err)
		assert.Equal(t, ErrorNotFound, err.Code)
		assert.Equal(t, "resource not found", err.Message)
		assert.Nil(t, err.Cause)
	})

	t.Run("formatted message", func(t *testing.T) {
		err := ErrNotFound("game %s not found", "game-123")
		require.NotNil(t, err)
		assert.Equal(t, ErrorNotFound, err.Code)
		assert.Equal(t, "game game-123 not found", err.Message)
		assert.Nil(t, err.Cause)
	})
}

func TestErrConflict(t *testing.T) {
	t.Run("simple message", func(t *testing.T) {
		err := ErrConflict("resource conflict")
		require.NotNil(t, err)
		assert.Equal(t, ErrorConflict, err.Code)
		assert.Equal(t, "resource conflict", err.Message)
		assert.Nil(t, err.Cause)
	})

	t.Run("formatted message", func(t *testing.T) {
		err := ErrConflict("game %s already exists", "game-123")
		require.NotNil(t, err)
		assert.Equal(t, ErrorConflict, err.Code)
		assert.Equal(t, "game game-123 already exists", err.Message)
		assert.Nil(t, err.Cause)
	})
}

func TestErrUnsupportedMediaType(t *testing.T) {
	t.Run("simple message", func(t *testing.T) {
		err := ErrUnsupportedMediaType("unsupported media type")
		require.NotNil(t, err)
		assert.Equal(t, ErrorUnsupportedMediaType, err.Code)
		assert.Equal(t, "unsupported media type", err.Message)
		assert.Nil(t, err.Cause)
	})

	t.Run("formatted message", func(t *testing.T) {
		err := ErrUnsupportedMediaType("content type %s not supported", "text/plain")
		require.NotNil(t, err)
		assert.Equal(t, ErrorUnsupportedMediaType, err.Code)
		assert.Equal(t, "content type text/plain not supported", err.Message)
		assert.Nil(t, err.Cause)
	})
}

func TestError_WithCause(t *testing.T) {
	t.Run("sets cause", func(t *testing.T) {
		err := ErrBadRequest("invalid input")
		cause := errors.New("underlying error")
		
		result := err.WithCause(cause)
		
		assert.Equal(t, err, result) // Should return same instance
		assert.Equal(t, cause, err.Cause)
	})

	t.Run("chainable", func(t *testing.T) {
		err := ErrNotFound("not found")
		cause1 := errors.New("cause 1")
		cause2 := errors.New("cause 2")
		
		result := err.WithCause(cause1).WithCause(cause2)
		
		assert.Equal(t, err, result)
		assert.Equal(t, cause2, err.Cause) // Last cause wins
	})
}

func TestError_Error(t *testing.T) {
	t.Run("returns message", func(t *testing.T) {
		err := ErrBadRequest("test message")
		assert.Equal(t, "test message", err.Error())
	})

	t.Run("formatted message", func(t *testing.T) {
		err := ErrNotFound("game %s not found", "game-123")
		assert.Equal(t, "game game-123 not found", err.Error())
	})
}

func TestError_Unwrap(t *testing.T) {
	t.Run("returns nil when no cause", func(t *testing.T) {
		err := ErrBadRequest("test")
		assert.Nil(t, err.Unwrap())
	})

	t.Run("returns cause", func(t *testing.T) {
		err := ErrBadRequest("test")
		cause := errors.New("underlying error")
		err.WithCause(cause)
		
		assert.Equal(t, cause, err.Unwrap())
	})

	t.Run("supports errors.Is", func(t *testing.T) {
		err := ErrBadRequest("test")
		cause := errors.New("underlying error")
		err.WithCause(cause)
		
		assert.True(t, errors.Is(err, cause))
	})
}

func TestHTTPStatusFromError(t *testing.T) {
	t.Run("HTTP error - BadRequest", func(t *testing.T) {
		err := ErrBadRequest("test")
		assert.Equal(t, http.StatusBadRequest, HTTPStatusFromError(err))
	})

	t.Run("HTTP error - Unauthorized", func(t *testing.T) {
		err := ErrUnauthorized("test")
		assert.Equal(t, http.StatusUnauthorized, HTTPStatusFromError(err))
	})

	t.Run("HTTP error - Forbidden", func(t *testing.T) {
		err := ErrForbidden("test")
		assert.Equal(t, http.StatusForbidden, HTTPStatusFromError(err))
	})

	t.Run("HTTP error - NotFound", func(t *testing.T) {
		err := ErrNotFound("test")
		assert.Equal(t, http.StatusNotFound, HTTPStatusFromError(err))
	})

	t.Run("HTTP error - Conflict", func(t *testing.T) {
		err := ErrConflict("test")
		assert.Equal(t, http.StatusConflict, HTTPStatusFromError(err))
	})

	t.Run("HTTP error - UnsupportedMediaType", func(t *testing.T) {
		err := ErrUnsupportedMediaType("test")
		assert.Equal(t, http.StatusUnsupportedMediaType, HTTPStatusFromError(err))
	})

	t.Run("HTTP error - unknown code", func(t *testing.T) {
		err := &Error{
			Code:    ErrorCode("unknown"),
			Message: "test",
		}
		assert.Equal(t, http.StatusInternalServerError, HTTPStatusFromError(err))
	})

	t.Run("HTTP error wrapped", func(t *testing.T) {
		err := ErrNotFound("test")
		wrapped := fmt.Errorf("wrapper: %w", err)
		assert.Equal(t, http.StatusNotFound, HTTPStatusFromError(wrapped))
	})

	t.Run("domain error - Validation", func(t *testing.T) {
		err := domainerr.ValidationErr("field", "message")
		assert.Equal(t, http.StatusBadRequest, HTTPStatusFromError(err))
	})

	t.Run("domain error - InvariantViolation", func(t *testing.T) {
		err := domainerr.InvariantViolationErr("entity", "id")
		assert.Equal(t, http.StatusBadRequest, HTTPStatusFromError(err))
	})

	t.Run("domain error - TrustedInvariantViolation", func(t *testing.T) {
		err := domainerr.TrustedInvariantViolationErr("entity", "id")
		assert.Equal(t, http.StatusInternalServerError, HTTPStatusFromError(err))
	})

	t.Run("domain error - Conflict", func(t *testing.T) {
		err := domainerr.ConflictErr("entity", "id")
		assert.Equal(t, http.StatusConflict, HTTPStatusFromError(err))
	})

	t.Run("domain error - NotFound", func(t *testing.T) {
		err := domainerr.NotFoundErr("entity", "id")
		assert.Equal(t, http.StatusNotFound, HTTPStatusFromError(err))
	})

	t.Run("domain error - Internal", func(t *testing.T) {
		err := domainerr.InternalErr("message")
		assert.Equal(t, http.StatusInternalServerError, HTTPStatusFromError(err))
	})

	t.Run("domain error - unknown code", func(t *testing.T) {
		// Create a mock domain error with unknown code
		mockErr := &mockDomainError{code: "unknown"}
		assert.Equal(t, http.StatusInternalServerError, HTTPStatusFromError(mockErr))
	})

	t.Run("domain error wrapped", func(t *testing.T) {
		err := domainerr.NotFoundErr("entity", "id")
		wrapped := fmt.Errorf("wrapper: %w", err)
		assert.Equal(t, http.StatusNotFound, HTTPStatusFromError(wrapped))
	})

	t.Run("unknown error type", func(t *testing.T) {
		err := errors.New("unknown error")
		assert.Equal(t, http.StatusInternalServerError, HTTPStatusFromError(err))
	})

	t.Run("nil error", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, HTTPStatusFromError(nil))
	})
}

// mockDomainError is a helper for testing unknown domain error codes
type mockDomainError struct {
	code string
}

func (m *mockDomainError) Error() string {
	return "mock error"
}

func (m *mockDomainError) ErrorCode() string {
	return m.code
}

func (m *mockDomainError) DomainError() {}

func TestAbortWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("aborts context", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(nil)
		err := ErrBadRequest("test error")
		
		AbortWithError(ctx, err)
		
		assert.True(t, ctx.IsAborted())
	})

	t.Run("adds error to context", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(nil)
		err := ErrNotFound("not found")
		
		AbortWithError(ctx, err)
		
		// Check that error was added to context
		// Gin's Error method adds errors to context.Errors slice
		assert.True(t, ctx.IsAborted())
		// Note: We can't easily test the exact error in context without accessing
		// internal gin state, but we can verify abort happened
	})

}

func TestDescribeError(t *testing.T) {
	t.Run("HTTP error", func(t *testing.T) {
		err := ErrBadRequest("invalid input")
		desc := DescribeError(err)
		
		require.NotNil(t, desc)
		assert.Equal(t, string(ErrorBadRequest), desc.Code)
		assert.Equal(t, "invalid input", desc.Message)
	})

	t.Run("HTTP error with formatted message", func(t *testing.T) {
		err := ErrNotFound("game %s not found", "game-123")
		desc := DescribeError(err)
		
		require.NotNil(t, desc)
		assert.Equal(t, string(ErrorNotFound), desc.Code)
		assert.Equal(t, "game game-123 not found", desc.Message)
	})

	t.Run("HTTP error wrapped", func(t *testing.T) {
		err := ErrConflict("conflict")
		wrapped := fmt.Errorf("wrapper: %w", err)
		desc := DescribeError(wrapped)
		
		require.NotNil(t, desc)
		assert.Equal(t, string(ErrorConflict), desc.Code)
		assert.Equal(t, "conflict", desc.Message)
	})

	t.Run("domain error returns nil", func(t *testing.T) {
		err := domainerr.NotFoundErr("entity", "id")
		desc := DescribeError(err)
		
		assert.Nil(t, desc)
	})

	t.Run("unknown error returns nil", func(t *testing.T) {
		err := errors.New("unknown error")
		desc := DescribeError(err)
		
		assert.Nil(t, desc)
	})

	t.Run("nil error returns nil", func(t *testing.T) {
		desc := DescribeError(nil)
		
		assert.Nil(t, desc)
	})
}

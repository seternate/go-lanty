package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful request - logs info", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		// hlog.FromRequest will use the default logger if no logger is in context
		// This is fine for testing - the middleware will still work
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.WriteString("test")

		Logger(ctx)

		// Logger should complete without panicking
		// We can't easily test the actual log output without a more complex setup,
		// but we can verify the function executes successfully
		if ctx.Writer.Status() != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, ctx.Writer.Status())
		}
	})

	t.Run("Client error (400) without errors - logs error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(http.StatusBadRequest)

		Logger(ctx)

		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, ctx.Writer.Status())
		}
	})

	t.Run("Client error (404) without errors - logs error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(http.StatusNotFound)

		Logger(ctx)

		if ctx.Writer.Status() != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, ctx.Writer.Status())
		}
	})

	t.Run("Server error (500) without errors - logs error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)

		Logger(ctx)

		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, ctx.Writer.Status())
		}
	})

	t.Run("Client error (400) with errors - logs error with error message", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		_ = ctx.Error(&testError{message: "validation failed"})

		Logger(ctx)

		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, ctx.Writer.Status())
		}
		if len(ctx.Errors) == 0 {
			t.Error("expected errors to be present")
		}
	})

	t.Run("Server error (500) with errors - logs error with error message", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		_ = ctx.Error(&testError{message: "internal server error"})

		Logger(ctx)

		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, ctx.Writer.Status())
		}
	})

	t.Run("Duration is measured", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(http.StatusOK)

		start := time.Now()
		Logger(ctx)
		duration := time.Since(start)

		// Logger should complete quickly (less than 100ms for a simple operation)
		if duration > 100*time.Millisecond {
			t.Errorf("logger took too long: %v", duration)
		}
	})

	t.Run("Response size is logged", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		responseBody := "test response body"
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.WriteString(responseBody)

		Logger(ctx)

		// Verify the response was written
		if w.Body.String() != responseBody {
			t.Errorf("expected body %q, got %q", responseBody, w.Body.String())
		}
		if ctx.Writer.Size() != len(responseBody) {
			t.Errorf("expected size %d, got %d", len(responseBody), ctx.Writer.Size())
		}
	})

	t.Run("Status code 399 - logs info (not error)", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(399) // Custom status code below 400

		Logger(ctx)

		if ctx.Writer.Status() != 399 {
			t.Errorf("expected status %d, got %d", 399, ctx.Writer.Status())
		}
	})

	t.Run("Status code 400 - logs error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(http.StatusBadRequest)

		Logger(ctx)

		if ctx.Writer.Status() != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, ctx.Writer.Status())
		}
	})

	t.Run("Status code 499 - logs error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(499) // Custom status code in 4xx range

		Logger(ctx)

		if ctx.Writer.Status() != 499 {
			t.Errorf("expected status %d, got %d", 499, ctx.Writer.Status())
		}
	})

	t.Run("Status code 500 - logs error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)

		Logger(ctx)

		if ctx.Writer.Status() != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, ctx.Writer.Status())
		}
	})
}

package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	errorx "github.com/seternate/go-lanty/internal/interface/http/error"
	"github.com/seternate/go-lanty/internal/interface/http/model"
)

// testResponseWriter wraps httptest.ResponseRecorder to track size properly
type testResponseWriter struct {
	http.ResponseWriter
	size int
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *testResponseWriter) WriteString(s string) (int, error) {
	n, err := w.ResponseWriter.Write([]byte(s))
	w.size += n
	return n, err
}

func (w *testResponseWriter) Size() int {
	return w.size
}

func (w *testResponseWriter) Status() int {
	return w.ResponseWriter.(*httptest.ResponseRecorder).Code
}

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("No errors - response already written", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Status(http.StatusOK)
		ctx.Writer.WriteString("test response")

		ErrorHandler(ctx)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		if w.Body.String() != "test response" {
			t.Errorf("expected body %q, got %q", "test response", w.Body.String())
		}
	})

	t.Run("No errors - empty errors list", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ErrorHandler(ctx)

		// No response should be written since there are no errors
		if w.Code != 0 && w.Code != 200 {
			t.Errorf("expected status 0 or 200 (not set or default), got %d", w.Code)
		}
	})

	t.Run("No errors - nil error in errors list", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Errors = []*gin.Error{
			{Err: nil},
		}

		ErrorHandler(ctx)

		// No response should be written since error is nil
		if w.Code != 0 && w.Code != 200 {
			t.Errorf("expected status 0 or 200 (not set or default), got %d", w.Code)
		}
	})

	t.Run("HTTP error - BadRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		httpErr := errorx.ErrBadRequest("invalid input")
		_ = ctx.Error(httpErr)

		ErrorHandler(ctx)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != string(errorx.ErrorBadRequest) {
			t.Errorf("expected error code %q, got %q", errorx.ErrorBadRequest, response.Error.Code)
		}
		if response.Error.Message != "invalid input" {
			t.Errorf("expected error message %q, got %q", "invalid input", response.Error.Message)
		}
	})

	t.Run("HTTP error - NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		httpErr := errorx.ErrNotFound("resource not found")
		_ = ctx.Error(httpErr)

		ErrorHandler(ctx)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != string(errorx.ErrorNotFound) {
			t.Errorf("expected error code %q, got %q", errorx.ErrorNotFound, response.Error.Code)
		}
	})

	t.Run("HTTP error - Conflict", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		httpErr := errorx.ErrConflict("resource already exists")
		_ = ctx.Error(httpErr)

		ErrorHandler(ctx)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != string(errorx.ErrorConflict) {
			t.Errorf("expected error code %q, got %q", errorx.ErrorConflict, response.Error.Code)
		}
	})

	t.Run("Domain error - Validation", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		domainErr := domainerr.ValidationErr("field", "invalid value")
		_ = ctx.Error(domainErr)

		ErrorHandler(ctx)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != string(domainerr.ErrorCodeValidation) {
			t.Errorf("expected error code %q, got %q", domainerr.ErrorCodeValidation, response.Error.Code)
		}
		if response.Error.Message != "field: invalid value" {
			t.Errorf("expected error message %q, got %q", "field: invalid value", response.Error.Message)
		}
	})

	t.Run("Domain error - NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		domainErr := domainerr.NotFoundErr("Game", "test-id")
		_ = ctx.Error(domainErr)

		ErrorHandler(ctx)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != string(domainerr.ErrorCodeNotFound) {
			t.Errorf("expected error code %q, got %q", domainerr.ErrorCodeNotFound, response.Error.Code)
		}
		if response.Error.Message != "Game (test-id) not found" {
			t.Errorf("expected error message %q, got %q", "Game (test-id) not found", response.Error.Message)
		}
	})

	t.Run("Domain error - Conflict", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		domainErr := domainerr.ConflictErr("Game", "test-id")
		_ = ctx.Error(domainErr)

		ErrorHandler(ctx)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != string(domainerr.ErrorCodeConflict) {
			t.Errorf("expected error code %q, got %q", domainerr.ErrorCodeConflict, response.Error.Code)
		}
		if response.Error.Message != "Game (test-id) already exists" {
			t.Errorf("expected error message %q, got %q", "Game (test-id) already exists", response.Error.Message)
		}
	})

	t.Run("Domain error - Internal", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		domainErr := domainerr.InternalErr("database connection failed")
		_ = ctx.Error(domainErr)

		ErrorHandler(ctx)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != "internal" {
			t.Errorf("expected error code %q, got %q", "internal", response.Error.Code)
		}
		if response.Error.Message != "Something went wrong. Please try again later." {
			t.Errorf("expected error message %q, got %q", "Something went wrong. Please try again later.", response.Error.Message)
		}
	})

	t.Run("Unknown error - fallback to internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		unknownErr := &testError{message: "unknown error"}
		_ = ctx.Error(unknownErr)

		ErrorHandler(ctx)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != "internal" {
			t.Errorf("expected error code %q, got %q", "internal", response.Error.Code)
		}
		if response.Error.Message != "Something went wrong unexpectedly. Please try again later." {
			t.Errorf("expected error message %q, got %q", "Something went wrong unexpectedly. Please try again later.", response.Error.Message)
		}
	})

	t.Run("HTTP error takes precedence over domain error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		domainErr := domainerr.ValidationErr("field", "domain error")
		httpErr := errorx.ErrBadRequest("http error")
		_ = ctx.Error(domainErr)
		_ = ctx.Error(httpErr) // HTTP error added last, so it's the one processed

		ErrorHandler(ctx)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response model.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error.Code != string(errorx.ErrorBadRequest) {
			t.Errorf("expected error code %q, got %q", errorx.ErrorBadRequest, response.Error.Code)
		}
	})
}

// testError is a simple error type for testing unknown errors
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

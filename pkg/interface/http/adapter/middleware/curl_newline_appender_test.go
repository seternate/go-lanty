package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCurlNewlineAppender(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Curl user agent - appends newline", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Request.Header.Set("User-Agent", "curl/7.68.0")
		ctx.Writer.WriteString("test response")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		if !strings.HasSuffix(body, "\n") {
			t.Errorf("expected body to end with newline, got %q", body)
		}
		if body != "test response\n" {
			t.Errorf("expected body %q, got %q", "test response\n", body)
		}
	})

	t.Run("Curl user agent with version - appends newline", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Request.Header.Set("User-Agent", "curl/8.1.2")
		ctx.Writer.WriteString("response")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		if !strings.HasSuffix(body, "\n") {
			t.Errorf("expected body to end with newline, got %q", body)
		}
	})

	t.Run("Curl user agent with additional info - appends newline", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Request.Header.Set("User-Agent", "curl/7.68.0 (x86_64-pc-linux-gnu) libcurl/7.68.0")
		ctx.Writer.WriteString("data")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		if !strings.HasSuffix(body, "\n") {
			t.Errorf("expected body to end with newline, got %q", body)
		}
	})

	t.Run("Non-curl user agent - no newline appended", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Request.Header.Set("User-Agent", "Mozilla/5.0")
		ctx.Writer.WriteString("test response")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		if strings.HasSuffix(body, "\n") {
			t.Errorf("expected body not to end with newline, got %q", body)
		}
		if body != "test response" {
			t.Errorf("expected body %q, got %q", "test response", body)
		}
	})

	t.Run("Empty user agent - no newline appended", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Writer.WriteString("test response")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		if strings.HasSuffix(body, "\n") {
			t.Errorf("expected body not to end with newline, got %q", body)
		}
		if body != "test response" {
			t.Errorf("expected body %q, got %q", "test response", body)
		}
	})

	t.Run("User agent starting with curl but not exact match - appends newline", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Request.Header.Set("User-Agent", "curlbot/1.0")
		ctx.Writer.WriteString("test")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		if !strings.HasSuffix(body, "\n") {
			t.Errorf("expected body to end with newline (prefix match), got %q", body)
		}
	})

	t.Run("Case sensitive - lowercase curl", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Request.Header.Set("User-Agent", "curl/7.68.0")
		ctx.Writer.WriteString("test")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		if !strings.HasSuffix(body, "\n") {
			t.Errorf("expected body to end with newline, got %q", body)
		}
	})

	t.Run("Case sensitive - uppercase CURL does not match", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Request.Header.Set("User-Agent", "CURL/7.68.0")
		ctx.Writer.WriteString("test")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		// strings.HasPrefix is case-sensitive, so "CURL" won't match "curl"
		if strings.HasSuffix(body, "\n") {
			t.Errorf("expected body not to end with newline (case sensitive), got %q", body)
		}
	})

	t.Run("Empty response body - still appends newline for curl", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/test", nil)
		ctx.Request.Header.Set("User-Agent", "curl/7.68.0")

		CurlNewlineAppender(ctx)

		body := w.Body.String()
		if body != "\n" {
			t.Errorf("expected body to be just newline, got %q", body)
		}
	})
}

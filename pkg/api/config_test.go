package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("valid HTTP URL", func(t *testing.T) {
		cfg, err := NewConfig("http://example.com")
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "http", cfg.BaseURL.Scheme)
		assert.Equal(t, "example.com", cfg.BaseURL.Host)
		assert.Equal(t, 30*time.Second, cfg.Timeout)
	})

	t.Run("valid HTTPS URL", func(t *testing.T) {
		cfg, err := NewConfig("https://api.example.com")
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "https", cfg.BaseURL.Scheme)
		assert.Equal(t, "api.example.com", cfg.BaseURL.Host)
	})

	t.Run("valid URL with path", func(t *testing.T) {
		cfg, err := NewConfig("https://example.com/api/v1")
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "/api/v1", cfg.BaseURL.Path)
	})

	t.Run("invalid URL - no scheme", func(t *testing.T) {
		cfg, err := NewConfig("example.com")
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "baseURL must contain a scheme")
	})

	t.Run("invalid URL - unsupported scheme", func(t *testing.T) {
		cfg, err := NewConfig("ftp://example.com")
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "scheme must be http or https")
	})

	t.Run("invalid URL - no host", func(t *testing.T) {
		cfg, err := NewConfig("http://")
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "baseURL must contain a host")
	})

	t.Run("invalid URL - malformed", func(t *testing.T) {
		cfg, err := NewConfig("://invalid")
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "failed to parse baseURL")
	})
}

func TestWithTimeout(t *testing.T) {
	cfg, err := NewConfig("http://example.com", WithTimeout(60*time.Second))
	require.NoError(t, err)
	assert.Equal(t, 60*time.Second, cfg.Timeout)
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	cfg, err := NewConfig("http://example.com", WithHTTPClient(customClient))
	require.NoError(t, err)
	assert.Equal(t, customClient, cfg.HTTPClient)
}

func TestWithHeader(t *testing.T) {
	cfg, err := NewConfig("http://example.com",
		WithHeader("Authorization", "Bearer token123"),
		WithHeader("X-Custom", "value"),
	)
	require.NoError(t, err)
	assert.NotNil(t, cfg.Headers)
	assert.Equal(t, "Bearer token123", cfg.Headers["Authorization"])
	assert.Equal(t, "value", cfg.Headers["X-Custom"])
}

func TestWithHeaders(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token123",
		"X-Custom":      "value",
	}

	cfg, err := NewConfig("http://example.com", WithHeaders(headers))
	require.NoError(t, err)
	assert.NotNil(t, cfg.Headers)
	assert.Equal(t, "Bearer token123", cfg.Headers["Authorization"])
	assert.Equal(t, "value", cfg.Headers["X-Custom"])
}

func TestWithHeaders_OverwritesExisting(t *testing.T) {
	headers1 := map[string]string{
		"Authorization": "Bearer token123",
		"X-Custom":      "value1",
	}
	headers2 := map[string]string{
		"X-Custom": "value2",
		"X-New":    "newvalue",
	}

	cfg, err := NewConfig("http://example.com",
		WithHeaders(headers1),
		WithHeaders(headers2),
	)
	require.NoError(t, err)
	assert.NotNil(t, cfg.Headers)
	assert.Equal(t, "Bearer token123", cfg.Headers["Authorization"])
	assert.Equal(t, "value2", cfg.Headers["X-Custom"]) // overwritten
	assert.Equal(t, "newvalue", cfg.Headers["X-New"])
}

func TestConfig_OptionChaining(t *testing.T) {
	customClient := &http.Client{Timeout: 5 * time.Second}
	cfg, err := NewConfig("https://api.example.com",
		WithTimeout(45*time.Second),
		WithHTTPClient(customClient),
		WithHeader("Authorization", "Bearer token"),
		WithHeader("Content-Type", "application/json"),
	)
	require.NoError(t, err)
	assert.Equal(t, "https", cfg.BaseURL.Scheme)
	assert.Equal(t, "api.example.com", cfg.BaseURL.Host)
	assert.Equal(t, customClient, cfg.HTTPClient)
	assert.Equal(t, "Bearer token", cfg.Headers["Authorization"])
	assert.Equal(t, "application/json", cfg.Headers["Content-Type"])
}

func TestConfig_WithHeaderCreatesMap(t *testing.T) {
	cfg, err := NewConfig("http://example.com", WithHeader("X-Test", "value"))
	require.NoError(t, err)
	assert.NotNil(t, cfg.Headers)
	assert.Equal(t, "value", cfg.Headers["X-Test"])
}

func TestConfig_WithHeadersCreatesMap(t *testing.T) {
	cfg, err := NewConfig("http://example.com", WithHeaders(map[string]string{"X-Test": "value"}))
	require.NoError(t, err)
	assert.NotNil(t, cfg.Headers)
	assert.Equal(t, "value", cfg.Headers["X-Test"])
}

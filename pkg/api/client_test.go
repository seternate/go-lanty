package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seternate/go-lanty/pkg/api/paths"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, err := New("http://example.com")
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.gamesClient)
		assert.NotNil(t, client.usersClient)
		assert.NotNil(t, client.httpClient)
		assert.NotNil(t, client.config)
	})

	t.Run("invalid URL", func(t *testing.T) {
		client, err := New("invalid-url")
		require.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("with options", func(t *testing.T) {
		customClient := &http.Client{}
		client, err := New("http://example.com",
			WithTimeout(60),
			WithHTTPClient(customClient),
			WithHeader("X-Test", "value"),
		)
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, customClient, client.httpClient)
		assert.Equal(t, "value", client.config.Headers["X-Test"])
	})
}

func TestClient_BuildURL(t *testing.T) {
	client, err := New("http://example.com")
	require.NoError(t, err)

	t.Run("path without params", func(t *testing.T) {
		url, err := client.BuildURL(paths.GetGames, nil)
		require.NoError(t, err)
		assert.Equal(t, "http://example.com/games", url)
	})

	t.Run("path with params", func(t *testing.T) {
		url, err := client.BuildURL(paths.GetGame, paths.GameParams{Slug: "test-game"})
		require.NoError(t, err)
		assert.Equal(t, "http://example.com/games/test-game", url)
	})

	t.Run("path with leading slash", func(t *testing.T) {
		url, err := client.BuildURL(paths.Path{Template: "/api/v1/games"}, nil)
		require.NoError(t, err)
		assert.Equal(t, "http://example.com/api/v1/games", url)
	})

	t.Run("base URL with path", func(t *testing.T) {
		client, err := New("http://example.com/api/v1")
		require.NoError(t, err)

		url, err := client.BuildURL(paths.GetGames, nil)
		require.NoError(t, err)
		assert.Equal(t, "http://example.com/api/v1/games", url)
	})

	t.Run("error from path.Build", func(t *testing.T) {
		url, err := client.BuildURL(paths.GetGame, paths.GameParams{Slug: ""})
		require.Error(t, err)
		assert.Empty(t, url)
	})
}

func TestClient_NewRESTRequestWithContext(t *testing.T) {
	client, err := New("http://example.com", WithHeader("X-Config-Header", "config-value"))
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("basic request", func(t *testing.T) {
		req, err := client.NewRESTRequestWithContext(ctx, http.MethodGet, "http://example.com/test", nil, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "http://example.com/test", req.URL.String())
		assert.Equal(t, "config-value", req.Header.Get("X-Config-Header"))
	})

	t.Run("request with headers", func(t *testing.T) {
		headers := map[string]string{
			"X-Custom":     "custom-value",
			"Content-Type": "application/json",
		}
		req, err := client.NewRESTRequestWithContext(ctx, http.MethodPost, "http://example.com/test", headers, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, "custom-value", req.Header.Get("X-Custom"))
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, "config-value", req.Header.Get("X-Config-Header")) // config headers should still be present
	})

	t.Run("request with query params", func(t *testing.T) {
		queryParams := map[string]string{
			"limit":  "10",
			"offset": "20",
		}
		req, err := client.NewRESTRequestWithContext(ctx, http.MethodGet, "http://example.com/test", nil, queryParams, nil)
		require.NoError(t, err)
		assert.Equal(t, "10", req.URL.Query().Get("limit"))
		assert.Equal(t, "20", req.URL.Query().Get("offset"))
	})

	t.Run("request with body", func(t *testing.T) {
		body := bytes.NewReader([]byte("test body"))
		req, err := client.NewRESTRequestWithContext(ctx, http.MethodPost, "http://example.com/test", nil, nil, body)
		require.NoError(t, err)
		assert.NotNil(t, req.Body)
	})

	t.Run("headers override config headers", func(t *testing.T) {
		headers := map[string]string{
			"X-Config-Header": "override-value",
		}
		req, err := client.NewRESTRequestWithContext(ctx, http.MethodGet, "http://example.com/test", headers, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, "override-value", req.Header.Get("X-Config-Header"))
	})

	t.Run("context is set", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		req, err := client.NewRESTRequestWithContext(ctx, http.MethodGet, "http://example.com/test", nil, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, ctx, req.Context())
	})
}

func TestClient_Do(t *testing.T) {
	t.Run("delegates to httpClient", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))
		defer server.Close()

		client, err := New(server.URL)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, server.URL+"/test", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})
}

func TestClient_RESTRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// RESTRequest is hardcoded to use paths.GetGames, so method and path are fixed
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/games", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"slug":"test"}]`))
		}))
		defer server.Close()

		client, err := New(server.URL)
		require.NoError(t, err)

		// Note: RESTRequest ignores the path and method parameters and uses paths.GetGames
		resp, err := client.RESTRequest(context.Background(), http.MethodGet, paths.GetGames, nil, nil, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("error status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client, err := New(server.URL)
		require.NoError(t, err)

		resp, err := client.RESTRequest(context.Background(), http.MethodGet, paths.GetGames, nil, nil, nil, nil)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "API response error")
	})

	t.Run("error building URL", func(t *testing.T) {
		client, err := New("http://example.com")
		require.NoError(t, err)

		// RESTRequest uses paths.GetGames but passes pathParams, which may cause error if invalid
		resp, err := client.RESTRequest(context.Background(), http.MethodGet, paths.GetGames, paths.GameParams{Slug: ""}, nil, nil, nil)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to build URL")
	})

	t.Run("error creating request", func(t *testing.T) {
		client, err := New("http://example.com")
		require.NoError(t, err)

		// Use invalid method to trigger error
		resp, err := client.RESTRequest(context.Background(), "INVALID METHOD\n", paths.GetGames, nil, nil, nil, nil)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestClient_Headers(t *testing.T) {
	client, err := New("http://example.com", WithHeader("X-Test", "value"))
	require.NoError(t, err)

	headers := client.Headers()
	assert.Equal(t, "value", headers["X-Test"])
}

func TestClient_Games(t *testing.T) {
	client, err := New("http://example.com")
	require.NoError(t, err)

	gamesClient := client.Games()
	assert.NotNil(t, gamesClient)
}

func TestClient_Users(t *testing.T) {
	client, err := New("http://example.com")
	require.NoError(t, err)

	usersClient := client.Users()
	assert.NotNil(t, usersClient)
}

func TestClient_HTTPClient(t *testing.T) {
	client, err := New("http://example.com")
	require.NoError(t, err)

	httpClient := client.HTTPClient()
	assert.NotNil(t, httpClient)
	assert.Equal(t, client.httpClient, httpClient)
}

func TestClient_BuildURL_ComplexPaths(t *testing.T) {
	client, err := New("http://example.com/api/v1")
	require.NoError(t, err)

	t.Run("user path", func(t *testing.T) {
		url, err := client.BuildURL(paths.PutUser, paths.UserParams{IPv4Address: "192.168.1.1"})
		require.NoError(t, err)
		assert.Equal(t, "http://example.com/api/v1/users/192.168.1.1", url)
	})

	t.Run("game icon path", func(t *testing.T) {
		url, err := client.BuildURL(paths.GetGameIcon, paths.GameParams{Slug: "test-game"})
		require.NoError(t, err)
		assert.Equal(t, "http://example.com/api/v1/games/test-game/icon", url)
	})
}

func TestClient_NewRESTRequestWithContext_InvalidContext(t *testing.T) {
	client, err := New("http://example.com")
	require.NoError(t, err)

	// Test with canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, err := client.NewRESTRequestWithContext(ctx, http.MethodGet, "http://example.com/test", nil, nil, nil)
	// Request creation should succeed even with canceled context
	require.NoError(t, err)
	assert.NotNil(t, req)
}

func TestClient_RESTRequest_WithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), "test-data")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New(server.URL)
	require.NoError(t, err)

	body := bytes.NewReader([]byte(`{"data":"test-data"}`))
	resp, err := client.RESTRequest(context.Background(), http.MethodPost, paths.GetGames, nil,
		map[string]string{"Content-Type": "application/json"}, nil, body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

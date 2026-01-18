package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seternate/go-lanty/pkg/api/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAPIResponse(t *testing.T) {
	t.Run("success status code", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)

		err := CheckAPIResponse(resp.Result())
		assert.NoError(t, err)
	})

	t.Run("error status code", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusNotFound)

		err := CheckAPIResponse(resp.Result())
		require.Error(t, err)
		assert.IsType(t, &internal.APIError{}, err)
	})

	t.Run("delegates to internal.CheckResponse", func(t *testing.T) {
		// Test that it calls internal.CheckResponse by testing with known responses
		t.Run("200 OK", func(t *testing.T) {
			resp := httptest.NewRecorder()
			resp.WriteHeader(http.StatusOK)
			err := CheckAPIResponse(resp.Result())
			assert.NoError(t, err)
		})

		t.Run("404 Not Found", func(t *testing.T) {
			resp := httptest.NewRecorder()
			resp.WriteHeader(http.StatusNotFound)
			err := CheckAPIResponse(resp.Result())
			require.Error(t, err)
			apiErr, ok := err.(*internal.APIError)
			require.True(t, ok)
			assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
		})

		t.Run("500 Internal Server Error", func(t *testing.T) {
			resp := httptest.NewRecorder()
			resp.WriteHeader(http.StatusInternalServerError)
			err := CheckAPIResponse(resp.Result())
			require.Error(t, err)
			apiErr, ok := err.(*internal.APIError)
			require.True(t, ok)
			assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		})
	})
}

func TestAPIError(t *testing.T) {
	t.Run("type alias works", func(t *testing.T) {
		// APIError is a type alias for internal.APIError
		var apiErr APIError
		assert.IsType(t, &internal.APIError{}, &apiErr)

		apiErr = APIError{
			StatusCode: http.StatusNotFound,
			Code:       "NOT_FOUND",
			Message:    "Resource not found",
		}

		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
		assert.Equal(t, "Resource not found", apiErr.Message)
	})
}

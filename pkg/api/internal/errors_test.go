package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrormodel "github.com/seternate/go-lanty/pkg/api/models/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIError_Error(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		err := &APIError{
			StatusCode: http.StatusNotFound,
			Code:       "NOT_FOUND",
			Message:    "Resource not found",
		}
		expected := "API error [404]: NOT_FOUND - Resource not found"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("without message", func(t *testing.T) {
		err := &APIError{
			StatusCode: http.StatusBadRequest,
			Code:       "BAD_REQUEST",
		}
		expected := "API error [400]: BAD_REQUEST"
		assert.Equal(t, expected, err.Error())
	})
}

func TestAPIError_IsNotFound(t *testing.T) {
	t.Run("404 status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusNotFound}
		assert.True(t, err.IsNotFound())
	})

	t.Run("other status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusBadRequest}
		assert.False(t, err.IsNotFound())
	})
}

func TestAPIError_IsBadRequest(t *testing.T) {
	t.Run("400 status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusBadRequest}
		assert.True(t, err.IsBadRequest())
	})

	t.Run("other status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusNotFound}
		assert.False(t, err.IsBadRequest())
	})
}

func TestAPIError_IsConflict(t *testing.T) {
	t.Run("409 status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusConflict}
		assert.True(t, err.IsConflict())
	})

	t.Run("other status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusNotFound}
		assert.False(t, err.IsConflict())
	})
}

func TestAPIError_IsUnauthorized(t *testing.T) {
	t.Run("401 status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusUnauthorized}
		assert.True(t, err.IsUnauthorized())
	})

	t.Run("other status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusNotFound}
		assert.False(t, err.IsUnauthorized())
	})
}

func TestAPIError_IsForbidden(t *testing.T) {
	t.Run("403 status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusForbidden}
		assert.True(t, err.IsForbidden())
	})

	t.Run("other status", func(t *testing.T) {
		err := &APIError{StatusCode: http.StatusNotFound}
		assert.False(t, err.IsForbidden())
	})
}

func TestParseErrorResponse(t *testing.T) {
	t.Run("valid JSON error response", func(t *testing.T) {
		errorResp := apierrormodel.Error{
			Error: apierrormodel.ErrorDescriptor{
				Code:    "CUSTOM_ERROR",
				Message: "Custom error message",
			},
		}
		body, _ := json.Marshal(errorResp)

		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusBadRequest)
		resp.Body.Write(body)
		httpResp := resp.Result()
		httpResp.Body = io.NopCloser(resp.Body)

		apiErr := ParseErrorResponse(httpResp)
		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		assert.Equal(t, "CUSTOM_ERROR", apiErr.Code)
		assert.Equal(t, "Custom error message", apiErr.Message)
		assert.NotNil(t, apiErr.Response)
	})

	t.Run("invalid JSON - fallback to status text", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusNotFound)
		resp.Body.WriteString("not json")
		httpResp := resp.Result()
		httpResp.Body = io.NopCloser(resp.Body)

		apiErr := ParseErrorResponse(httpResp)
		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
		assert.Equal(t, "Not Found", apiErr.Code)
		assert.Equal(t, "not json", apiErr.Message)
		assert.Nil(t, apiErr.Response)
	})

	t.Run("empty body", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusInternalServerError)
		httpResp := resp.Result()

		apiErr := ParseErrorResponse(httpResp)
		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		assert.Equal(t, "Internal Server Error", apiErr.Code)
		assert.Empty(t, apiErr.Message)
	})

	t.Run("error response with code only", func(t *testing.T) {
		errorResp := apierrormodel.Error{
			Error: apierrormodel.ErrorDescriptor{
				Code: "ERROR_CODE",
			},
		}
		body, _ := json.Marshal(errorResp)

		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusBadRequest)
		resp.Body.Write(body)
		httpResp := resp.Result()
		httpResp.Body = io.NopCloser(resp.Body)

		apiErr := ParseErrorResponse(httpResp)
		require.NotNil(t, apiErr)
		assert.Equal(t, "ERROR_CODE", apiErr.Code)
		assert.Empty(t, apiErr.Message)
	})

	t.Run("error response with message only", func(t *testing.T) {
		errorResp := apierrormodel.Error{
			Error: apierrormodel.ErrorDescriptor{
				Message: "Error message only",
			},
		}
		body, _ := json.Marshal(errorResp)

		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusBadRequest)
		resp.Body.Write(body)
		httpResp := resp.Result()
		httpResp.Body = io.NopCloser(resp.Body)

		apiErr := ParseErrorResponse(httpResp)
		require.NotNil(t, apiErr)
		assert.Equal(t, "Bad Request", apiErr.Code) // falls back to status text
		assert.Equal(t, "Error message only", apiErr.Message)
	})
}

func TestCheckResponse(t *testing.T) {
	t.Run("success status codes", func(t *testing.T) {
		statusCodes := []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusPartialContent,
		}

		for _, statusCode := range statusCodes {
			t.Run(http.StatusText(statusCode), func(t *testing.T) {
				resp := httptest.NewRecorder()
				resp.WriteHeader(statusCode)
				httpResp := resp.Result()

				err := CheckResponse(httpResp)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("error status code - calls ParseErrorResponse", func(t *testing.T) {
		errorResp := apierrormodel.Error{
			Error: apierrormodel.ErrorDescriptor{
				Code:    "NOT_FOUND",
				Message: "Resource not found",
			},
		}
		body, _ := json.Marshal(errorResp)

		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusNotFound)
		resp.Body.Write(body)
		httpResp := resp.Result()
		httpResp.Body = io.NopCloser(resp.Body)

		err := CheckResponse(httpResp)
		require.Error(t, err)
		apiErr, ok := err.(*APIError)
		require.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
		assert.Equal(t, "Resource not found", apiErr.Message)
	})

	t.Run("error status code - 500", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusInternalServerError)
		httpResp := resp.Result()

		err := CheckResponse(httpResp)
		require.Error(t, err)
		apiErr, ok := err.(*APIError)
		require.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	})
}

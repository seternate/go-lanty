package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apierrormodel "github.com/seternate/go-lanty/pkg/api/models/error"
)

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Response   *apierrormodel.Error
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error [%d]: %s - %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("API error [%d]: %s", e.StatusCode, e.Code)
}

func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

func (e *APIError) IsBadRequest() bool {
	return e.StatusCode == http.StatusBadRequest
}

func (e *APIError) IsConflict() bool {
	return e.StatusCode == http.StatusConflict
}

func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

func ParseErrorResponse(resp *http.Response) *APIError {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Code:       http.StatusText(resp.StatusCode),
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		apiErr.Message = "failed to read error response"
		return apiErr
	}

	var errorResp apierrormodel.Error
	if err := json.Unmarshal(body, &errorResp); err == nil {
		apiErr.Response = &errorResp
		if errorResp.Error.Code != "" {
			apiErr.Code = errorResp.Error.Code
		}
		if errorResp.Error.Message != "" {
			apiErr.Message = errorResp.Error.Message
		}
		return apiErr
	}

	if len(body) > 0 {
		apiErr.Message = string(body)
	}

	return apiErr
}

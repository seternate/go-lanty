package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponse(t *testing.T) {
	t.Run("JSON marshaling", func(t *testing.T) {
		resp := ErrorResponse{
			Error: ErrorDescriptor{
				Code:    "bad request",
				Message: "invalid input",
			},
		}

		data, err := json.Marshal(resp)
		require.NoError(t, err)

		var unmarshaled ErrorResponse
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, resp.Error.Code, unmarshaled.Error.Code)
		assert.Equal(t, resp.Error.Message, unmarshaled.Error.Message)
	})

	t.Run("JSON structure", func(t *testing.T) {
		resp := ErrorResponse{
			Error: ErrorDescriptor{
				Code:    "not found",
				Message: "resource not found",
			},
		}

		data, err := json.Marshal(resp)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Contains(t, jsonMap, "error")
		errorObj, ok := jsonMap["error"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "not found", errorObj["code"])
		assert.Equal(t, "resource not found", errorObj["message"])
	})
}

func TestErrorDescriptor(t *testing.T) {
	t.Run("JSON marshaling", func(t *testing.T) {
		desc := ErrorDescriptor{
			Code:    "conflict",
			Message: "resource already exists",
		}

		data, err := json.Marshal(desc)
		require.NoError(t, err)

		var unmarshaled ErrorDescriptor
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, desc.Code, unmarshaled.Code)
		assert.Equal(t, desc.Message, unmarshaled.Message)
	})

	t.Run("JSON field names", func(t *testing.T) {
		desc := ErrorDescriptor{
			Code:    "test",
			Message: "test message",
		}

		data, err := json.Marshal(desc)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Contains(t, jsonMap, "code")
		assert.Contains(t, jsonMap, "message")
		assert.Equal(t, "test", jsonMap["code"])
		assert.Equal(t, "test message", jsonMap["message"])
	})
}

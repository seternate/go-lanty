package gamemodel

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutable_JSON(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		requiresAdmin := true
		format := "exe"
		separator := " "
		now := time.Now()

		executable := Executable{
			Role:              "client",
			Path:              "/path/to/client.exe",
			RequiresAdmin:     &requiresAdmin,
			Format:            &format,
			ArgumentSeperator: &separator,
			Args:              []Arg{},
			CreatedAt:         now,
		}

		data, err := json.Marshal(executable)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "client", jsonMap["role"])
		assert.Equal(t, "/path/to/client.exe", jsonMap["path"])
		assert.Equal(t, true, jsonMap["requiresadmin"])
		assert.Equal(t, "exe", jsonMap["format"])
		assert.Equal(t, " ", jsonMap["argumentseperator"])
		assert.NotNil(t, jsonMap["args"])
	})

	t.Run("unmarshaling", func(t *testing.T) {
		jsonData := `{
			"role":"client",
			"path":"/path/to/client.exe",
			"requiresadmin":true,
			"format":"exe",
			"argumentseperator":" ",
			"args":[],
			"createdat":"2024-01-01T00:00:00Z"
		}`

		var executable Executable
		err := json.Unmarshal([]byte(jsonData), &executable)
		require.NoError(t, err)

		assert.Equal(t, "client", executable.Role)
		assert.Equal(t, "/path/to/client.exe", executable.Path)
		assert.NotNil(t, executable.RequiresAdmin)
		assert.True(t, *executable.RequiresAdmin)
		assert.NotNil(t, executable.Format)
		assert.Equal(t, "exe", *executable.Format)
		assert.NotNil(t, executable.Args)
		assert.Len(t, executable.Args, 0)
	})

	t.Run("round trip", func(t *testing.T) {
		requiresAdmin := false
		now := time.Now()

		original := Executable{
			Role:          "server",
			Path:          "/path/to/server",
			RequiresAdmin: &requiresAdmin,
			Args:          []Arg{},
			CreatedAt:     now,
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled Executable
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, original.Role, unmarshaled.Role)
		assert.Equal(t, original.Path, unmarshaled.Path)
		assert.NotNil(t, unmarshaled.RequiresAdmin)
		assert.False(t, *unmarshaled.RequiresAdmin)
		assert.Len(t, unmarshaled.Args, 0)
	})

	t.Run("with args", func(t *testing.T) {
		now := time.Now()
		required := true

		executable := Executable{
			Role:      "client",
			Path:      "/path/to/client",
			Args: []Arg{
				{
					Role:     "test",
					Required: &required,
					Argument: "test-arg",
					OrderIndex: 1,
					CreatedAt: now,
				},
			},
			CreatedAt: now,
		}

		data, err := json.Marshal(executable)
		require.NoError(t, err)

		var unmarshaled Executable
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.Args, 1)
		assert.Equal(t, "test", unmarshaled.Args[0].Role)
		assert.Equal(t, "test-arg", unmarshaled.Args[0].Argument)
	})

	t.Run("nil pointer fields", func(t *testing.T) {
		executable := Executable{
			Role:      "client",
			Path:      "/path/to/client",
			Args:      []Arg{},
			CreatedAt: time.Now(),
		}

		data, err := json.Marshal(executable)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "client", jsonMap["role"])
		// Nil pointers should be omitted or null
		_, hasRequiresAdmin := jsonMap["requiresadmin"]
		assert.False(t, hasRequiresAdmin)
	})
}

func TestUpsertGameExecutableRequest_JSON(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		requiresAdmin := true
		format := "exe"
		separator := " "

		req := UpsertGameExecutableRequest{
			Role:              "client",
			Path:              "/path/to/client.exe",
			RequiresAdmin:     &requiresAdmin,
			Format:            &format,
			ArgumentSeperator: &separator,
			Args:              []UpsertGameExecutableArgRequest{},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "client", jsonMap["role"])
		assert.Equal(t, "/path/to/client.exe", jsonMap["path"])
		assert.Equal(t, true, jsonMap["requiresadmin"])
		assert.Equal(t, "exe", jsonMap["format"])
		assert.Equal(t, " ", jsonMap["argumentseperator"])
		assert.NotNil(t, jsonMap["args"])
	})

	t.Run("unmarshaling", func(t *testing.T) {
		jsonData := `{
			"role":"client",
			"path":"/path/to/client.exe",
			"requiresadmin":true,
			"format":"exe",
			"argumentseperator":" ",
			"args":[]
		}`

		var req UpsertGameExecutableRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, "client", req.Role)
		assert.Equal(t, "/path/to/client.exe", req.Path)
		assert.NotNil(t, req.RequiresAdmin)
		assert.True(t, *req.RequiresAdmin)
		assert.NotNil(t, req.Args)
		assert.Len(t, req.Args, 0)
	})

	t.Run("round trip", func(t *testing.T) {
		requiresAdmin := false

		original := UpsertGameExecutableRequest{
			Role:          "server",
			Path:          "/path/to/server",
			RequiresAdmin: &requiresAdmin,
			Args:          []UpsertGameExecutableArgRequest{},
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled UpsertGameExecutableRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, original.Role, unmarshaled.Role)
		assert.Equal(t, original.Path, unmarshaled.Path)
		assert.NotNil(t, unmarshaled.RequiresAdmin)
		assert.False(t, *unmarshaled.RequiresAdmin)
		assert.Len(t, unmarshaled.Args, 0)
	})

	t.Run("with args", func(t *testing.T) {
		required := true

		req := UpsertGameExecutableRequest{
			Role: "client",
			Path: "/path/to/client",
			Args: []UpsertGameExecutableArgRequest{
				{
					Role:     "test",
					Name:     "test-name",
					Required: &required,
					Argument: "test-arg",
					OrderIndex: 1,
				},
			},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled UpsertGameExecutableRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.Args, 1)
		assert.Equal(t, "test", unmarshaled.Args[0].Role)
		assert.Equal(t, "test-name", unmarshaled.Args[0].Name)
		assert.Equal(t, "test-arg", unmarshaled.Args[0].Argument)
	})
}

package gamemodel

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArg_JSON(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		required := true
		enabled := false
		format := "string"
		arg := Arg{
			Role:              "test",
			Required:          &required,
			Enabled:           &enabled,
			Format:            &format,
			Argument:          "test-arg",
			OrderIndex:        1,
			CreatedAt:         time.Now(),
		}

		data, err := json.Marshal(arg)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "test", jsonMap["role"])
		assert.Equal(t, true, jsonMap["required"])
		assert.Equal(t, false, jsonMap["enabled"])
		assert.Equal(t, "string", jsonMap["format"])
		assert.Equal(t, "test-arg", jsonMap["argument"])
		assert.Equal(t, float64(1), jsonMap["orderindex"])
	})

	t.Run("unmarshaling", func(t *testing.T) {
		jsonData := `{
			"role":"test",
			"required":true,
			"enabled":false,
			"format":"string",
			"argument":"test-arg",
			"orderindex":1,
			"createdat":"2024-01-01T00:00:00Z"
		}`

		var arg Arg
		err := json.Unmarshal([]byte(jsonData), &arg)
		require.NoError(t, err)

		assert.Equal(t, "test", arg.Role)
		assert.NotNil(t, arg.Required)
		assert.True(t, *arg.Required)
		assert.NotNil(t, arg.Enabled)
		assert.False(t, *arg.Enabled)
		assert.Equal(t, "test-arg", arg.Argument)
		assert.Equal(t, int64(1), arg.OrderIndex)
	})

	t.Run("nil pointer fields", func(t *testing.T) {
		arg := Arg{
			Role:      "test",
			Argument:  "test-arg",
			OrderIndex: 1,
			CreatedAt: time.Now(),
		}

		data, err := json.Marshal(arg)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "test", jsonMap["role"])
		// Nil pointers should be omitted from JSON
		_, hasRequired := jsonMap["required"]
		assert.False(t, hasRequired)
	})

	t.Run("all pointer fields", func(t *testing.T) {
		required := true
		enabled := true
		format := "int"
		separator := ","
		description := "Test description"
		defaultString := "default"
		defaultBool := true
		defaultInt := int64(42)
		defaultFloat := 3.14
		minInt := int64(0)
		maxInt := int64(100)
		minFloat := 0.0
		maxFloat := 100.0
		precision := int64(2)

		arg := Arg{
			Role:              "test",
			Required:          &required,
			Enabled:           &enabled,
			Format:            &format,
			ArgumentSeparator: &separator,
			Argument:          "test-arg",
			Description:       &description,
			DefaultString:     &defaultString,
			DefaultBool:       &defaultBool,
			DefaultInt:        &defaultInt,
			DefaultFloat:      &defaultFloat,
			EnumValues:        []string{"value1", "value2"},
			MinInt:            &minInt,
			MaxInt:            &maxInt,
			MinFloat:          &minFloat,
			MaxFloat:          &maxFloat,
			FloatPrecision:    &precision,
			OrderIndex:        1,
			CreatedAt:         time.Now(),
		}

		data, err := json.Marshal(arg)
		require.NoError(t, err)

		var unmarshaled Arg
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, arg.Role, unmarshaled.Role)
		assert.NotNil(t, unmarshaled.Required)
		assert.True(t, *unmarshaled.Required)
		assert.NotNil(t, unmarshaled.DefaultInt)
		assert.Equal(t, int64(42), *unmarshaled.DefaultInt)
		assert.Len(t, unmarshaled.EnumValues, 2)
	})
}

func TestUpsertGameExecutableArgRequest_JSON(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		required := true
		enabled := false
		req := UpsertGameExecutableArgRequest{
			Role:     "test",
			Name:     "test-name",
			Required: &required,
			Enabled:  &enabled,
			Argument: "test-arg",
			OrderIndex: 1,
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "test", jsonMap["role"])
		assert.Equal(t, "test-name", jsonMap["name"])
		assert.Equal(t, true, jsonMap["required"])
		assert.Equal(t, false, jsonMap["enabled"])
		assert.Equal(t, "test-arg", jsonMap["argument"])
	})

	t.Run("unmarshaling", func(t *testing.T) {
		jsonData := `{
			"role":"test",
			"name":"test-name",
			"required":true,
			"enabled":false,
			"argument":"test-arg",
			"orderindex":1
		}`

		var req UpsertGameExecutableArgRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, "test", req.Role)
		assert.Equal(t, "test-name", req.Name)
		assert.NotNil(t, req.Required)
		assert.True(t, *req.Required)
		assert.Equal(t, "test-arg", req.Argument)
		assert.Equal(t, int64(1), req.OrderIndex)
	})

	t.Run("round trip", func(t *testing.T) {
		required := true
		original := UpsertGameExecutableArgRequest{
			Role:      "test",
			Name:      "test-name",
			Required:  &required,
			Argument:  "test-arg",
			OrderIndex: 1,
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled UpsertGameExecutableArgRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, original.Role, unmarshaled.Role)
		assert.Equal(t, original.Name, unmarshaled.Name)
		assert.NotNil(t, unmarshaled.Required)
		assert.True(t, *unmarshaled.Required)
		assert.Equal(t, original.Argument, unmarshaled.Argument)
	})
}

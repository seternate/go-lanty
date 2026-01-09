package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appGameSrv "github.com/seternate/go-lanty/pkg/application/game"
)

func TestNewGameExecutableArgResponse(t *testing.T) {
	t.Run("all fields populated", func(t *testing.T) {
		required := true
		enabled := false
		format := "string"
		argSeparator := "="
		description := "test description"
		defaultString := "default"
		defaultBool := true
		defaultInt := int64(42)
		defaultFloat := 3.14
		enumValues := []string{"val1", "val2"}
		minInt := int64(0)
		maxInt := int64(100)
		minFloat := 0.0
		maxFloat := 100.0
		floatPrecision := int64(2)
		orderIndex := int64(1)
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		argView := appGameSrv.GameArgView{
			Role:           "test-role",
			Name:           "test-name",
			Required:       &required,
			Enabled:        &enabled,
			Format:         &format,
			ArgSeparator:   &argSeparator,
			Arg:            "--test",
			Description:    &description,
			DefaultString:  &defaultString,
			DefaultBool:    &defaultBool,
			DefaultInt:     &defaultInt,
			DefaultFloat:   &defaultFloat,
			Enums:          enumValues,
			MinInt:         &minInt,
			MaxInt:         &maxInt,
			MinFloat:       &minFloat,
			MaxFloat:       &maxFloat,
			FloatPrecision: &floatPrecision,
			OrderIndex:     orderIndex,
			CreatedAt:      createdAt,
		}

		result := NewGameExecutableArgResponse(argView)

		require.NotNil(t, result)
		assert.Equal(t, "test-role", result.Role)
		assert.Equal(t, &required, result.Required)
		assert.Equal(t, &enabled, result.Enabled)
		assert.Equal(t, &format, result.Format)
		assert.Equal(t, &argSeparator, result.ArgumentSeparator)
		assert.Equal(t, "--test", result.Argument)
		assert.Equal(t, &description, result.Description)
		assert.Equal(t, &defaultString, result.DefaultString)
		assert.Equal(t, &defaultBool, result.DefaultBool)
		assert.Equal(t, &defaultInt, result.DefaultInt)
		assert.Equal(t, &defaultFloat, result.DefaultFloat)
		assert.Equal(t, enumValues, result.EnumValues)
		assert.Equal(t, &minInt, result.MinInt)
		assert.Equal(t, &maxInt, result.MaxInt)
		assert.Equal(t, &minFloat, result.MinFloat)
		assert.Equal(t, &maxFloat, result.MaxFloat)
		assert.Equal(t, &floatPrecision, result.FloatPrecision)
		assert.Equal(t, orderIndex, result.OrderIndex)
		assert.Equal(t, createdAt, result.CreatedAt)
	})

	t.Run("minimal fields", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		argView := appGameSrv.GameArgView{
			Role:       "test-role",
			Name:       "test-name",
			Arg:        "--test",
			OrderIndex: 0,
			CreatedAt:  createdAt,
		}

		result := NewGameExecutableArgResponse(argView)

		require.NotNil(t, result)
		assert.Equal(t, "test-role", result.Role)
		assert.Nil(t, result.Required)
		assert.Nil(t, result.Enabled)
		assert.Nil(t, result.Format)
		assert.Nil(t, result.ArgumentSeparator)
		assert.Equal(t, "--test", result.Argument)
		assert.Nil(t, result.Description)
		assert.Nil(t, result.DefaultString)
		assert.Nil(t, result.DefaultBool)
		assert.Nil(t, result.DefaultInt)
		assert.Nil(t, result.DefaultFloat)
		assert.Nil(t, result.EnumValues)
		assert.Nil(t, result.MinInt)
		assert.Nil(t, result.MaxInt)
		assert.Nil(t, result.MinFloat)
		assert.Nil(t, result.MaxFloat)
		assert.Nil(t, result.FloatPrecision)
		assert.Equal(t, int64(0), result.OrderIndex)
		assert.Equal(t, createdAt, result.CreatedAt)
	})

	t.Run("nil pointer fields", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		argView := appGameSrv.GameArgView{
			Role:           "test-role",
			Name:           "test-name",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--test",
			Description:    nil,
			DefaultString:  nil,
			DefaultBool:    nil,
			DefaultInt:     nil,
			DefaultFloat:   nil,
			Enums:          nil,
			MinInt:         nil,
			MaxInt:         nil,
			MinFloat:       nil,
			MaxFloat:       nil,
			FloatPrecision: nil,
			OrderIndex:     0,
			CreatedAt:      createdAt,
		}

		result := NewGameExecutableArgResponse(argView)

		require.NotNil(t, result)
		assert.Nil(t, result.Required)
		assert.Nil(t, result.Enabled)
		assert.Nil(t, result.Format)
		assert.Nil(t, result.ArgumentSeparator)
		assert.Nil(t, result.Description)
		assert.Nil(t, result.DefaultString)
		assert.Nil(t, result.DefaultBool)
		assert.Nil(t, result.DefaultInt)
		assert.Nil(t, result.DefaultFloat)
		assert.Nil(t, result.EnumValues)
		assert.Nil(t, result.MinInt)
		assert.Nil(t, result.MaxInt)
		assert.Nil(t, result.MinFloat)
		assert.Nil(t, result.MaxFloat)
		assert.Nil(t, result.FloatPrecision)
	})

	t.Run("empty enum values", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		argView := appGameSrv.GameArgView{
			Role:       "test-role",
			Name:       "test-name",
			Arg:        "--test",
			Enums:      []string{},
			OrderIndex: 0,
			CreatedAt:  createdAt,
		}

		result := NewGameExecutableArgResponse(argView)

		require.NotNil(t, result)
		assert.NotNil(t, result.EnumValues)
		assert.Len(t, result.EnumValues, 0)
	})
}

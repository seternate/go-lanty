package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appGameSrv "github.com/seternate/go-lanty/pkg/application/game"
)

func TestNewGameExecutableResponse(t *testing.T) {
	t.Run("all fields populated with args", func(t *testing.T) {
		requiresAdmin := true
		format := "exe"
		argSeparator := " "
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		required := true
		argView := appGameSrv.GameArgView{
			Role:       "arg-role",
			Name:       "arg-name",
			Required:   &required,
			Arg:        "--test",
			OrderIndex: 0,
			CreatedAt:  createdAt,
		}

		execView := appGameSrv.GameExecView{
			GameSlug:      "test-game",
			Role:          "launcher",
			Path:          "/path/to/executable",
			RequiresAdmin: &requiresAdmin,
			Format:        &format,
			ArgSeperator:  &argSeparator,
			CreatedAt:     createdAt,
			Args:          []appGameSrv.GameArgView{argView},
		}

		result := NewGameExecutableResponse(execView)

		require.NotNil(t, result)
		assert.Equal(t, "launcher", result.Role)
		assert.Equal(t, "/path/to/executable", result.Path)
		assert.Equal(t, &requiresAdmin, result.RequiresAdmin)
		assert.Equal(t, &format, result.Format)
		assert.Equal(t, &argSeparator, result.ArgumentSeperator)
		assert.Equal(t, createdAt, result.CreatedAt)
		require.Len(t, result.Args, 1)
		assert.Equal(t, "arg-role", result.Args[0].Role)
		assert.Equal(t, "--test", result.Args[0].Argument)
	})

	t.Run("minimal fields no args", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		execView := appGameSrv.GameExecView{
			GameSlug:  "test-game",
			Role:      "launcher",
			Path:      "/path/to/executable",
			CreatedAt: createdAt,
			Args:      []appGameSrv.GameArgView{},
		}

		result := NewGameExecutableResponse(execView)

		require.NotNil(t, result)
		assert.Equal(t, "launcher", result.Role)
		assert.Equal(t, "/path/to/executable", result.Path)
		assert.Nil(t, result.RequiresAdmin)
		assert.Nil(t, result.Format)
		assert.Nil(t, result.ArgumentSeperator)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.NotNil(t, result.Args)
		assert.Len(t, result.Args, 0)
	})

	t.Run("nil pointer fields", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		execView := appGameSrv.GameExecView{
			GameSlug:      "test-game",
			Role:          "launcher",
			Path:          "/path/to/executable",
			RequiresAdmin: nil,
			Format:        nil,
			ArgSeperator:  nil,
			CreatedAt:     createdAt,
			Args:          nil,
		}

		result := NewGameExecutableResponse(execView)

		require.NotNil(t, result)
		assert.Nil(t, result.RequiresAdmin)
		assert.Nil(t, result.Format)
		assert.Nil(t, result.ArgumentSeperator)
		assert.NotNil(t, result.Args)
		assert.Len(t, result.Args, 0)
	})

	t.Run("multiple args", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		arg1 := appGameSrv.GameArgView{
			Role:       "arg1",
			Name:       "name1",
			Arg:        "--arg1",
			OrderIndex: 0,
			CreatedAt:  createdAt,
		}

		arg2 := appGameSrv.GameArgView{
			Role:       "arg2",
			Name:       "name2",
			Arg:        "--arg2",
			OrderIndex: 1,
			CreatedAt:  createdAt,
		}

		execView := appGameSrv.GameExecView{
			GameSlug:  "test-game",
			Role:      "launcher",
			Path:      "/path/to/executable",
			CreatedAt: createdAt,
			Args:      []appGameSrv.GameArgView{arg1, arg2},
		}

		result := NewGameExecutableResponse(execView)

		require.NotNil(t, result)
		require.Len(t, result.Args, 2)
		assert.Equal(t, "arg1", result.Args[0].Role)
		assert.Equal(t, "--arg1", result.Args[0].Argument)
		assert.Equal(t, "arg2", result.Args[1].Role)
		assert.Equal(t, "--arg2", result.Args[1].Argument)
	})
}

package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appGameSrv "github.com/seternate/go-lanty/pkg/application/game"
)

func TestNewGameResponse(t *testing.T) {
	t.Run("all fields populated with executables", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		exec1 := appGameSrv.GameExecView{
			GameSlug:  "test-game",
			Role:      "launcher",
			Path:      "/path/to/launcher",
			CreatedAt: createdAt,
			Args:      []appGameSrv.GameArgView{},
		}

		exec2 := appGameSrv.GameExecView{
			GameSlug:  "test-game",
			Role:      "uninstaller",
			Path:      "/path/to/uninstaller",
			CreatedAt: createdAt,
			Args:      []appGameSrv.GameArgView{},
		}

		gameView := &appGameSrv.GameView{
			Slug:      "test-game",
			Name:      "Test Game",
			CreatedAt: createdAt,
			Execs: map[string]appGameSrv.GameExecView{
				"launcher":    exec1,
				"uninstaller": exec2,
			},
		}

		result := NewGameResponse(gameView)

		require.NotNil(t, result)
		assert.Equal(t, "test-game", result.Slug)
		assert.Equal(t, "Test Game", result.Name)
		assert.Equal(t, createdAt, result.CreatedAt)
		require.Len(t, result.Executables, 2)

		// Check that both executables are present (order may vary)
		roles := make(map[string]bool)
		for _, exec := range result.Executables {
			roles[exec.Role] = true
		}
		assert.True(t, roles["launcher"])
		assert.True(t, roles["uninstaller"])
	})

	t.Run("minimal fields no executables", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		gameView := &appGameSrv.GameView{
			Slug:      "test-game",
			Name:      "Test Game",
			CreatedAt: createdAt,
			Execs:     map[string]appGameSrv.GameExecView{},
		}

		result := NewGameResponse(gameView)

		require.NotNil(t, result)
		assert.Equal(t, "test-game", result.Slug)
		assert.Equal(t, "Test Game", result.Name)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.NotNil(t, result.Executables)
		assert.Len(t, result.Executables, 0)
	})

	t.Run("executable with args", func(t *testing.T) {
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

		exec := appGameSrv.GameExecView{
			GameSlug:  "test-game",
			Role:      "launcher",
			Path:      "/path/to/launcher",
			CreatedAt: createdAt,
			Args:      []appGameSrv.GameArgView{argView},
		}

		gameView := &appGameSrv.GameView{
			Slug:      "test-game",
			Name:      "Test Game",
			CreatedAt: createdAt,
			Execs: map[string]appGameSrv.GameExecView{
				"launcher": exec,
			},
		}

		result := NewGameResponse(gameView)

		require.NotNil(t, result)
		require.Len(t, result.Executables, 1)
		assert.Equal(t, "launcher", result.Executables[0].Role)
		require.Len(t, result.Executables[0].Args, 1)
		assert.Equal(t, "arg-role", result.Executables[0].Args[0].Role)
		assert.Equal(t, "--test", result.Executables[0].Args[0].Argument)
	})
}

func TestUpsertGameRequest_ToCommand(t *testing.T) {
	t.Run("all fields populated", func(t *testing.T) {
		requiresAdmin := true
		format := "exe"
		argSeparator := " "
		required := true
		enabled := false
		defaultString := "default"
		defaultBool := true
		defaultInt := int64(42)
		defaultFloat := 3.14
		minInt := int64(0)
		maxInt := int64(100)
		minFloat := 0.0
		maxFloat := 100.0
		floatPrecision := int64(2)

		req := &UpsertGameRequest{
			Name: "Test Game",
			Executables: []UpsertGameExecutableRequest{
				{
					Role:              "launcher",
					Path:              "/path/to/launcher",
					RequiresAdmin:     &requiresAdmin,
					Format:            &format,
					ArgumentSeperator: &argSeparator,
					Args: []UpsertGameExecutableArgRequest{
						{
							Role:              "arg-role",
							Name:              "arg-name",
							Required:          &required,
							Enabled:           &enabled,
							Format:            &format,
							ArgumentSeparator: &argSeparator,
							Argument:          "--test",
							DefaultString:     &defaultString,
							DefaultBool:       &defaultBool,
							DefaultInt:        &defaultInt,
							DefaultFloat:      &defaultFloat,
							EnumValues:        []string{"val1", "val2"},
							MinInt:            &minInt,
							MaxInt:            &maxInt,
							MinFloat:          &minFloat,
							MaxFloat:          &maxFloat,
							FloatPrecision:    &floatPrecision,
							OrderIndex:        0,
						},
					},
				},
			},
		}

		slug := "test-game"
		cmd := req.ToCommand(slug)

		assert.Equal(t, slug, cmd.Slug)
		assert.Equal(t, "Test Game", cmd.Name)
		require.Len(t, cmd.Executables, 1)
		assert.Equal(t, "launcher", cmd.Executables[0].Role)
		assert.Equal(t, "/path/to/launcher", cmd.Executables[0].Path)
		assert.Equal(t, &requiresAdmin, cmd.Executables[0].RequiresAdmin)
		assert.Equal(t, &format, cmd.Executables[0].Format)
		assert.Equal(t, &argSeparator, cmd.Executables[0].ArgumentSeperator)
		require.Len(t, cmd.Executables[0].Args, 1)
		assert.Equal(t, "arg-role", cmd.Executables[0].Args[0].Role)
		assert.Equal(t, "arg-name", cmd.Executables[0].Args[0].Name)
		assert.Equal(t, &required, cmd.Executables[0].Args[0].Required)
		assert.Equal(t, &enabled, cmd.Executables[0].Args[0].Enabled)
		assert.Equal(t, &format, cmd.Executables[0].Args[0].Format)
		assert.Equal(t, &argSeparator, cmd.Executables[0].Args[0].ArgumentSeparator)
		assert.Equal(t, "--test", cmd.Executables[0].Args[0].Argument)
		assert.Equal(t, &defaultString, cmd.Executables[0].Args[0].DefaultString)
		assert.Equal(t, &defaultBool, cmd.Executables[0].Args[0].DefaultBool)
		assert.Equal(t, &defaultInt, cmd.Executables[0].Args[0].DefaultInt)
		assert.Equal(t, &defaultFloat, cmd.Executables[0].Args[0].DefaultFloat)
		assert.Equal(t, []string{"val1", "val2"}, cmd.Executables[0].Args[0].EnumValues)
		assert.Equal(t, &minInt, cmd.Executables[0].Args[0].MinInt)
		assert.Equal(t, &maxInt, cmd.Executables[0].Args[0].MaxInt)
		assert.Equal(t, &minFloat, cmd.Executables[0].Args[0].MinFloat)
		assert.Equal(t, &maxFloat, cmd.Executables[0].Args[0].MaxFloat)
		assert.Equal(t, &floatPrecision, cmd.Executables[0].Args[0].FloatPrecision)
		assert.Equal(t, int64(0), cmd.Executables[0].Args[0].OrderIndex)
	})

	t.Run("minimal fields", func(t *testing.T) {
		req := &UpsertGameRequest{
			Name:        "Test Game",
			Executables: []UpsertGameExecutableRequest{},
		}

		slug := "test-game"
		cmd := req.ToCommand(slug)

		assert.Equal(t, slug, cmd.Slug)
		assert.Equal(t, "Test Game", cmd.Name)
		assert.NotNil(t, cmd.Executables)
		assert.Len(t, cmd.Executables, 0)
	})

	t.Run("multiple executables", func(t *testing.T) {
		req := &UpsertGameRequest{
			Name: "Test Game",
			Executables: []UpsertGameExecutableRequest{
				{
					Role: "launcher",
					Path: "/path/to/launcher",
					Args: []UpsertGameExecutableArgRequest{},
				},
				{
					Role: "uninstaller",
					Path: "/path/to/uninstaller",
					Args: []UpsertGameExecutableArgRequest{},
				},
			},
		}

		slug := "test-game"
		cmd := req.ToCommand(slug)

		assert.Equal(t, slug, cmd.Slug)
		assert.Equal(t, "Test Game", cmd.Name)
		require.Len(t, cmd.Executables, 2)
		assert.Equal(t, "launcher", cmd.Executables[0].Role)
		assert.Equal(t, "/path/to/launcher", cmd.Executables[0].Path)
		assert.Equal(t, "uninstaller", cmd.Executables[1].Role)
		assert.Equal(t, "/path/to/uninstaller", cmd.Executables[1].Path)
	})

	t.Run("executable with multiple args", func(t *testing.T) {
		req := &UpsertGameRequest{
			Name: "Test Game",
			Executables: []UpsertGameExecutableRequest{
				{
					Role: "launcher",
					Path: "/path/to/launcher",
					Args: []UpsertGameExecutableArgRequest{
						{
							Role:     "arg1",
							Name:     "name1",
							Argument: "--arg1",
							OrderIndex: 0,
						},
						{
							Role:     "arg2",
							Name:     "name2",
							Argument: "--arg2",
							OrderIndex: 1,
						},
					},
				},
			},
		}

		slug := "test-game"
		cmd := req.ToCommand(slug)

		require.Len(t, cmd.Executables, 1)
		require.Len(t, cmd.Executables[0].Args, 2)
		assert.Equal(t, "arg1", cmd.Executables[0].Args[0].Role)
		assert.Equal(t, "--arg1", cmd.Executables[0].Args[0].Argument)
		assert.Equal(t, "arg2", cmd.Executables[0].Args[1].Role)
		assert.Equal(t, "--arg2", cmd.Executables[0].Args[1].Argument)
	})

	t.Run("nil pointer fields", func(t *testing.T) {
		req := &UpsertGameRequest{
			Name: "Test Game",
			Executables: []UpsertGameExecutableRequest{
				{
					Role:              "launcher",
					Path:              "/path/to/launcher",
					RequiresAdmin:     nil,
					Format:            nil,
					ArgumentSeperator: nil,
					Args: []UpsertGameExecutableArgRequest{
						{
							Role:              "arg-role",
							Name:              "arg-name",
							Required:          nil,
							Enabled:           nil,
							Format:            nil,
							ArgumentSeparator: nil,
							Argument:          "--test",
							DefaultString:     nil,
							DefaultBool:       nil,
							DefaultInt:        nil,
							DefaultFloat:      nil,
							EnumValues:        nil,
							MinInt:            nil,
							MaxInt:            nil,
							MinFloat:          nil,
							MaxFloat:          nil,
							FloatPrecision:    nil,
							OrderIndex:        0,
						},
					},
				},
			},
		}

		slug := "test-game"
		cmd := req.ToCommand(slug)

		assert.Equal(t, slug, cmd.Slug)
		require.Len(t, cmd.Executables, 1)
		assert.Nil(t, cmd.Executables[0].RequiresAdmin)
		assert.Nil(t, cmd.Executables[0].Format)
		assert.Nil(t, cmd.Executables[0].ArgumentSeperator)
		require.Len(t, cmd.Executables[0].Args, 1)
		assert.Nil(t, cmd.Executables[0].Args[0].Required)
		assert.Nil(t, cmd.Executables[0].Args[0].Enabled)
		assert.Nil(t, cmd.Executables[0].Args[0].Format)
		assert.Nil(t, cmd.Executables[0].Args[0].ArgumentSeparator)
		assert.Nil(t, cmd.Executables[0].Args[0].DefaultString)
		assert.Nil(t, cmd.Executables[0].Args[0].DefaultBool)
		assert.Nil(t, cmd.Executables[0].Args[0].DefaultInt)
		assert.Nil(t, cmd.Executables[0].Args[0].DefaultFloat)
		assert.Nil(t, cmd.Executables[0].Args[0].EnumValues)
		assert.Nil(t, cmd.Executables[0].Args[0].MinInt)
		assert.Nil(t, cmd.Executables[0].Args[0].MaxInt)
		assert.Nil(t, cmd.Executables[0].Args[0].MinFloat)
		assert.Nil(t, cmd.Executables[0].Args[0].MaxFloat)
		assert.Nil(t, cmd.Executables[0].Args[0].FloatPrecision)
	})
}

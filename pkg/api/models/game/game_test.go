package gamemodel

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGame_JSON(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		now := time.Now()
		game := Game{
			Slug:        "test-game",
			Name:        "Test Game",
			Executables: []Executable{},
			CreatedAt:   now,
		}

		data, err := json.Marshal(game)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "test-game", jsonMap["slug"])
		assert.Equal(t, "Test Game", jsonMap["name"])
		assert.NotNil(t, jsonMap["executables"])
	})

	t.Run("unmarshaling", func(t *testing.T) {
		jsonData := `{"slug":"test-game","name":"Test Game","executables":[],"createdat":"2024-01-01T00:00:00Z"}`

		var game Game
		err := json.Unmarshal([]byte(jsonData), &game)
		require.NoError(t, err)

		assert.Equal(t, "test-game", game.Slug)
		assert.Equal(t, "Test Game", game.Name)
		assert.NotNil(t, game.Executables)
	})

	t.Run("round trip", func(t *testing.T) {
		now := time.Now()
		original := Game{
			Slug:        "roundtrip",
			Name:        "Round Trip Game",
			Executables: []Executable{},
			CreatedAt:   now,
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled Game
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, original.Slug, unmarshaled.Slug)
		assert.Equal(t, original.Name, unmarshaled.Name)
		assert.Len(t, unmarshaled.Executables, 0)
	})

	t.Run("with executables", func(t *testing.T) {
		now := time.Now()
		game := Game{
			Slug: "test-game",
			Name: "Test Game",
			Executables: []Executable{
				{
					Role:      "client",
					Path:      "/path/to/client",
					CreatedAt: now,
				},
			},
			CreatedAt: now,
		}

		data, err := json.Marshal(game)
		require.NoError(t, err)

		var unmarshaled Game
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.Executables, 1)
		assert.Equal(t, "client", unmarshaled.Executables[0].Role)
		assert.Equal(t, "/path/to/client", unmarshaled.Executables[0].Path)
	})
}

func TestUpsertGameRequest_JSON(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		req := UpsertGameRequest{
			Name:        "New Game",
			Executables: []UpsertGameExecutableRequest{},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "New Game", jsonMap["name"])
		assert.NotNil(t, jsonMap["executables"])
	})

	t.Run("unmarshaling", func(t *testing.T) {
		jsonData := `{"name":"New Game","executables":[]}`

		var req UpsertGameRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, "New Game", req.Name)
		assert.NotNil(t, req.Executables)
	})

	t.Run("round trip", func(t *testing.T) {
		original := UpsertGameRequest{
			Name:        "Round Trip",
			Executables: []UpsertGameExecutableRequest{},
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled UpsertGameRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, original.Name, unmarshaled.Name)
		assert.Len(t, unmarshaled.Executables, 0)
	})

	t.Run("with executables", func(t *testing.T) {
		req := UpsertGameRequest{
			Name: "New Game",
			Executables: []UpsertGameExecutableRequest{
				{
					Role: "client",
					Path: "/path/to/client",
				},
			},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled UpsertGameRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.Executables, 1)
		assert.Equal(t, "client", unmarshaled.Executables[0].Role)
		assert.Equal(t, "/path/to/client", unmarshaled.Executables[0].Path)
	})
}

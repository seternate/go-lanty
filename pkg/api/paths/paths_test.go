package paths

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPath_Build(t *testing.T) {
	t.Run("with nil params", func(t *testing.T) {
		path := Path{Template: "/games"}
		result, err := path.Build(nil)
		require.NoError(t, err)
		assert.Equal(t, "/games", result)
	})

	t.Run("with params", func(t *testing.T) {
		path := Path{Template: "/games/:slug"}
		params := GameParams{Slug: "test-game"}
		result, err := path.Build(params)
		require.NoError(t, err)
		assert.Equal(t, "/games/test-game", result)
	})

	t.Run("with params that error", func(t *testing.T) {
		path := Path{Template: "/games/:slug"}
		params := GameParams{Slug: ""} // Empty slug should error
		result, err := path.Build(params)
		require.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("path without params", func(t *testing.T) {
		path := Path{Template: "/games"}
		params := GameParams{Slug: "test"} // Params provided but not needed
		result, err := path.Build(params)
		require.NoError(t, err)
		// Should replace slug even though not in template
		assert.Equal(t, "/games", result)
	})
}

func TestReplaceParams(t *testing.T) {
	t.Run("single param", func(t *testing.T) {
		result := replaceParams("/games/:slug", "slug", "test-game")
		assert.Equal(t, "/games/test-game", result)
	})

	t.Run("multiple occurrences", func(t *testing.T) {
		result := replaceParams("/games/:slug/:slug", "slug", "test-game")
		assert.Equal(t, "/games/test-game/test-game", result)
	})

	t.Run("no param in path", func(t *testing.T) {
		result := replaceParams("/games", "slug", "test-game")
		assert.Equal(t, "/games", result)
	})

	t.Run("different param", func(t *testing.T) {
		result := replaceParams("/games/:slug", "id", "123")
		assert.Equal(t, "/games/:slug", result)
	})

	t.Run("empty value", func(t *testing.T) {
		result := replaceParams("/games/:slug", "slug", "")
		assert.Equal(t, "/games/", result)
	})
}

func TestPath_GetGames(t *testing.T) {
	result, err := GetGames.Build(nil)
	require.NoError(t, err)
	assert.Equal(t, "/games", result)
	assert.Equal(t, "GET", GetGames.Method)
}

func TestPath_GetGame(t *testing.T) {
	result, err := GetGame.Build(GameParams{Slug: "test"})
	require.NoError(t, err)
	assert.Equal(t, "/games/test", result)
	assert.Equal(t, "GET", GetGame.Method)
}

func TestPath_PutGame(t *testing.T) {
	result, err := PutGame.Build(GameParams{Slug: "test"})
	require.NoError(t, err)
	assert.Equal(t, "/games/test", result)
	assert.Equal(t, "PUT", PutGame.Method)
}

func TestPath_DeleteGame(t *testing.T) {
	result, err := DeleteGame.Build(GameParams{Slug: "test"})
	require.NoError(t, err)
	assert.Equal(t, "/games/test", result)
	assert.Equal(t, "DELETE", DeleteGame.Method)
}

func TestPath_GetGameIcon(t *testing.T) {
	result, err := GetGameIcon.Build(GameParams{Slug: "test"})
	require.NoError(t, err)
	assert.Equal(t, "/games/test/icon", result)
	assert.Equal(t, "GET", GetGameIcon.Method)
}

func TestPath_PutGameIcon(t *testing.T) {
	result, err := PutGameIcon.Build(GameParams{Slug: "test"})
	require.NoError(t, err)
	assert.Equal(t, "/games/test/icon", result)
	assert.Equal(t, "PUT", PutGameIcon.Method)
}

func TestPath_GetGameBlob(t *testing.T) {
	result, err := GetGameBlob.Build(GameParams{Slug: "test"})
	require.NoError(t, err)
	assert.Equal(t, "/games/test/blob", result)
	assert.Equal(t, "GET", GetGameBlob.Method)
}

func TestPath_PutGameBlob(t *testing.T) {
	result, err := PutGameBlob.Build(GameParams{Slug: "test"})
	require.NoError(t, err)
	assert.Equal(t, "/games/test/blob", result)
	assert.Equal(t, "PUT", PutGameBlob.Method)
}

func TestPath_GetUsers(t *testing.T) {
	result, err := GetUsers.Build(nil)
	require.NoError(t, err)
	assert.Equal(t, "/users", result)
	assert.Equal(t, "GET", GetUsers.Method)
}

func TestPath_PutUser(t *testing.T) {
	result, err := PutUser.Build(UserParams{IPv4Address: "192.168.1.1"})
	require.NoError(t, err)
	assert.Equal(t, "/users/192.168.1.1", result)
	assert.Equal(t, "PUT", PutUser.Method)
}

func TestPath_DeleteUser(t *testing.T) {
	result, err := DeleteUser.Build(UserParams{IPv4Address: "192.168.1.1"})
	require.NoError(t, err)
	assert.Equal(t, "/users/192.168.1.1", result)
	assert.Equal(t, "DELETE", DeleteUser.Method)
}

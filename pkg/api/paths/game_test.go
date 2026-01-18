package paths

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameParams_Apply(t *testing.T) {
	t.Run("valid slug", func(t *testing.T) {
		params := GameParams{Slug: "test-game"}
		result, err := params.Apply("/games/:slug")
		require.NoError(t, err)
		assert.Equal(t, "/games/test-game", result)
	})

	t.Run("empty slug", func(t *testing.T) {
		params := GameParams{Slug: ""}
		result, err := params.Apply("/games/:slug")
		require.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "slug is required")
	})

	t.Run("path without slug param", func(t *testing.T) {
		params := GameParams{Slug: "test"}
		result, err := params.Apply("/games")
		require.NoError(t, err)
		assert.Equal(t, "/games", result)
	})

	t.Run("path with multiple params", func(t *testing.T) {
		params := GameParams{Slug: "test-game"}
		result, err := params.Apply("/games/:slug/blob")
		require.NoError(t, err)
		assert.Equal(t, "/games/test-game/blob", result)
	})

	t.Run("path with slug in middle", func(t *testing.T) {
		params := GameParams{Slug: "test-game"}
		result, err := params.Apply("/api/v1/games/:slug")
		require.NoError(t, err)
		assert.Equal(t, "/api/v1/games/test-game", result)
	})

	t.Run("complex path", func(t *testing.T) {
		params := GameParams{Slug: "my-test-game"}
		result, err := params.Apply("/games/:slug/icon")
		require.NoError(t, err)
		assert.Equal(t, "/games/my-test-game/icon", result)
	})
}

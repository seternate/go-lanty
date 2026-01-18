package paths

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserParams_Apply(t *testing.T) {
	t.Run("valid IPv4 address", func(t *testing.T) {
		params := UserParams{IPv4Address: "192.168.1.1"}
		result, err := params.Apply("/users/:ipv4Address")
		require.NoError(t, err)
		assert.Equal(t, "/users/192.168.1.1", result)
	})

	t.Run("empty IPv4 address", func(t *testing.T) {
		params := UserParams{IPv4Address: ""}
		result, err := params.Apply("/users/:ipv4Address")
		require.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "ipv4Address is required")
	})

	t.Run("path without ipv4Address param", func(t *testing.T) {
		params := UserParams{IPv4Address: "192.168.1.1"}
		result, err := params.Apply("/users")
		require.NoError(t, err)
		assert.Equal(t, "/users", result)
	})

	t.Run("path with multiple params", func(t *testing.T) {
		params := UserParams{IPv4Address: "192.168.1.1"}
		result, err := params.Apply("/users/:ipv4Address/info")
		require.NoError(t, err)
		assert.Equal(t, "/users/192.168.1.1/info", result)
	})

	t.Run("path with ipv4Address in middle", func(t *testing.T) {
		params := UserParams{IPv4Address: "10.0.0.1"}
		result, err := params.Apply("/api/v1/users/:ipv4Address")
		require.NoError(t, err)
		assert.Equal(t, "/api/v1/users/10.0.0.1", result)
	})

	t.Run("different IPv4 addresses", func(t *testing.T) {
		testCases := []string{
			"127.0.0.1",
			"0.0.0.0",
			"255.255.255.255",
			"10.20.30.40",
		}

		for _, addr := range testCases {
			t.Run(addr, func(t *testing.T) {
				params := UserParams{IPv4Address: addr}
				result, err := params.Apply("/users/:ipv4Address")
				require.NoError(t, err)
				assert.Equal(t, "/users/"+addr, result)
			})
		}
	})
}

package usermodel

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_JSON(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		user := User{
			IPv4Address: "192.168.1.1",
			Username:    "testuser",
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "192.168.1.1", jsonMap["ipv4Address"])
		assert.Equal(t, "testuser", jsonMap["username"])
	})

	t.Run("unmarshaling", func(t *testing.T) {
		jsonData := `{"ipv4Address":"192.168.1.1","username":"testuser"}`

		var user User
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, "192.168.1.1", user.IPv4Address)
		assert.Equal(t, "testuser", user.Username)
	})

	t.Run("round trip", func(t *testing.T) {
		original := User{
			IPv4Address: "10.0.0.1",
			Username:    "roundtrip",
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled User
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, original.IPv4Address, unmarshaled.IPv4Address)
		assert.Equal(t, original.Username, unmarshaled.Username)
	})
}

func TestUpsertUserRequest_JSON(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		req := UpsertUserRequest{
			Username: "newuser",
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "newuser", jsonMap["username"])
	})

	t.Run("unmarshaling", func(t *testing.T) {
		jsonData := `{"username":"newuser"}`

		var req UpsertUserRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, "newuser", req.Username)
	})

	t.Run("round trip", func(t *testing.T) {
		original := UpsertUserRequest{
			Username: "roundtrip",
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var unmarshaled UpsertUserRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, original.Username, unmarshaled.Username)
	})
}

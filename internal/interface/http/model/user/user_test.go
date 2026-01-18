package user

import (
	"net"
	"testing"
	"time"

	appUserSrv "github.com/seternate/go-lanty/internal/application/user"
	apimodel "github.com/seternate/go-lanty/pkg/api/models/user"
	domainuser "github.com/seternate/go-lanty/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("all fields populated", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		userView := &appUserSrv.UserView{
			IPv4Address: "192.168.1.1",
			Username:    "testuser",
			CreatedAt:   createdAt,
		}

		result := NewUser(userView)

		require.NotNil(t, result)
		assert.Equal(t, "192.168.1.1", result.IPv4Address)
		assert.Equal(t, "testuser", result.Username)
	})

	t.Run("minimal fields", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		userView := &appUserSrv.UserView{
			IPv4Address: "10.0.0.1",
			Username:    "minimaluser",
			CreatedAt:   createdAt,
		}

		result := NewUser(userView)

		require.NotNil(t, result)
		assert.Equal(t, "10.0.0.1", result.IPv4Address)
		assert.Equal(t, "minimaluser", result.Username)
	})

	t.Run("empty strings", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		userView := &appUserSrv.UserView{
			IPv4Address: "",
			Username:    "",
			CreatedAt:   createdAt,
		}

		result := NewUser(userView)

		require.NotNil(t, result)
		assert.Equal(t, "", result.IPv4Address)
		assert.Equal(t, "", result.Username)
	})
}

func TestNewUserFromDomain(t *testing.T) {
	t.Run("all fields populated", func(t *testing.T) {
		ip := net.ParseIP("192.168.1.1").To4()
		require.NotNil(t, ip)

		domainUser := &domainuser.User{
			IPv4Address: ip,
			Username:    "testuser",
		}

		result := NewUserFromDomain(domainUser)

		require.NotNil(t, result)
		assert.Equal(t, "192.168.1.1", result.IPv4Address)
		assert.Equal(t, "testuser", result.Username)
	})

	t.Run("different IPv4 address", func(t *testing.T) {
		ip := net.ParseIP("10.0.0.1").To4()
		require.NotNil(t, ip)

		domainUser := &domainuser.User{
			IPv4Address: ip,
			Username:    "anotheruser",
		}

		result := NewUserFromDomain(domainUser)

		require.NotNil(t, result)
		assert.Equal(t, "10.0.0.1", result.IPv4Address)
		assert.Equal(t, "anotheruser", result.Username)
	})

	t.Run("localhost IPv4 address", func(t *testing.T) {
		ip := net.ParseIP("127.0.0.1").To4()
		require.NotNil(t, ip)

		domainUser := &domainuser.User{
			IPv4Address: ip,
			Username:    "localhostuser",
		}

		result := NewUserFromDomain(domainUser)

		require.NotNil(t, result)
		assert.Equal(t, "127.0.0.1", result.IPv4Address)
		assert.Equal(t, "localhostuser", result.Username)
	})

	t.Run("empty username", func(t *testing.T) {
		ip := net.ParseIP("192.168.1.1").To4()
		require.NotNil(t, ip)

		domainUser := &domainuser.User{
			IPv4Address: ip,
			Username:    "",
		}

		result := NewUserFromDomain(domainUser)

		require.NotNil(t, result)
		assert.Equal(t, "192.168.1.1", result.IPv4Address)
		assert.Equal(t, "", result.Username)
	})
}

func TestToCommand(t *testing.T) {
	t.Run("all fields populated", func(t *testing.T) {
		req := &apimodel.UpsertUserRequest{
			Username: "testuser",
		}

		ipv4Address := "192.168.1.1"
		cmd := ToCommand(req, ipv4Address)

		assert.Equal(t, ipv4Address, cmd.IPv4Address)
		assert.Equal(t, "testuser", cmd.Username)
	})

	t.Run("different username", func(t *testing.T) {
		req := &apimodel.UpsertUserRequest{
			Username: "anotheruser",
		}

		ipv4Address := "10.0.0.1"
		cmd := ToCommand(req, ipv4Address)

		assert.Equal(t, ipv4Address, cmd.IPv4Address)
		assert.Equal(t, "anotheruser", cmd.Username)
	})

	t.Run("empty username", func(t *testing.T) {
		req := &apimodel.UpsertUserRequest{
			Username: "",
		}

		ipv4Address := "192.168.1.1"
		cmd := ToCommand(req, ipv4Address)

		assert.Equal(t, ipv4Address, cmd.IPv4Address)
		assert.Equal(t, "", cmd.Username)
	})

	t.Run("empty IPv4 address", func(t *testing.T) {
		req := &apimodel.UpsertUserRequest{
			Username: "testuser",
		}

		ipv4Address := ""
		cmd := ToCommand(req, ipv4Address)

		assert.Equal(t, "", cmd.IPv4Address)
		assert.Equal(t, "testuser", cmd.Username)
	})

	t.Run("localhost IPv4 address", func(t *testing.T) {
		req := &apimodel.UpsertUserRequest{
			Username: "localhostuser",
		}

		ipv4Address := "127.0.0.1"
		cmd := ToCommand(req, ipv4Address)

		assert.Equal(t, "127.0.0.1", cmd.IPv4Address)
		assert.Equal(t, "localhostuser", cmd.Username)
	})
}

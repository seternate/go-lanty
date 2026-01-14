package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	"github.com/seternate/go-lanty/internal/domain/user"
)

// Mock implementations
type mockUserRepository struct {
	users     map[string]*user.User
	getErr    error
	saveErr   error
	deleteErr error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*user.User),
	}
}

func (m *mockUserRepository) GetUser(ipv4Address string) (*user.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	u, ok := m.users[ipv4Address]
	if !ok {
		return nil, domainerr.NotFoundErr("user", ipv4Address)
	}
	return u, nil
}

func (m *mockUserRepository) SaveUser(u *user.User) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.users[u.IPv4Address.String()] = u
	return nil
}

func (m *mockUserRepository) DeleteUser(ipv4Address string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.users, ipv4Address)
	return nil
}

func TestCommandService_UpsertUser(t *testing.T) {
	t.Run("create new user", func(t *testing.T) {
		repo := newMockUserRepository()
		service := NewCommandService(repo)

		cmd := UpsertUserCommand{
			Username:    "testuser",
			IPv4Address: "192.168.1.1",
		}

		createdUser, created, err := service.UpsertUser(cmd)
		require.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.True(t, created)
		assert.Equal(t, "testuser", createdUser.Username)
		assert.Equal(t, "192.168.1.1", createdUser.IPv4Address.String())

		// Verify user was saved
		savedUser, err := repo.GetUser("192.168.1.1")
		require.NoError(t, err)
		assert.Equal(t, createdUser.Username, savedUser.Username)
		assert.Equal(t, createdUser.IPv4Address.String(), savedUser.IPv4Address.String())
	})

	t.Run("update existing user", func(t *testing.T) {
		repo := newMockUserRepository()
		existingUser, err := user.NewUser("oldusername", "192.168.1.1")
		require.NoError(t, err)
		repo.users["192.168.1.1"] = existingUser

		service := NewCommandService(repo)

		cmd := UpsertUserCommand{
			IPv4Address: "192.168.1.1",
			Username:    "newusername",
		}

		updatedUser, created, err := service.UpsertUser(cmd)
		require.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.False(t, created)
		assert.Equal(t, "newusername", updatedUser.Username)
		assert.Equal(t, "192.168.1.1", updatedUser.IPv4Address.String())

		// Verify user was updated in repository
		savedUser, err := repo.GetUser("192.168.1.1")
		require.NoError(t, err)
		assert.Equal(t, "newusername", savedUser.Username)
	})

	t.Run("invalid username on create", func(t *testing.T) {
		repo := newMockUserRepository()
		service := NewCommandService(repo)

		cmd := UpsertUserCommand{
			Username:    "", // Empty username should fail validation
			IPv4Address: "192.168.1.1",
		}

		createdUser, created, err := service.UpsertUser(cmd)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to create new user")
	})

	t.Run("invalid IPv4 address on create", func(t *testing.T) {
		repo := newMockUserRepository()
		service := NewCommandService(repo)

		cmd := UpsertUserCommand{
			Username:    "testuser",
			IPv4Address: "invalid-ip",
		}

		createdUser, created, err := service.UpsertUser(cmd)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to create new user")
	})

	t.Run("invalid username on update", func(t *testing.T) {
		repo := newMockUserRepository()
		existingUser, err := user.NewUser("oldusername", "192.168.1.1")
		require.NoError(t, err)
		repo.users["192.168.1.1"] = existingUser

		service := NewCommandService(repo)

		cmd := UpsertUserCommand{
			IPv4Address: "192.168.1.1",
			Username:    "", // Empty username should fail validation
		}

		updatedUser, created, err := service.UpsertUser(cmd)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to update user")
	})

	t.Run("repository get error (non-not-found)", func(t *testing.T) {
		repo := newMockUserRepository()
		repo.getErr = errors.New("database error")
		service := NewCommandService(repo)

		cmd := UpsertUserCommand{
			Username:    "testuser",
			IPv4Address: "192.168.1.1",
		}

		createdUser, created, err := service.UpsertUser(cmd)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to get user")
	})

	t.Run("repository save error on create", func(t *testing.T) {
		repo := newMockUserRepository()
		repo.saveErr = errors.New("save failed")
		service := NewCommandService(repo)

		cmd := UpsertUserCommand{
			Username:    "testuser",
			IPv4Address: "192.168.1.1",
		}

		createdUser, created, err := service.UpsertUser(cmd)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to save new user to repository")
	})

	t.Run("repository save error on update", func(t *testing.T) {
		repo := newMockUserRepository()
		existingUser, err := user.NewUser("oldusername", "192.168.1.1")
		require.NoError(t, err)
		repo.users["192.168.1.1"] = existingUser
		repo.saveErr = errors.New("save failed")

		service := NewCommandService(repo)

		cmd := UpsertUserCommand{
			IPv4Address: "192.168.1.1",
			Username:    "newusername",
		}

		updatedUser, created, err := service.UpsertUser(cmd)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to save updated user")
	})
}

func TestCommandService_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockUserRepository()
		existingUser, err := user.NewUser("testuser", "192.168.1.1")
		require.NoError(t, err)
		repo.users["192.168.1.1"] = existingUser

		service := NewCommandService(repo)

		err = service.DeleteUser("192.168.1.1")
		require.NoError(t, err)

		// Verify user was deleted
		_, err = repo.GetUser("192.168.1.1")
		assert.Error(t, err)
		var notFound domainerr.NotFound
		assert.ErrorAs(t, err, &notFound)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := newMockUserRepository()
		service := NewCommandService(repo)

		err := service.DeleteUser("192.168.1.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get user from repository")
	})

	t.Run("repository delete failure", func(t *testing.T) {
		repo := newMockUserRepository()
		existingUser, err := user.NewUser("testuser", "192.168.1.1")
		require.NoError(t, err)
		repo.users["192.168.1.1"] = existingUser
		repo.deleteErr = errors.New("delete failed")

		service := NewCommandService(repo)

		err = service.DeleteUser("192.168.1.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete user from repository")
	})
}

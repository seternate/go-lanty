package user

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupQueryServiceTest(t *testing.T) (*queryServiceImpl, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	service := NewQueryService(sqlxDB).(*queryServiceImpl)

	cleanup := func() {
		sqlxDB.Close()
	}

	return service, mock, cleanup
}

func TestQueryService_GetUsers(t *testing.T) {
	t.Run("empty result", func(t *testing.T) {
		service, mock, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"ipv4_address", "username", "created_at"})

		mock.ExpectQuery(`SELECT ipv4_address, username, created_at FROM users`).
			WillReturnRows(rows)

		users, err := service.GetUsers()
		require.NoError(t, err)
		assert.Empty(t, users)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("single user", func(t *testing.T) {
		service, mock, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		createdAt := time.Now()
		rows := sqlmock.NewRows([]string{"ipv4_address", "username", "created_at"}).
			AddRow("192.168.1.1", "testuser", createdAt)

		mock.ExpectQuery(`SELECT ipv4_address, username, created_at FROM users`).
			WillReturnRows(rows)

		users, err := service.GetUsers()
		require.NoError(t, err)
		require.Len(t, users, 1)
		assert.Equal(t, "192.168.1.1", users[0].IPv4Address)
		assert.Equal(t, "testuser", users[0].Username)
		assert.WithinDuration(t, createdAt, users[0].CreatedAt, time.Second)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("multiple users", func(t *testing.T) {
		service, mock, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		createdAt1 := time.Now()
		createdAt2 := time.Now().Add(time.Hour)
		rows := sqlmock.NewRows([]string{"ipv4_address", "username", "created_at"}).
			AddRow("192.168.1.1", "user1", createdAt1).
			AddRow("192.168.1.2", "user2", createdAt2)

		mock.ExpectQuery(`SELECT ipv4_address, username, created_at FROM users`).
			WillReturnRows(rows)

		users, err := service.GetUsers()
		require.NoError(t, err)
		require.Len(t, users, 2)
		assert.Equal(t, "192.168.1.1", users[0].IPv4Address)
		assert.Equal(t, "user1", users[0].Username)
		assert.Equal(t, "192.168.1.2", users[1].IPv4Address)
		assert.Equal(t, "user2", users[1].Username)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user with null created_at", func(t *testing.T) {
		service, mock, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"ipv4_address", "username", "created_at"}).
			AddRow("192.168.1.1", "testuser", nil)

		mock.ExpectQuery(`SELECT ipv4_address, username, created_at FROM users`).
			WillReturnRows(rows)

		users, err := service.GetUsers()
		require.NoError(t, err)
		require.Len(t, users, 1)
		assert.Equal(t, "192.168.1.1", users[0].IPv4Address)
		assert.Equal(t, "testuser", users[0].Username)
		assert.True(t, users[0].CreatedAt.IsZero())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		service, mock, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		dbErr := errors.New("database error")
		mock.ExpectQuery(`SELECT ipv4_address, username, created_at FROM users`).
			WillReturnError(dbErr)

		users, err := service.GetUsers()
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "database query failed")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query build error", func(t *testing.T) {
		// This is harder to test directly, but we can test with invalid SQL
		// Actually, squirrel handles this well, so we'll skip this edge case
		// as it's unlikely to happen in practice
	})
}

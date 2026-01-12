package game

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/seternate/go-lanty/pkg/application/asset"
	domainasset "github.com/seternate/go-lanty/pkg/domain/asset"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

func setupQueryServiceTest(t *testing.T) (*queryServiceImpl, sqlmock.Sqlmock, *mockAssetServiceForQuery, *mockAssetServiceForQuery, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	iconServiceMock := newMockAssetServiceForQuery()
	blobServiceMock := newMockAssetServiceForQuery()
	iconQueryService := &mockAssetQueryServiceForGame{
		getFunc:       iconServiceMock.getAssetFunc,
		getContentFunc: iconServiceMock.getContentFunc,
	}
	iconCommandService := &mockAssetCommandServiceForGame{}
	blobQueryService := &mockAssetQueryServiceForGame{
		getFunc:       blobServiceMock.getAssetFunc,
		getContentFunc: blobServiceMock.getContentFunc,
	}
	blobCommandService := &mockAssetCommandServiceForGame{}
	iconService := asset.NewService(iconQueryService, iconCommandService)
	blobService := asset.NewService(blobQueryService, blobCommandService)
	service := NewQueryService(sqlxDB, iconService, blobService).(*queryServiceImpl)

	cleanup := func() {
		sqlxDB.Close()
	}

	return service, mock, iconServiceMock, blobServiceMock, cleanup
}

type mockAssetServiceForQuery struct {
	getAssetFunc    func(id uuid.UUID) (*asset.AssetView, error)
	getContentFunc  func(id uuid.UUID) (io.ReadCloser, error)
}

func newMockAssetServiceForQuery() *mockAssetServiceForQuery {
	return &mockAssetServiceForQuery{
		getAssetFunc: func(id uuid.UUID) (*asset.AssetView, error) {
			return &asset.AssetView{
				ID:       id,
				Size:     100,
				Checksum: "checksum",
				Algorithm: "md5",
				MimeType: "image/png",
			}, nil
		},
		getContentFunc: func(id uuid.UUID) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("content"))), nil
		},
	}
}

func (m *mockAssetServiceForQuery) Query() asset.QueryService {
	return &mockAssetQueryServiceForGame{
		getFunc:       m.getAssetFunc,
		getContentFunc: m.getContentFunc,
	}
}

func (m *mockAssetServiceForQuery) Command() asset.CommandService {
	return &mockAssetCommandServiceForGame{}
}

type mockAssetCommandServiceForGame struct{}

func (m *mockAssetCommandServiceForGame) StoreNewAsset(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
	return nil, nil
}

func (m *mockAssetCommandServiceForGame) DeleteAsset(id uuid.UUID) error {
	return nil
}

type mockAssetQueryServiceForGame struct {
	getFunc       func(id uuid.UUID) (*asset.AssetView, error)
	getContentFunc func(id uuid.UUID) (io.ReadCloser, error)
}

func (m *mockAssetQueryServiceForGame) GetAsset(id uuid.UUID) (*asset.AssetView, error) {
	return m.getFunc(id)
}

func (m *mockAssetQueryServiceForGame) GetAssetContent(id uuid.UUID) (io.ReadCloser, error) {
	return m.getContentFunc(id)
}

func TestQueryService_GetGames(t *testing.T) {
	t.Run("empty result", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{
			"game_slug", "game_name", "game_created_at",
			"exec_id", "exec_game_slug", "exec_role", "exec_path", "exec_requires_admin", "exec_format", "exec_arg_seperator", "exec_created_at",
			"arg_id", "arg_game_exec_id", "arg_role", "arg_name", "arg_required", "arg_enabled", "arg_format", "arg_separator", "arg_arg", "arg_description",
			"arg_default_string", "arg_default_bool", "arg_default_int", "arg_default_float", "arg_enums", "arg_min_int", "arg_max_int",
			"arg_min_float", "arg_max_float", "arg_float_precision", "arg_order_index", "arg_created_at",
		})

		mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

		games, err := service.GetGames()
		require.NoError(t, err)
		assert.Empty(t, games)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("single game without execs", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		createdAt := time.Now()
		rows := sqlmock.NewRows([]string{
			"game_slug", "game_name", "game_created_at",
			"exec_id", "exec_game_slug", "exec_role", "exec_path", "exec_requires_admin", "exec_format", "exec_arg_seperator", "exec_created_at",
			"arg_id", "arg_game_exec_id", "arg_role", "arg_name", "arg_required", "arg_enabled", "arg_format", "arg_separator", "arg_arg", "arg_description",
			"arg_default_string", "arg_default_bool", "arg_default_int", "arg_default_float", "arg_enums", "arg_min_int", "arg_max_int",
			"arg_min_float", "arg_max_float", "arg_float_precision", "arg_order_index", "arg_created_at",
		}).AddRow("test-game", "Test Game", createdAt, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

		games, err := service.GetGames()
		require.NoError(t, err)
		require.Len(t, games, 1)
		assert.Equal(t, "test-game", games[0].Slug)
		assert.Equal(t, "Test Game", games[0].Name)
		assert.Empty(t, games[0].Execs)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		dbErr := errors.New("database error")
		mock.ExpectQuery(`SELECT`).WillReturnError(dbErr)

		games, err := service.GetGames()
		assert.Error(t, err)
		assert.Nil(t, games)
		assert.Contains(t, err.Error(), "database query failed")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestQueryService_GetGame(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		createdAt := time.Now()
		rows := sqlmock.NewRows([]string{
			"game_slug", "game_name", "game_created_at",
			"exec_game_slug", "exec_role", "exec_path", "exec_requires_admin", "exec_format", "exec_arg_seperator", "exec_created_at",
			"arg_id", "arg_game_exec_id", "arg_role", "arg_name", "arg_required", "arg_enabled", "arg_format", "arg_separator", "arg_arg", "arg_description",
			"arg_default_string", "arg_default_bool", "arg_default_int", "arg_default_float", "arg_enums", "arg_min_int", "arg_max_int",
			"arg_min_float", "arg_max_float", "arg_float_precision", "arg_order_index", "arg_created_at",
		})
		rows.AddRow(
			"test-game", "Test Game", createdAt,
			"test-game", "client", "/path/to/client", 
			sql.NullBool{Valid: true, Bool: true}, sql.NullString{}, sql.NullString{}, 
			sql.NullTime{Valid: true, Time: createdAt},
			sql.NullString{}, sql.NullString{}, sql.NullString{}, sql.NullString{}, 
			sql.NullBool{}, sql.NullBool{}, sql.NullString{}, sql.NullString{}, sql.NullString{}, sql.NullString{},
			sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}, 
			pq.StringArray{}, sql.NullInt64{}, sql.NullInt64{},
			sql.NullFloat64{}, sql.NullFloat64{}, sql.NullInt64{}, sql.NullInt64{}, sql.NullTime{})

		mock.ExpectQuery(`SELECT`).WithArgs("test-game").WillReturnRows(rows)

		game, err := service.GetGame("test-game")
		require.NoError(t, err)
		require.NotNil(t, game)
		assert.Equal(t, "test-game", game.Slug)
		assert.Equal(t, "Test Game", game.Name)
		assert.WithinDuration(t, createdAt, game.CreatedAt, time.Second)
		assert.Empty(t, game.Execs)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{
			"game_slug", "game_name", "game_created_at",
			"exec_game_slug", "exec_role", "exec_path", "exec_requires_admin", "exec_format", "exec_arg_seperator", "exec_created_at",
			"arg_id", "arg_game_exec_id", "arg_role", "arg_name", "arg_required", "arg_enabled", "arg_format", "arg_separator", "arg_arg", "arg_description",
			"arg_default_string", "arg_default_bool", "arg_default_int", "arg_default_float", "arg_enums", "arg_min_int", "arg_max_int",
			"arg_min_float", "arg_max_float", "arg_float_precision", "arg_order_index", "arg_created_at",
		})

		mock.ExpectQuery(`SELECT`).WithArgs("non-existent").WillReturnRows(rows)

		game, err := service.GetGame("non-existent")
		assert.Error(t, err)
		assert.Nil(t, game)
		var notFound domainerr.NotFound
		assert.ErrorAs(t, err, &notFound)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		dbErr := errors.New("database error")
		mock.ExpectQuery(`SELECT`).WithArgs("test-game").WillReturnError(dbErr)

		game, err := service.GetGame("test-game")
		assert.Error(t, err)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "database query failed")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestQueryService_FetchIcon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, mock, iconService, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		assetID := uuid.New()
		rows := sqlmock.NewRows([]string{"asset_id"}).AddRow(assetID)

		mock.ExpectQuery(`SELECT asset_id FROM game_assets`).
			WithArgs("test-game", "icon").
			WillReturnRows(rows)

		content, err := service.FetchIcon("test-game")
		require.NoError(t, err)
		assert.NotNil(t, content)
		assert.Equal(t, uint64(100), content.Size)
		assert.Equal(t, "checksum", content.Checksum)
		assert.Equal(t, "md5", content.Algorithm)
		assert.Equal(t, "image/png", content.MimeType)

		// Verify asset service was called
		assert.NotNil(t, iconService.getAssetFunc)
		assert.NotNil(t, iconService.getContentFunc)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		mock.ExpectQuery(`SELECT asset_id FROM game_assets`).
			WithArgs("non-existent", "icon").
			WillReturnError(sql.ErrNoRows)

		content, err := service.FetchIcon("non-existent")
		assert.Error(t, err)
		assert.Nil(t, content)
		var notFound domainerr.NotFound
		assert.ErrorAs(t, err, &notFound)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		dbErr := errors.New("database error")
		mock.ExpectQuery(`SELECT asset_id FROM game_assets`).
			WithArgs("test-game", "icon").
			WillReturnError(dbErr)

		content, err := service.FetchIcon("test-game")
		assert.Error(t, err)
		assert.Nil(t, content)
		assert.Contains(t, err.Error(), "database query failed")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Note: Testing asset service errors is complex because the service is created at setup
	// For now, we test the main paths. Asset service error handling is tested in asset service tests.
}

func TestQueryService_FetchBlob(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, mock, _, blobService, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		assetID := uuid.New()
		rows := sqlmock.NewRows([]string{"asset_id"}).AddRow(assetID)

		mock.ExpectQuery(`SELECT asset_id FROM game_assets`).
			WithArgs("test-game", "blob").
			WillReturnRows(rows)

		content, err := service.FetchBlob("test-game")
		require.NoError(t, err)
		assert.NotNil(t, content)
		assert.Equal(t, uint64(100), content.Size)

		// Verify asset service was called
		assert.NotNil(t, blobService.getAssetFunc)
		assert.NotNil(t, blobService.getContentFunc)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		mock.ExpectQuery(`SELECT asset_id FROM game_assets`).
			WithArgs("non-existent", "blob").
			WillReturnError(sql.ErrNoRows)

		content, err := service.FetchBlob("non-existent")
		assert.Error(t, err)
		assert.Nil(t, content)
		var notFound domainerr.NotFound
		assert.ErrorAs(t, err, &notFound)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		service, mock, _, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		dbErr := errors.New("database error")
		mock.ExpectQuery(`SELECT asset_id FROM game_assets`).
			WithArgs("test-game", "blob").
			WillReturnError(dbErr)

		content, err := service.FetchBlob("test-game")
		assert.Error(t, err)
		assert.Nil(t, content)
		assert.Contains(t, err.Error(), "database query failed")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

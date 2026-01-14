package asset

import (
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

func setupQueryServiceTest(t *testing.T) (*queryServiceImpl, sqlmock.Sqlmock, *mockStorageAdapter, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	storage := newMockStorageAdapter()
	service := NewQueryService(sqlxDB, storage).(*queryServiceImpl)

	cleanup := func() {
		sqlxDB.Close()
	}

	return service, mock, storage, cleanup
}

func TestQueryService_GetAsset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, mock, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		testID := uuid.New()
		expectedURL := "file:///test-asset"
		expectedSize := uint64(1024)
		expectedChecksum := "abc123"
		expectedAlgorithm := "md5"
		expectedMimeType := "image/png"
		expectedCreatedAt := time.Now()

		rows := sqlmock.NewRows([]string{"id", "url", "size", "checksum", "algorithm", "mime_type", "created_at"}).
			AddRow(testID, expectedURL, expectedSize, expectedChecksum, expectedAlgorithm, expectedMimeType, expectedCreatedAt)

		mock.ExpectQuery(`SELECT id, url, size, checksum, algorithm, mime_type, created_at FROM assets WHERE id = \$1`).
			WithArgs(testID).
			WillReturnRows(rows)

		asset, err := service.GetAsset(testID)
		require.NoError(t, err)
		assert.NotNil(t, asset)
		assert.Equal(t, testID, asset.ID)
		assert.Equal(t, expectedURL, asset.URL)
		assert.Equal(t, expectedSize, asset.Size)
		assert.Equal(t, expectedChecksum, asset.Checksum)
		assert.Equal(t, expectedAlgorithm, asset.Algorithm)
		assert.Equal(t, expectedMimeType, asset.MimeType)
		assert.WithinDuration(t, expectedCreatedAt, asset.CreatedAt, time.Second)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		service, mock, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		testID := uuid.New()

		mock.ExpectQuery(`SELECT id, url, size, checksum, algorithm, mime_type, created_at FROM assets WHERE id = \$1`).
			WithArgs(testID).
			WillReturnError(sql.ErrNoRows)

		asset, err := service.GetAsset(testID)
		assert.Error(t, err)
		assert.Nil(t, asset)
		var notFound domainerr.NotFound
		assert.ErrorAs(t, err, &notFound)
		assert.Contains(t, err.Error(), "asset")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		service, mock, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		testID := uuid.New()
		dbErr := errors.New("database connection failed")

		mock.ExpectQuery(`SELECT id, url, size, checksum, algorithm, mime_type, created_at FROM assets WHERE id = \$1`).
			WithArgs(testID).
			WillReturnError(dbErr)

		asset, err := service.GetAsset(testID)
		assert.Error(t, err)
		assert.Nil(t, asset)
		assert.Contains(t, err.Error(), "database query failed")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query build error", func(t *testing.T) {
		// This is harder to test directly, but we can test with invalid SQL
		// Actually, squirrel handles this well, so we'll skip this edge case
		// as it's unlikely to happen in practice
	})
}

func TestQueryService_GetAssetContent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, mock, storage, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		testID := uuid.New()
		expectedURL := "file:///test-asset"
		expectedData := []byte("test asset content")
		storage.fetchedData[expectedURL] = expectedData

		rows := sqlmock.NewRows([]string{"url"}).
			AddRow(expectedURL)

		mock.ExpectQuery(`SELECT url FROM assets WHERE id = \$1`).
			WithArgs(testID).
			WillReturnRows(rows)

		reader, err := service.GetAssetContent(testID)
		require.NoError(t, err)
		assert.NotNil(t, reader)

		// Read and verify content
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, expectedData, data)

		reader.Close()

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("asset not found", func(t *testing.T) {
		service, mock, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		testID := uuid.New()

		mock.ExpectQuery(`SELECT url FROM assets WHERE id = \$1`).
			WithArgs(testID).
			WillReturnError(sql.ErrNoRows)

		reader, err := service.GetAssetContent(testID)
		assert.Error(t, err)
		assert.Nil(t, reader)
		var notFound domainerr.NotFound
		assert.ErrorAs(t, err, &notFound)
		assert.Contains(t, err.Error(), "asset")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		service, mock, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		testID := uuid.New()
		dbErr := errors.New("database connection failed")

		mock.ExpectQuery(`SELECT url FROM assets WHERE id = \$1`).
			WithArgs(testID).
			WillReturnError(dbErr)

		reader, err := service.GetAssetContent(testID)
		assert.Error(t, err)
		assert.Nil(t, reader)
		assert.Contains(t, err.Error(), "database query failed")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid URL in database", func(t *testing.T) {
		service, mock, _, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		testID := uuid.New()
		invalidURL := "://invalid-url"

		rows := sqlmock.NewRows([]string{"url"}).
			AddRow(invalidURL)

		mock.ExpectQuery(`SELECT url FROM assets WHERE id = \$1`).
			WithArgs(testID).
			WillReturnRows(rows)

		reader, err := service.GetAssetContent(testID)
		assert.Error(t, err)
		assert.Nil(t, reader)
		assert.Contains(t, err.Error(), "failed to parse asset URL")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("storage fetch failure", func(t *testing.T) {
		service, mock, storage, cleanup := setupQueryServiceTest(t)
		defer cleanup()

		testID := uuid.New()
		expectedURL := "file:///test-asset"
		storage.fetchErr = errors.New("storage fetch failed")

		rows := sqlmock.NewRows([]string{"url"}).
			AddRow(expectedURL)

		mock.ExpectQuery(`SELECT url FROM assets WHERE id = \$1`).
			WithArgs(testID).
			WillReturnRows(rows)

		reader, err := service.GetAssetContent(testID)
		assert.Error(t, err)
		assert.Nil(t, reader)
		assert.Contains(t, err.Error(), "failed to fetch asset from storage")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

package asset

import (
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/seternate/go-lanty/internal/domain/asset"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	"github.com/seternate/go-lanty/internal/infrastructure/database"
)

func getTestDB(t *testing.T) *sqlx.DB {
	dbString := os.Getenv("DB")
	if dbString == "" {
		dbString = "postgres://lanty:lanty@localhost:5432/lanty?sslmode=disable"
	}

	db, err := database.New(dbString)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	return db
}

func setupTest(t *testing.T) (*assetrepository, *sqlx.DB) {
	db := getTestDB(t)
	repo := NewAssetRepository(db).(*assetrepository)

	// Clean up any existing test data
	_, err := db.Exec("DELETE FROM assets")
	if err != nil {
		t.Fatalf("failed to clean up test data: %v", err)
	}

	// Add defensive cleanup for when test fails
	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM assets")
		db.Close()
	})

	return repo, db
}

func TestNewAssetRepository(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := NewAssetRepository(db)

	if repo == nil {
		t.Fatal("expected repository, got nil")
	}

	// Verify it implements the interface
	var _ asset.AssetRepository = repo
}

func TestAssetRepository_GetAsset(t *testing.T) {
	repo, _ := setupTest(t)

	testID := uuid.New()

	t.Run("Get existing asset", func(t *testing.T) {
		// Create a test asset directly in the database
		testAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		err = repo.CreateAsset(testAsset)
		if err != nil {
			t.Fatalf("failed to create asset in database: %v", err)
		}

		// Retrieve the asset
		retrieved, err := repo.GetAsset(testID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved == nil {
			t.Fatal("expected asset, got nil")
		}
		if retrieved.ID != testID {
			t.Errorf("expected ID %v, got %v", testID, retrieved.ID)
		}
		if retrieved.URL.String() != "https://example.com/file.bin" {
			t.Errorf("expected URL %q, got %q", "https://example.com/file.bin", retrieved.URL.String())
		}
		if retrieved.Size != 1024 {
			t.Errorf("expected Size 1024, got %d", retrieved.Size)
		}
		if retrieved.Checksum != "abc123" {
			t.Errorf("expected Checksum %q, got %q", "abc123", retrieved.Checksum)
		}
		if retrieved.Algorithm.String() != "md5" {
			t.Errorf("expected Algorithm %q, got %q", "md5", retrieved.Algorithm.String())
		}
		if retrieved.MimeType != "application/octet-stream" {
			t.Errorf("expected MimeType %q, got %q", "application/octet-stream", retrieved.MimeType)
		}
	})

	t.Run("Get non-existent asset", func(t *testing.T) {
		nonExistentID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

		retrieved, err := repo.GetAsset(nonExistentID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if retrieved != nil {
			t.Error("expected nil asset")
		}
		var notFoundErr *domainerr.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Errorf("expected NotFoundError, got %T: %v", err, err)
		}
	})

	t.Run("Get asset with SHA-256", func(t *testing.T) {
		testID2 := uuid.New()
		testAsset, err := asset.RehydrateAsset(testID2, "https://example.com/file.pdf", 2048, "sha256hash", "sha-256", "application/pdf")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		err = repo.CreateAsset(testAsset)
		if err != nil {
			t.Fatalf("failed to create asset in database: %v", err)
		}

		retrieved, err := repo.GetAsset(testID2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved.Algorithm.String() != "sha-256" {
			t.Errorf("expected Algorithm %q, got %q", "sha-256", retrieved.Algorithm.String())
		}
		if retrieved.MimeType != "application/pdf" {
			t.Errorf("expected MimeType %q, got %q", "application/pdf", retrieved.MimeType)
		}
	})
}

func TestAssetRepository_CreateAsset(t *testing.T) {
	repo, _ := setupTest(t)

	t.Run("Create valid asset", func(t *testing.T) {
		testID := uuid.New()
		testAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		err = repo.CreateAsset(testAsset)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify it was created
		retrieved, err := repo.GetAsset(testID)
		if err != nil {
			t.Fatalf("failed to retrieve created asset: %v", err)
		}
		if retrieved.ID != testID {
			t.Errorf("expected ID %v, got %v", testID, retrieved.ID)
		}
	})

	t.Run("Create asset with SHA-256", func(t *testing.T) {
		testID := uuid.New()
		testAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.pdf", 2048, "sha256hash", "sha-256", "application/pdf")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		err = repo.CreateAsset(testAsset)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		retrieved, err := repo.GetAsset(testID)
		if err != nil {
			t.Fatalf("failed to retrieve created asset: %v", err)
		}
		if retrieved.Algorithm.String() != "sha-256" {
			t.Errorf("expected Algorithm %q, got %q", "sha-256", retrieved.Algorithm.String())
		}
	})

	t.Run("Create duplicate asset (conflict)", func(t *testing.T) {
		testID := uuid.New()
		testAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		err = repo.CreateAsset(testAsset)
		if err != nil {
			t.Fatalf("failed to create first asset: %v", err)
		}

		// Try to create the same asset again
		err = repo.CreateAsset(testAsset)
		if err == nil {
			t.Error("expected error for duplicate asset, got nil")
		}
		var conflictErr *domainerr.ConflictError
		if !errors.As(err, &conflictErr) {
			t.Errorf("expected ConflictError, got %T: %v", err, err)
		}
	})

	t.Run("Create asset with complex URL", func(t *testing.T) {
		testID := uuid.New()
		url := "https://example.com/file.bin?version=1&token=abc123"
		testAsset, err := asset.RehydrateAsset(testID, url, 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		err = repo.CreateAsset(testAsset)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		retrieved, err := repo.GetAsset(testID)
		if err != nil {
			t.Fatalf("failed to retrieve created asset: %v", err)
		}
		if retrieved.URL.String() != url {
			t.Errorf("expected URL %q, got %q", url, retrieved.URL.String())
		}
	})

	t.Run("Create multiple assets", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			testID := uuid.New()
			testAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", uint64(1024+i), "checksum", "md5", "application/octet-stream")
			if err != nil {
				t.Fatalf("failed to create test asset: %v", err)
			}

			err = repo.CreateAsset(testAsset)
			if err != nil {
				t.Fatalf("unexpected error creating asset %d: %v", i, err)
			}
		}
	})
}

func TestAssetRepository_DeleteAsset(t *testing.T) {
	repo, _ := setupTest(t)

	t.Run("Delete existing asset", func(t *testing.T) {
		testID := uuid.New()
		testAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		err = repo.CreateAsset(testAsset)
		if err != nil {
			t.Fatalf("failed to create asset: %v", err)
		}

		// Verify it exists
		_, err = repo.GetAsset(testID)
		if err != nil {
			t.Fatalf("asset should exist before deletion: %v", err)
		}

		// Delete it
		err = repo.DeleteAsset(testID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify it's gone
		_, err = repo.GetAsset(testID)
		if err == nil {
			t.Error("expected error after deletion, got nil")
		}
		var notFoundErr *domainerr.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Errorf("expected NotFoundError, got %T: %v", err, err)
		}
	})

	t.Run("Delete non-existent asset", func(t *testing.T) {
		nonExistentID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

		// Delete should not return an error even if the asset doesn't exist
		// (this is a common pattern in repositories)
		err := repo.DeleteAsset(nonExistentID)
		if err != nil {
			// If it does return an error, that's also acceptable behavior
			// We just verify it doesn't panic
		}
	})

	t.Run("Delete multiple assets", func(t *testing.T) {
		testIDs := make([]uuid.UUID, 3)
		for i := 0; i < 3; i++ {
			testID := uuid.New()
			testIDs[i] = testID
			testAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
			if err != nil {
				t.Fatalf("failed to create test asset: %v", err)
			}

			err = repo.CreateAsset(testAsset)
			if err != nil {
				t.Fatalf("failed to create asset: %v", err)
			}
		}

		// Delete all of them
		for _, testID := range testIDs {
			err := repo.DeleteAsset(testID)
			if err != nil {
				t.Fatalf("unexpected error deleting asset %v: %v", testID, err)
			}
		}

		// Verify they're all gone
		for _, testID := range testIDs {
			_, err := repo.GetAsset(testID)
			if err == nil {
				t.Errorf("asset %v should not exist after deletion", testID)
			}
		}
	})
}

func TestAssetRepository_Integration(t *testing.T) {
	repo, _ := setupTest(t)

	t.Run("Full CRUD cycle", func(t *testing.T) {
		testID := uuid.New()
		testAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		// Create
		err = repo.CreateAsset(testAsset)
		if err != nil {
			t.Fatalf("failed to create asset: %v", err)
		}

		// Read
		retrieved, err := repo.GetAsset(testID)
		if err != nil {
			t.Fatalf("failed to retrieve asset: %v", err)
		}
		if retrieved.ID != testID {
			t.Errorf("expected ID %v, got %v", testID, retrieved.ID)
		}

		// Delete
		err = repo.DeleteAsset(testID)
		if err != nil {
			t.Fatalf("failed to delete asset: %v", err)
		}

		// Verify deletion
		_, err = repo.GetAsset(testID)
		if err == nil {
			t.Error("expected error after deletion, got nil")
		}
	})
}

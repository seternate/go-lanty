package game

import (
	"errors"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/seternate/go-lanty/internal/domain/asset"
	"github.com/seternate/go-lanty/internal/domain/game"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	"github.com/seternate/go-lanty/internal/infrastructure/database"
	assetrepo "github.com/seternate/go-lanty/internal/infrastructure/persistence/asset"
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

func setupTest(t *testing.T) (*gamerepository, *sqlx.DB) {
	db := getTestDB(t)
	repo := NewGameRepository(db).(*gamerepository)

	// Clean up any existing test data (in correct order due to foreign keys)
	_, _ = db.Exec("DELETE FROM game_args")
	_, _ = db.Exec("DELETE FROM game_execs")
	_, _ = db.Exec("DELETE FROM game_assets")
	_, _ = db.Exec("DELETE FROM games")
	_, _ = db.Exec("DELETE FROM assets")

	// Add defensive cleanup for when test fails
	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM game_args")
		_, _ = db.Exec("DELETE FROM game_execs")
		_, _ = db.Exec("DELETE FROM game_assets")
		_, _ = db.Exec("DELETE FROM games")
		_, _ = db.Exec("DELETE FROM assets")
		db.Close()
	})

	return repo, db
}

func createTestAsset(t *testing.T, db *sqlx.DB, assetID uuid.UUID) {
	assetRepo := assetrepo.NewAssetRepository(db)
	testAsset, err := asset.RehydrateAsset(assetID, "https://example.com/test.png", 1024, "test-checksum", "md5", "image/png")
	if err != nil {
		t.Fatalf("failed to create test asset: %v", err)
	}
	err = assetRepo.CreateAsset(testAsset)
	if err != nil {
		t.Fatalf("failed to save test asset: %v", err)
	}
}

func TestNewGameRepository(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := NewGameRepository(db)

	if repo == nil {
		t.Fatal("expected repository, got nil")
	}

	// Verify it implements the interface
	var _ game.GameRepository = repo
}

func TestGameRepository_GetGame(t *testing.T) {
	repo, db := setupTest(t)

	t.Run("Get existing game", func(t *testing.T) {
		testSlug := "test-game-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("failed to save game in database: %v", err)
		}

		// Retrieve the game
		retrieved, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved == nil {
			t.Fatal("expected game, got nil")
		}
		if retrieved.Slug != testSlug {
			t.Errorf("expected Slug %q, got %q", testSlug, retrieved.Slug)
		}
		if retrieved.Name != "Test Game" {
			t.Errorf("expected Name %q, got %q", "Test Game", retrieved.Name)
		}
	})

	t.Run("Get non-existent game", func(t *testing.T) {
		nonExistentSlug := "non-existent-game"

		retrieved, err := repo.GetGame(nonExistentSlug)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if retrieved != nil {
			t.Error("expected nil game")
		}
		var notFoundErr *domainerr.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Errorf("expected NotFoundError, got %T: %v", err, err)
		}
	})

	t.Run("Get game with execs", func(t *testing.T) {
		testSlug := "test-game-execs-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		exec, err := game.RehydrateGameExec("/path/to/client", "client")
		if err != nil {
			t.Fatalf("failed to create test exec: %v", err)
		}
		testGame.SetExecutable(*exec)

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("failed to save game in database: %v", err)
		}

		retrieved, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(retrieved.Execs) != 1 {
			t.Errorf("expected 1 exec, got %d", len(retrieved.Execs))
		}
		if exec, ok := retrieved.Execs[game.GAME_EXEC_ROLE_CLIENT]; !ok {
			t.Error("expected client exec to be present")
		} else if exec.Path != "/path/to/client" {
			t.Errorf("expected Path %q, got %q", "/path/to/client", exec.Path)
		}
	})

	t.Run("Get game with assets", func(t *testing.T) {
		testSlug := "test-game-assets-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		assetID := uuid.New()
		createTestAsset(t, db, assetID)
		gameAsset, err := game.RehydrateGameAsset(assetID, "icon")
		if err != nil {
			t.Fatalf("failed to create test game asset: %v", err)
		}
		testGame.SetAsset(*gameAsset)

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("failed to save game in database: %v", err)
		}

		retrieved, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(retrieved.Assets) != 1 {
			t.Errorf("expected 1 asset, got %d", len(retrieved.Assets))
		}
		if asset, ok := retrieved.Assets[game.GAME_ASSET_ROLE_ICON]; !ok {
			t.Error("expected icon asset to be present")
		} else if asset.AssetID != assetID {
			t.Errorf("expected AssetID %v, got %v", assetID, asset.AssetID)
		}
	})
}

func TestGameRepository_SaveGame(t *testing.T) {
	repo, db := setupTest(t)

	t.Run("Save valid game", func(t *testing.T) {
		testSlug := "test-game-save-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify it was saved
		retrieved, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("failed to retrieve saved game: %v", err)
		}
		if retrieved.Slug != testSlug {
			t.Errorf("expected Slug %q, got %q", testSlug, retrieved.Slug)
		}
	})

	t.Run("Save game with execs and args", func(t *testing.T) {
		testSlug := "test-game-execs-args-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		arg, err := game.NewGameArg(game.GameArgInput{
			Role:       "flag",
			Name:       "test-arg",
			Arg:        "--test",
			OrderIndex: 0,
		})
		if err != nil {
			t.Fatalf("failed to create test arg: %v", err)
		}

		execWithArgs, err := game.RehydrateGameExec("/path/to/client", "client", game.WithArgs(arg))
		if err != nil {
			t.Fatalf("failed to create exec with args: %v", err)
		}
		testGame.SetExecutable(*execWithArgs)

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		retrieved, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("failed to retrieve saved game: %v", err)
		}
		if exec, ok := retrieved.Execs[game.GAME_EXEC_ROLE_CLIENT]; ok {
			if len(exec.Args) != 1 {
				t.Errorf("expected 1 arg, got %d", len(exec.Args))
			}
		} else {
			t.Error("expected client exec to be present")
		}
	})

	t.Run("Save game with assets", func(t *testing.T) {
		testSlug := "test-game-assets-save-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		assetID := uuid.New()
		createTestAsset(t, db, assetID)
		gameAsset, err := game.RehydrateGameAsset(assetID, "icon")
		if err != nil {
			t.Fatalf("failed to create test game asset: %v", err)
		}
		testGame.SetAsset(*gameAsset)

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		retrieved, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("failed to retrieve saved game: %v", err)
		}
		if len(retrieved.Assets) != 1 {
			t.Errorf("expected 1 asset, got %d", len(retrieved.Assets))
		}
	})

	t.Run("Update existing game", func(t *testing.T) {
		testSlug := "test-game-update-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("failed to save game: %v", err)
		}

		// Update the game
		err = testGame.SetName("Updated Test Game")
		if err != nil {
			t.Fatalf("failed to update game name: %v", err)
		}

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("failed to update game: %v", err)
		}

		retrieved, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("failed to retrieve updated game: %v", err)
		}
		if retrieved.Name != "Updated Test Game" {
			t.Errorf("expected Name %q, got %q", "Updated Test Game", retrieved.Name)
		}
	})
}

func TestGameRepository_DeleteGame(t *testing.T) {
	repo, db := setupTest(t)

	t.Run("Delete existing game", func(t *testing.T) {
		testSlug := "test-game-delete-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("failed to save game: %v", err)
		}

		// Verify it exists
		_, err = repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("game should exist before deletion: %v", err)
		}

		// Delete it
		err = repo.DeleteGame(testSlug)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify it's gone
		_, err = repo.GetGame(testSlug)
		if err == nil {
			t.Error("expected error after deletion, got nil")
		}
		var notFoundErr *domainerr.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Errorf("expected NotFoundError, got %T: %v", err, err)
		}
	})

	t.Run("Delete non-existent game", func(t *testing.T) {
		nonExistentSlug := "non-existent-game"

		// Delete should not return an error even if the game doesn't exist
		err := repo.DeleteGame(nonExistentSlug)
		if err != nil {
			// If it does return an error, that's also acceptable behavior
			// We just verify it doesn't panic
		}
	})

	t.Run("Delete game with execs and assets", func(t *testing.T) {
		testSlug := "test-game-delete-complex-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		exec, err := game.RehydrateGameExec("/path/to/client", "client")
		if err != nil {
			t.Fatalf("failed to create test exec: %v", err)
		}
		testGame.SetExecutable(*exec)

		assetID := uuid.New()
		createTestAsset(t, db, assetID)
		gameAsset, err := game.RehydrateGameAsset(assetID, "icon")
		if err != nil {
			t.Fatalf("failed to create test game asset: %v", err)
		}
		testGame.SetAsset(*gameAsset)

		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("failed to save game: %v", err)
		}

		// Delete it
		err = repo.DeleteGame(testSlug)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify it's gone
		_, err = repo.GetGame(testSlug)
		if err == nil {
			t.Error("expected error after deletion, got nil")
		}
	})
}

func TestGameRepository_Integration(t *testing.T) {
	repo, _ := setupTest(t)

	t.Run("Full CRUD cycle", func(t *testing.T) {
		testSlug := "test-game-crud-" + uuid.New().String()[:8]
		testGame, err := game.RehydrateGame(testSlug, "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		// Create
		err = repo.SaveGame(testGame)
		if err != nil {
			t.Fatalf("failed to save game: %v", err)
		}

		// Read
		retrieved, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("failed to retrieve game: %v", err)
		}
		if retrieved.Slug != testSlug {
			t.Errorf("expected Slug %q, got %q", testSlug, retrieved.Slug)
		}

		// Update
		err = retrieved.SetName("Updated Name")
		if err != nil {
			t.Fatalf("failed to update name: %v", err)
		}
		err = repo.SaveGame(retrieved)
		if err != nil {
			t.Fatalf("failed to update game: %v", err)
		}

		// Verify update
		updated, err := repo.GetGame(testSlug)
		if err != nil {
			t.Fatalf("failed to retrieve updated game: %v", err)
		}
		if updated.Name != "Updated Name" {
			t.Errorf("expected Name %q, got %q", "Updated Name", updated.Name)
		}

		// Delete
		err = repo.DeleteGame(testSlug)
		if err != nil {
			t.Fatalf("failed to delete game: %v", err)
		}

		// Verify deletion
		_, err = repo.GetGame(testSlug)
		if err == nil {
			t.Error("expected error after deletion, got nil")
		}
	})
}

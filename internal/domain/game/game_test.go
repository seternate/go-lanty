package game

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

func TestNewGame(t *testing.T) {
	t.Run("Valid game", func(t *testing.T) {
		game, err := NewGame("test-game", "Test Game")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if game == nil {
			t.Fatal("expected game, got nil")
		}
		if game.Slug != "test-game" {
			t.Errorf("expected Slug %q, got %q", "test-game", game.Slug)
		}
		if game.Name != "Test Game" {
			t.Errorf("expected Name %q, got %q", "Test Game", game.Name)
		}
		if game.Execs == nil {
			t.Error("expected Execs to be initialized")
		}
		if len(game.Execs) != 0 {
			t.Errorf("expected empty Execs, got %d", len(game.Execs))
		}
		if game.Assets == nil {
			t.Error("expected Assets to be initialized")
		}
		if len(game.Assets) != 0 {
			t.Errorf("expected empty Assets, got %d", len(game.Assets))
		}
	})

	t.Run("Empty slug", func(t *testing.T) {
		game, err := NewGame("", "Test Game")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if game != nil {
			t.Error("expected nil game")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid slug format", func(t *testing.T) {
		invalidSlugs := []string{
			"Test-Game",  // uppercase
			"test_game",  // underscore
			"test game",  // space
			"test-game-", // trailing dash
			"-test-game", // leading dash
			"test--game", // double dash
		}
		for _, slug := range invalidSlugs {
			t.Run(slug, func(t *testing.T) {
				game, err := NewGame(slug, "Test Game")
				if err == nil {
					t.Error("expected error, got nil")
				}
				if game != nil {
					t.Error("expected nil game")
				}
			})
		}
	})

	t.Run("Valid slug formats", func(t *testing.T) {
		validSlugs := []string{
			"test",
			"test-game",
			"test123",
			"123test",
			"test-game-123",
			"a",
			"a-b-c-d",
		}
		for _, slug := range validSlugs {
			t.Run(slug, func(t *testing.T) {
				game, err := NewGame(slug, "Test Game")
				if err != nil {
					t.Errorf("unexpected error for slug %q: %v", slug, err)
				}
				if game == nil {
					t.Errorf("expected game for slug %q, got nil", slug)
				}
			})
		}
	})

	t.Run("Empty name", func(t *testing.T) {
		game, err := NewGame("test-game", "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if game != nil {
			t.Error("expected nil game")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})
}

func TestRehydrateGame(t *testing.T) {
	t.Run("Valid game", func(t *testing.T) {
		game, err := RehydrateGame("test-game", "Test Game")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if game == nil {
			t.Fatal("expected game, got nil")
		}
		if game.Slug != "test-game" {
			t.Errorf("expected Slug %q, got %q", "test-game", game.Slug)
		}
		if game.Name != "Test Game" {
			t.Errorf("expected Name %q, got %q", "Test Game", game.Name)
		}
	})

	t.Run("Invalid slug returns TrustedInvariantViolationError", func(t *testing.T) {
		game, err := RehydrateGame("", "Test Game")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if game != nil {
			t.Error("expected nil game")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})

	t.Run("Empty name returns TrustedInvariantViolationError", func(t *testing.T) {
		game, err := RehydrateGame("test-game", "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if game != nil {
			t.Error("expected nil game")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})
}

func TestGame_SetName(t *testing.T) {
	t.Run("Valid name", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		err := game.SetName("New Name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if game.Name != "New Name" {
			t.Errorf("expected Name %q, got %q", "New Name", game.Name)
		}
	})

	t.Run("Empty name", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		err := game.SetName("")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
		// Name should not be changed
		if game.Name != "Test Game" {
			t.Errorf("expected Name to remain %q, got %q", "Test Game", game.Name)
		}
	})
}

func TestGame_SetExecutable(t *testing.T) {
	t.Run("Set client executable", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		exec, _ := NewGameExec("/path/to/client", "client")
		game.SetExecutable(*exec)
		if len(game.Execs) != 1 {
			t.Errorf("expected 1 exec, got %d", len(game.Execs))
		}
		if game.Execs[GAME_EXEC_ROLE_CLIENT].Path != "/path/to/client" {
			t.Errorf("expected client exec path %q, got %q", "/path/to/client", game.Execs[GAME_EXEC_ROLE_CLIENT].Path)
		}
	})

	t.Run("Set server executable", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		exec, _ := NewGameExec("/path/to/server", "server")
		game.SetExecutable(*exec)
		if len(game.Execs) != 1 {
			t.Errorf("expected 1 exec, got %d", len(game.Execs))
		}
		if game.Execs[GAME_EXEC_ROLE_SERVER].Path != "/path/to/server" {
			t.Errorf("expected server exec path %q, got %q", "/path/to/server", game.Execs[GAME_EXEC_ROLE_SERVER].Path)
		}
	})

	t.Run("Replace existing executable", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		exec1, _ := NewGameExec("/path/to/client1", "client")
		game.SetExecutable(*exec1)
		exec2, _ := NewGameExec("/path/to/client2", "client")
		game.SetExecutable(*exec2)
		if len(game.Execs) != 1 {
			t.Errorf("expected 1 exec, got %d", len(game.Execs))
		}
		if game.Execs[GAME_EXEC_ROLE_CLIENT].Path != "/path/to/client2" {
			t.Errorf("expected client exec path %q, got %q", "/path/to/client2", game.Execs[GAME_EXEC_ROLE_CLIENT].Path)
		}
	})

	t.Run("Set multiple executables", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		clientExec, _ := NewGameExec("/path/to/client", "client")
		serverExec, _ := NewGameExec("/path/to/server", "server")
		game.SetExecutable(*clientExec)
		game.SetExecutable(*serverExec)
		if len(game.Execs) != 2 {
			t.Errorf("expected 2 execs, got %d", len(game.Execs))
		}
		if game.Execs[GAME_EXEC_ROLE_CLIENT].Path != "/path/to/client" {
			t.Errorf("expected client exec path %q, got %q", "/path/to/client", game.Execs[GAME_EXEC_ROLE_CLIENT].Path)
		}
		if game.Execs[GAME_EXEC_ROLE_SERVER].Path != "/path/to/server" {
			t.Errorf("expected server exec path %q, got %q", "/path/to/server", game.Execs[GAME_EXEC_ROLE_SERVER].Path)
		}
	})
}

func TestGame_SetAsset(t *testing.T) {
	testID := uuid.New()

	t.Run("Set icon asset", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		asset, _ := NewGameAsset(testID, "icon", "image/png")
		oldAsset := game.SetAsset(*asset)
		if oldAsset != nil {
			t.Error("expected nil old asset")
		}
		if len(game.Assets) != 1 {
			t.Errorf("expected 1 asset, got %d", len(game.Assets))
		}
		if game.Assets[GAME_ASSET_ROLE_ICON].AssetID != testID {
			t.Errorf("expected icon asset ID %v, got %v", testID, game.Assets[GAME_ASSET_ROLE_ICON].AssetID)
		}
	})

	t.Run("Set blob asset", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		asset, _ := NewGameAsset(testID, "blob", "application/octet-stream")
		oldAsset := game.SetAsset(*asset)
		if oldAsset != nil {
			t.Error("expected nil old asset")
		}
		if len(game.Assets) != 1 {
			t.Errorf("expected 1 asset, got %d", len(game.Assets))
		}
		if game.Assets[GAME_ASSET_ROLE_BLOB].AssetID != testID {
			t.Errorf("expected blob asset ID %v, got %v", testID, game.Assets[GAME_ASSET_ROLE_BLOB].AssetID)
		}
	})

	t.Run("Replace existing asset", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		testID1 := uuid.New()
		testID2 := uuid.New()
		asset1, _ := NewGameAsset(testID1, "icon", "image/png")
		asset2, _ := NewGameAsset(testID2, "icon", "image/jpeg")
		game.SetAsset(*asset1)
		oldAsset := game.SetAsset(*asset2)
		if oldAsset == nil {
			t.Fatal("expected non-nil old asset")
		}
		if oldAsset.AssetID != testID1 {
			t.Errorf("expected old asset ID %v, got %v", testID1, oldAsset.AssetID)
		}
		if len(game.Assets) != 1 {
			t.Errorf("expected 1 asset, got %d", len(game.Assets))
		}
		if game.Assets[GAME_ASSET_ROLE_ICON].AssetID != testID2 {
			t.Errorf("expected icon asset ID %v, got %v", testID2, game.Assets[GAME_ASSET_ROLE_ICON].AssetID)
		}
	})

	t.Run("Set multiple assets", func(t *testing.T) {
		game, _ := NewGame("test-game", "Test Game")
		iconID := uuid.New()
		blobID := uuid.New()
		iconAsset, _ := NewGameAsset(iconID, "icon", "image/png")
		blobAsset, _ := NewGameAsset(blobID, "blob", "application/octet-stream")
		game.SetAsset(*iconAsset)
		game.SetAsset(*blobAsset)
		if len(game.Assets) != 2 {
			t.Errorf("expected 2 assets, got %d", len(game.Assets))
		}
		if game.Assets[GAME_ASSET_ROLE_ICON].AssetID != iconID {
			t.Errorf("expected icon asset ID %v, got %v", iconID, game.Assets[GAME_ASSET_ROLE_ICON].AssetID)
		}
		if game.Assets[GAME_ASSET_ROLE_BLOB].AssetID != blobID {
			t.Errorf("expected blob asset ID %v, got %v", blobID, game.Assets[GAME_ASSET_ROLE_BLOB].AssetID)
		}
	})
}

func TestValidateSlug(t *testing.T) {
	t.Run("Valid slugs", func(t *testing.T) {
		validSlugs := []string{
			"test",
			"test-game",
			"test123",
			"123test",
			"test-game-123",
			"a",
			"a-b-c-d",
			"my-game-2024",
		}
		for _, slug := range validSlugs {
			t.Run(slug, func(t *testing.T) {
				err := validateSlug(slug)
				if err != nil {
					t.Errorf("unexpected error for slug %q: %v", slug, err)
				}
			})
		}
	})

	t.Run("Invalid slugs", func(t *testing.T) {
		invalidSlugs := []string{
			"",
			"Test-Game",  // uppercase
			"test_game",  // underscore
			"test game",  // space
			"test-game-", // trailing dash
			"-test-game", // leading dash
			"test--game", // double dash
			"test.game",  // dot
			"test@game",  // special char
		}
		for _, slug := range invalidSlugs {
			t.Run(slug, func(t *testing.T) {
				err := validateSlug(slug)
				if err == nil {
					t.Errorf("expected error for slug %q, got nil", slug)
					return
				}
				var validationErrs *domainerr.ValidationErrors
				if !errors.As(err, &validationErrs) {
					t.Errorf("expected ValidationErrors, got %T", err)
				}
			})
		}
	})
}

func TestValidateName(t *testing.T) {
	t.Run("Valid names", func(t *testing.T) {
		validNames := []string{
			"Test Game",
			"My Awesome Game",
			"Game 123",
			"a",
			"Very Long Game Name With Many Words",
		}
		for _, name := range validNames {
			t.Run(name, func(t *testing.T) {
				err := validateName(name)
				if err != nil {
					t.Errorf("unexpected error for name %q: %v", name, err)
				}
			})
		}
	})

	t.Run("Empty name", func(t *testing.T) {
		err := validateName("")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil && validationErr.Field != "name" {
			t.Errorf("expected Field %q, got %q", "name", validationErr.Field)
		}
	})
}

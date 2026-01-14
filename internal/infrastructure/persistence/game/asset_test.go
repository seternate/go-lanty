package game

import (
	"testing"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/internal/domain/game"
)

func TestDisassembleGameAsset(t *testing.T) {
	t.Run("Valid game asset with icon role", func(t *testing.T) {
		testAssetID := uuid.New()
		gameAsset, err := game.RehydrateGameAsset(testAssetID, "icon")
		if err != nil {
			t.Fatalf("failed to create test game asset: %v", err)
		}

		row := disassembleGameAsset("test-game", *gameAsset)
		if row == nil {
			t.Fatal("expected row, got nil")
		}
		if row.GameSlug != "test-game" {
			t.Errorf("expected GameSlug %q, got %q", "test-game", row.GameSlug)
		}
		if row.AssetID != testAssetID {
			t.Errorf("expected AssetID %v, got %v", testAssetID, row.AssetID)
		}
		if row.Role != "icon" {
			t.Errorf("expected Role %q, got %q", "icon", row.Role)
		}
	})

	t.Run("Valid game asset with blob role", func(t *testing.T) {
		testAssetID := uuid.New()
		gameAsset, err := game.RehydrateGameAsset(testAssetID, "blob")
		if err != nil {
			t.Fatalf("failed to create test game asset: %v", err)
		}

		row := disassembleGameAsset("test-game", *gameAsset)
		if row == nil {
			t.Fatal("expected row, got nil")
		}
		if row.Role != "blob" {
			t.Errorf("expected Role %q, got %q", "blob", row.Role)
		}
	})

	t.Run("Different game slugs", func(t *testing.T) {
		testAssetID := uuid.New()
		gameAsset, err := game.RehydrateGameAsset(testAssetID, "icon")
		if err != nil {
			t.Fatalf("failed to create test game asset: %v", err)
		}

		row1 := disassembleGameAsset("game-1", *gameAsset)
		row2 := disassembleGameAsset("game-2", *gameAsset)

		if row1.GameSlug != "game-1" {
			t.Errorf("expected GameSlug %q, got %q", "game-1", row1.GameSlug)
		}
		if row2.GameSlug != "game-2" {
			t.Errorf("expected GameSlug %q, got %q", "game-2", row2.GameSlug)
		}
		if row1.AssetID != row2.AssetID {
			t.Error("expected same AssetID for both rows")
		}
	})
}

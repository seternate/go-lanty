package game

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/internal/domain/game"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

func TestGameRow_Assemble(t *testing.T) {
	t.Run("Valid game row with no execs or assets", func(t *testing.T) {
		row := GameRow{
			Slug: "test-game",
			Name: "Test Game",
		}

		assembled, err := row.Assemble(nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled == nil {
			t.Fatal("expected game, got nil")
		}
		if assembled.Slug != "test-game" {
			t.Errorf("expected Slug %q, got %q", "test-game", assembled.Slug)
		}
		if assembled.Name != "Test Game" {
			t.Errorf("expected Name %q, got %q", "Test Game", assembled.Name)
		}
		if len(assembled.Execs) != 0 {
			t.Errorf("expected 0 execs, got %d", len(assembled.Execs))
		}
		if len(assembled.Assets) != 0 {
			t.Errorf("expected 0 assets, got %d", len(assembled.Assets))
		}
	})

	t.Run("Valid game row with execs", func(t *testing.T) {
		row := GameRow{
			Slug: "test-game",
			Name: "Test Game",
		}

		execID := uuid.New()
		execRow := GameExecRow{
			ID:            execID,
			GameSlug:      "test-game",
			Role:          "client",
			Path:          "/path/to/client",
			RequiresAdmin: nil,
			Format:        nil,
			ArgSeperator:  nil,
		}

		assembled, err := row.Assemble([]GameExecRow{execRow}, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(assembled.Execs) != 1 {
			t.Errorf("expected 1 exec, got %d", len(assembled.Execs))
		}
		if exec, ok := assembled.Execs[game.GAME_EXEC_ROLE_CLIENT]; !ok {
			t.Error("expected client exec to be present")
		} else if exec.Path != "/path/to/client" {
			t.Errorf("expected Path %q, got %q", "/path/to/client", exec.Path)
		}
	})

	t.Run("Valid game row with assets", func(t *testing.T) {
		row := GameRow{
			Slug: "test-game",
			Name: "Test Game",
		}

		assetID := uuid.New()
		assetRow := GameAssetRow{
			GameSlug: "test-game",
			AssetID:  assetID,
			Role:     "icon",
		}

		assembled, err := row.Assemble(nil, nil, []GameAssetRow{assetRow})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(assembled.Assets) != 1 {
			t.Errorf("expected 1 asset, got %d", len(assembled.Assets))
		}
		if asset, ok := assembled.Assets[game.GAME_ASSET_ROLE_ICON]; !ok {
			t.Error("expected icon asset to be present")
		} else if asset.AssetID != assetID {
			t.Errorf("expected AssetID %v, got %v", assetID, asset.AssetID)
		}
	})

	t.Run("Valid game row with execs and args", func(t *testing.T) {
		row := GameRow{
			Slug: "test-game",
			Name: "Test Game",
		}

		execID := uuid.New()
		execRow := GameExecRow{
			ID:            execID,
			GameSlug:      "test-game",
			Role:          "client",
			Path:          "/path/to/client",
			RequiresAdmin: nil,
			Format:        nil,
			ArgSeperator:  nil,
		}

		argRow := GameArgRow{
			GameExecID:     execID,
			Role:           "flag",
			Name:           "test-arg",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--test",
			Description:    nil,
			DefaultString:  nil,
			DefaultBool:    nil,
			DefaultInt:     nil,
			DefaultFloat:   nil,
			Enums:          nil,
			MinInt:         nil,
			MaxInt:         nil,
			MinFloat:       nil,
			MaxFloat:       nil,
			FloatPrecision: nil,
			OrderIndex:     0,
		}

		assembled, err := row.Assemble([]GameExecRow{execRow}, []GameArgRow{argRow}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec, ok := assembled.Execs[game.GAME_EXEC_ROLE_CLIENT]; ok {
			if len(exec.Args) != 1 {
				t.Errorf("expected 1 arg, got %d", len(exec.Args))
			}
		} else {
			t.Error("expected client exec to be present")
		}
	})

	t.Run("Invalid slug", func(t *testing.T) {
		row := GameRow{
			Slug: "",
			Name: "Test Game",
		}

		assembled, err := row.Assemble(nil, nil, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil game")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		row := GameRow{
			Slug: "test-game",
			Name: "",
		}

		assembled, err := row.Assemble(nil, nil, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil game")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})
}

func TestDisassemble(t *testing.T) {
	t.Run("Valid game with no execs or assets", func(t *testing.T) {
		game, err := game.RehydrateGame("test-game", "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		disassembled := Disassemble(game)
		if disassembled == nil {
			t.Fatal("expected disassembled game, got nil")
		}
		if disassembled.GameRow.Slug != "test-game" {
			t.Errorf("expected Slug %q, got %q", "test-game", disassembled.GameRow.Slug)
		}
		if disassembled.GameRow.Name != "Test Game" {
			t.Errorf("expected Name %q, got %q", "Test Game", disassembled.GameRow.Name)
		}
		if len(disassembled.ExecRows) != 0 {
			t.Errorf("expected 0 exec rows, got %d", len(disassembled.ExecRows))
		}
		if len(disassembled.ArgRows) != 0 {
			t.Errorf("expected 0 arg rows, got %d", len(disassembled.ArgRows))
		}
		if len(disassembled.AssetRows) != 0 {
			t.Errorf("expected 0 asset rows, got %d", len(disassembled.AssetRows))
		}
	})

	t.Run("Valid game with execs", func(t *testing.T) {
		testGame, err := game.RehydrateGame("test-game", "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		exec, err := game.RehydrateGameExec("/path/to/client", "client")
		if err != nil {
			t.Fatalf("failed to create test exec: %v", err)
		}
		testGame.SetExecutable(*exec)

		disassembled := Disassemble(testGame)
		if len(disassembled.ExecRows) != 1 {
			t.Errorf("expected 1 exec row, got %d", len(disassembled.ExecRows))
		}
		if disassembled.ExecRows[0].Role != "client" {
			t.Errorf("expected Role %q, got %q", "client", disassembled.ExecRows[0].Role)
		}
		if disassembled.ExecRows[0].Path != "/path/to/client" {
			t.Errorf("expected Path %q, got %q", "/path/to/client", disassembled.ExecRows[0].Path)
		}
	})

	t.Run("Valid game with assets", func(t *testing.T) {
		testGame, err := game.RehydrateGame("test-game", "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		assetID := uuid.New()
		asset, err := game.RehydrateGameAsset(assetID, "icon")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}
		testGame.SetAsset(*asset)

		disassembled := Disassemble(testGame)
		if len(disassembled.AssetRows) != 1 {
			t.Errorf("expected 1 asset row, got %d", len(disassembled.AssetRows))
		}
		if disassembled.AssetRows[0].Role != "icon" {
			t.Errorf("expected Role %q, got %q", "icon", disassembled.AssetRows[0].Role)
		}
		if disassembled.AssetRows[0].AssetID != assetID {
			t.Errorf("expected AssetID %v, got %v", assetID, disassembled.AssetRows[0].AssetID)
		}
	})

	t.Run("Round trip: Disassemble then Assemble", func(t *testing.T) {
		originalGame, err := game.RehydrateGame("test-game", "Test Game")
		if err != nil {
			t.Fatalf("failed to create test game: %v", err)
		}

		exec, err := game.RehydrateGameExec("/path/to/client", "client")
		if err != nil {
			t.Fatalf("failed to create test exec: %v", err)
		}
		originalGame.SetExecutable(*exec)

		assetID := uuid.New()
		asset, err := game.RehydrateGameAsset(assetID, "icon")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}
		originalGame.SetAsset(*asset)

		disassembled := Disassemble(originalGame)
		reassembled, err := disassembled.GameRow.Assemble(disassembled.ExecRows, disassembled.ArgRows, disassembled.AssetRows)
		if err != nil {
			t.Fatalf("failed to reassemble: %v", err)
		}

		if reassembled.Slug != originalGame.Slug {
			t.Errorf("expected Slug %q, got %q", originalGame.Slug, reassembled.Slug)
		}
		if reassembled.Name != originalGame.Name {
			t.Errorf("expected Name %q, got %q", originalGame.Name, reassembled.Name)
		}
		if len(reassembled.Execs) != len(originalGame.Execs) {
			t.Errorf("expected %d execs, got %d", len(originalGame.Execs), len(reassembled.Execs))
		}
		if len(reassembled.Assets) != len(originalGame.Assets) {
			t.Errorf("expected %d assets, got %d", len(originalGame.Assets), len(reassembled.Assets))
		}
	})
}

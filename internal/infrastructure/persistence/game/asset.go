package game

import (
	"github.com/google/uuid"
	"github.com/seternate/go-lanty/internal/domain/game"
)

type GameAssetRow struct {
	GameSlug string    `db:"game_slug"`
	AssetID  uuid.UUID `db:"asset_id"`
	Role     string    `db:"role"`
}

func disassembleGameAsset(gameSlug string, asset game.GameAsset) *GameAssetRow {
	return &GameAssetRow{
		GameSlug: gameSlug,
		AssetID:  asset.AssetID,
		Role:     asset.Role.String(),
	}
}

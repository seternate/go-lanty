package persistence

import (
	"github.com/jmoiron/sqlx"
	"github.com/seternate/go-lanty/pkg/domain/asset"
	"github.com/seternate/go-lanty/pkg/domain/game"
	assetrepo "github.com/seternate/go-lanty/pkg/infrastructure/persistence/asset"
	gamerepo "github.com/seternate/go-lanty/pkg/infrastructure/persistence/game"
)

type Repositories struct {
	Game  game.GameRepository
	Asset asset.AssetRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Game:  gamerepo.NewGameRepository(db),
		Asset: assetrepo.NewAssetRepository(db),
	}
}

package persistence

import (
	"github.com/jmoiron/sqlx"
	"github.com/seternate/go-lanty/internal/domain/asset"
	"github.com/seternate/go-lanty/internal/domain/game"
	"github.com/seternate/go-lanty/internal/domain/user"
	assetrepo "github.com/seternate/go-lanty/internal/infrastructure/persistence/asset"
	gamerepo "github.com/seternate/go-lanty/internal/infrastructure/persistence/game"
	userrepo "github.com/seternate/go-lanty/internal/infrastructure/persistence/user"
)

type Repositories struct {
	Game  game.GameRepository
	Asset asset.AssetRepository
	User  user.UserRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Game:  gamerepo.NewGameRepository(db),
		Asset: assetrepo.NewAssetRepository(db),
		User:  userrepo.NewUserRepository(db),
	}
}

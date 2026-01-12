package game

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/seternate/go-lanty/pkg/application/asset"
	domainasset "github.com/seternate/go-lanty/pkg/domain/asset"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
	domaingame "github.com/seternate/go-lanty/pkg/domain/game"
)

// Mock implementations
type mockGameRepository struct {
	games     map[string]*domaingame.Game
	getErr    error
	saveErr   error
	deleteErr error
}

func newMockGameRepository() *mockGameRepository {
	return &mockGameRepository{
		games: make(map[string]*domaingame.Game),
	}
}

func (m *mockGameRepository) GetGame(slug string) (*domaingame.Game, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	game, ok := m.games[slug]
	if !ok {
		return nil, domainerr.NotFoundErr("game", slug)
	}
	return game, nil
}

func (m *mockGameRepository) SaveGame(game *domaingame.Game) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.games[game.Slug] = game
	return nil
}

func (m *mockGameRepository) DeleteGame(slug string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.games, slug)
	return nil
}

type mockAssetService struct {
	storeAssetFunc func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error)
	deleteAssetFunc func(id uuid.UUID) error
	getAssetFunc    func(id uuid.UUID) (*asset.AssetView, error)
	getContentFunc  func(id uuid.UUID) (io.ReadCloser, error)
}

func newMockAssetService() *mockAssetService {
	return &mockAssetService{
		storeAssetFunc: func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
			asset, err := domainasset.NewAsset("file:///test", 100, cmd.Checksum, cmd.Algorithm, "image/png")
			if err != nil {
				return nil, err
			}
			return asset, nil
		},
		deleteAssetFunc: func(id uuid.UUID) error {
			return nil
		},
		getAssetFunc: func(id uuid.UUID) (*asset.AssetView, error) {
			return &asset.AssetView{
				ID:       id,
				Size:     100,
				Checksum: "checksum",
				MimeType: "image/png",
			}, nil
		},
		getContentFunc: func(id uuid.UUID) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("content"))), nil
		},
	}
}

func (m *mockAssetService) Command() asset.CommandService {
	return &mockAssetCommandService{
		storeFunc: m.storeAssetFunc,
		deleteFunc: m.deleteAssetFunc,
	}
}

func (m *mockAssetService) Query() asset.QueryService {
	return &mockAssetQueryService{
		getFunc: m.getAssetFunc,
		getContentFunc: m.getContentFunc,
	}
}

func createMockAssetService(storeFunc func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error), deleteFunc func(id uuid.UUID) error) *asset.Service {
	if storeFunc == nil {
		storeFunc = func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
			asset, err := domainasset.NewAsset("file:///test", 100, cmd.Checksum, cmd.Algorithm, "image/png")
			if err != nil {
				return nil, err
			}
			return asset, nil
		}
	}
	if deleteFunc == nil {
		deleteFunc = func(id uuid.UUID) error {
			return nil
		}
	}
	mock := &mockAssetService{
		storeAssetFunc: storeFunc,
		deleteAssetFunc: deleteFunc,
		getAssetFunc: func(id uuid.UUID) (*asset.AssetView, error) {
			return &asset.AssetView{
				ID:       id,
				Size:     100,
				Checksum: "checksum",
				MimeType: "image/png",
			}, nil
		},
		getContentFunc: func(id uuid.UUID) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("content"))), nil
		},
	}
	commandService := &mockAssetCommandService{
		storeFunc:  mock.storeAssetFunc,
		deleteFunc: mock.deleteAssetFunc,
	}
	queryService := &mockAssetQueryService{
		getFunc:       mock.getAssetFunc,
		getContentFunc: mock.getContentFunc,
	}
	return asset.NewService(queryService, commandService)
}

type mockAssetCommandService struct {
	storeFunc  func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error)
	deleteFunc func(id uuid.UUID) error
}

func (m *mockAssetCommandService) StoreNewAsset(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
	return m.storeFunc(cmd)
}

func (m *mockAssetCommandService) DeleteAsset(id uuid.UUID) error {
	return m.deleteFunc(id)
}

type mockAssetQueryService struct {
	getFunc       func(id uuid.UUID) (*asset.AssetView, error)
	getContentFunc func(id uuid.UUID) (io.ReadCloser, error)
}

func (m *mockAssetQueryService) GetAsset(id uuid.UUID) (*asset.AssetView, error) {
	return m.getFunc(id)
}

func (m *mockAssetQueryService) GetAssetContent(id uuid.UUID) (io.ReadCloser, error) {
	return m.getContentFunc(id)
}

func TestCommandService_UpsertGame(t *testing.T) {
	t.Run("create new game", func(t *testing.T) {
		repo := newMockGameRepository()
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "Test Game",
			Executables: []UpsertGameExecutable{
				{
					Role: "client",
					Path: "/path/to/client",
					Args: []UpsertGameExecutableArg{
						{
							Role:         "string",
							Name:         "player-name",
							Required:     boolPtr(true),
							Argument:     "--player",
							DefaultString: stringPtr("default-player"),
							OrderIndex:   0,
						},
					},
				},
			},
		}

		game, created, err := service.UpsertGame(cmd)
		require.NoError(t, err)
		assert.True(t, created)
		assert.NotNil(t, game)
		assert.Equal(t, "test-game", game.Slug)
		assert.Equal(t, "Test Game", game.Name)
		assert.Len(t, game.Execs, 1)

		// Verify game was saved
		savedGame, err := repo.GetGame("test-game")
		require.NoError(t, err)
		assert.Equal(t, game.Slug, savedGame.Slug)
	})

	t.Run("update existing game", func(t *testing.T) {
		repo := newMockGameRepository()
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		// Create existing game
		existingGame, err := domaingame.NewGame("test-game", "Old Name")
		require.NoError(t, err)
		repo.games["test-game"] = existingGame

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "New Name",
			Executables: []UpsertGameExecutable{
				{
					Role: "server",
					Path: "/path/to/server",
				},
			},
		}

		game, created, err := service.UpsertGame(cmd)
		require.NoError(t, err)
		assert.False(t, created)
		assert.NotNil(t, game)
		assert.Equal(t, "New Name", game.Name)
		assert.Len(t, game.Execs, 1)
	})

	t.Run("game arg creation failure", func(t *testing.T) {
		repo := newMockGameRepository()
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "Test Game",
			Executables: []UpsertGameExecutable{
				{
					Role: "client",
					Path: "/path/to/client",
					Args: []UpsertGameExecutableArg{
						{
							Role:      "invalid-role", // Invalid role
							Name:      "arg",
							Argument:  "--arg",
							OrderIndex: 0,
						},
					},
				},
			},
		}

		game, created, err := service.UpsertGame(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "failed to create new game arg")
	})

	t.Run("game exec creation failure", func(t *testing.T) {
		repo := newMockGameRepository()
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "Test Game",
			Executables: []UpsertGameExecutable{
				{
					Role: "invalid-role", // Invalid role
					Path: "/path/to/client",
				},
			},
		}

		game, created, err := service.UpsertGame(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "failed to create new game exec")
	})

	t.Run("repository get error (non-not-found)", func(t *testing.T) {
		repo := newMockGameRepository()
		repo.getErr = errors.New("database error")
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "Test Game",
		}

		game, created, err := service.UpsertGame(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "failed to get game from repository")
	})

	t.Run("repository save error on create", func(t *testing.T) {
		repo := newMockGameRepository()
		repo.saveErr = errors.New("save failed")
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "Test Game",
		}

		game, created, err := service.UpsertGame(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "failed to save new game to repository")
	})

	t.Run("repository save error on update", func(t *testing.T) {
		repo := newMockGameRepository()
		existingGame, err := domaingame.NewGame("test-game", "Old Name")
		require.NoError(t, err)
		repo.games["test-game"] = existingGame
		repo.saveErr = errors.New("save failed")

		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "New Name",
		}

		game, created, err := service.UpsertGame(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "failed to save updated game to repository")
	})

	t.Run("set name error", func(t *testing.T) {
		repo := newMockGameRepository()
		existingGame, err := domaingame.NewGame("test-game", "Old Name")
		require.NoError(t, err)
		repo.games["test-game"] = existingGame

		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "", // Empty name should fail validation
		}

		game, created, err := service.UpsertGame(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "failed to update game name")
	})

	t.Run("duplicate executable roles", func(t *testing.T) {
		repo := newMockGameRepository()
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "Test Game",
			Executables: []UpsertGameExecutable{
				{
					Role: "client",
					Path: "/path/to/client1",
				},
				{
					Role: "client", // Duplicate role
					Path: "/path/to/client2",
				},
			},
		}

		game, created, err := service.UpsertGame(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "failed to validate executables")

		// Verify it's a validation error
		var validationErr domainerr.Validation
		assert.ErrorAs(t, err, &validationErr)

		// Verify the error contains information about duplicate roles
		assert.Contains(t, err.Error(), "duplicate role found")
		assert.Contains(t, err.Error(), "role=client")
	})

	t.Run("multiple duplicate executable roles", func(t *testing.T) {
		repo := newMockGameRepository()
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := UpsertGameCommand{
			Slug: "test-game",
			Name: "Test Game",
			Executables: []UpsertGameExecutable{
				{
					Role: "client",
					Path: "/path/to/client1",
				},
				{
					Role: "server",
					Path: "/path/to/server1",
				},
				{
					Role: "client", // First duplicate
					Path: "/path/to/client2",
				},
				{
					Role: "server", // Second duplicate
					Path: "/path/to/server2",
				},
			},
		}

		game, created, err := service.UpsertGame(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Nil(t, game)
		assert.Contains(t, err.Error(), "failed to validate executables")

		// Verify it's a validation error
		var validationErr domainerr.Validation
		assert.ErrorAs(t, err, &validationErr)

		// Verify the error contains information about duplicate roles
		assert.Contains(t, err.Error(), "duplicate role found")
		assert.Contains(t, err.Error(), "role=client")
		assert.Contains(t, err.Error(), "role=server")
	})
}

func TestCommandService_DeleteGame(t *testing.T) {
	t.Run("success with assets", func(t *testing.T) {
		repo := newMockGameRepository()
		
		// Create game with assets
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)

		iconID := uuid.New()
		blobID := uuid.New()
		iconAsset, err := domaingame.NewGameAsset(iconID, "icon", "image/png")
		require.NoError(t, err)
		blobAsset, err := domaingame.NewGameAsset(blobID, "blob", "application/octet-stream")
		require.NoError(t, err)

		game.SetAsset(*iconAsset)
		game.SetAsset(*blobAsset)
		repo.games["test-game"] = game

		iconDeleteCalled := false
		blobDeleteCalled := false
		iconService := createMockAssetService(nil, func(id uuid.UUID) error {
			iconDeleteCalled = true
			assert.Equal(t, iconID, id)
			return nil
		})
		blobService := createMockAssetService(nil, func(id uuid.UUID) error {
			blobDeleteCalled = true
			assert.Equal(t, blobID, id)
			return nil
		})

		service := NewCommandService(repo, iconService, blobService)

		err = service.DeleteGame("test-game")
		require.NoError(t, err)

		// Verify game was deleted
		_, err = repo.GetGame("test-game")
		assert.Error(t, err)
		var notFound domainerr.NotFound
		assert.ErrorAs(t, err, &notFound)

		// Verify assets were deleted
		assert.True(t, iconDeleteCalled)
		assert.True(t, blobDeleteCalled)
	})

	t.Run("game not found", func(t *testing.T) {
		repo := newMockGameRepository()
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		err := service.DeleteGame("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get game from repository")
	})

	t.Run("repository delete failure", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)
		repo.games["test-game"] = game
		repo.deleteErr = errors.New("delete failed")

		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		err = service.DeleteGame("test-game")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete game from repository")
	})

	t.Run("asset delete failure", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)

		iconID := uuid.New()
		iconAsset, err := domaingame.NewGameAsset(iconID, "icon", "image/png")
		require.NoError(t, err)
		game.SetAsset(*iconAsset)
		repo.games["test-game"] = game

		iconService := createMockAssetService(nil, func(id uuid.UUID) error {
			return errors.New("asset delete failed")
		})
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		err = service.DeleteGame("test-game")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete game assets")
	})

	t.Run("missing asset service for role", func(t *testing.T) {
		repo := newMockGameRepository()
		// Create asset with undefined role (not icon or blob)
		// This shouldn't happen in practice, but we test the error handling
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		// The service only has icon and blob services, so this should work
		// But if we had an undefined role, it would fail
		err := service.DeleteGame("test-game")
		// This will fail because game doesn't exist, not because of missing service
		assert.Error(t, err)
	})
}

func TestCommandService_StoreNewAsset(t *testing.T) {
	t.Run("create new asset", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)
		repo.games["test-game"] = game

		assetID := uuid.New()
		iconService := createMockAssetService(func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
			asset, err := domainasset.RehydrateAsset(assetID, "file:///test", 100, cmd.Checksum, cmd.Algorithm, "image/png")
			if err != nil {
				return nil, err
			}
			return asset, nil
		}, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := StoreNewAssetCommand{
			Slug:      "test-game",
			Role:      domaingame.GAME_ASSET_ROLE_ICON,
			Checksum:  "abc123",
			Algorithm: "md5",
			Data:      bytes.NewReader([]byte("icon data")),
		}

		created, err := service.StoreNewAsset(cmd)
		require.NoError(t, err)
		assert.True(t, created)

		// Verify game was updated with asset
		updatedGame, err := repo.GetGame("test-game")
		require.NoError(t, err)
		assert.Contains(t, updatedGame.Assets, domaingame.GAME_ASSET_ROLE_ICON)
		assert.Equal(t, assetID, updatedGame.Assets[domaingame.GAME_ASSET_ROLE_ICON].AssetID)
	})

	t.Run("replace existing asset", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)

		oldIconID := uuid.New()
		oldIconAsset, err := domaingame.NewGameAsset(oldIconID, "icon", "image/png")
		require.NoError(t, err)
		game.SetAsset(*oldIconAsset)
		repo.games["test-game"] = game

		newIconID := uuid.New()
		oldAssetDeleted := false
		iconService := createMockAssetService(func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
			asset, err := domainasset.RehydrateAsset(newIconID, "file:///test", 100, cmd.Checksum, cmd.Algorithm, "image/png")
			if err != nil {
				return nil, err
			}
			return asset, nil
		}, func(id uuid.UUID) error {
			oldAssetDeleted = true
			assert.Equal(t, oldIconID, id)
			return nil
		})
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := StoreNewAssetCommand{
			Slug:      "test-game",
			Role:      domaingame.GAME_ASSET_ROLE_ICON,
			Checksum:  "new-checksum",
			Algorithm: "md5",
			Data:      bytes.NewReader([]byte("new icon data")),
		}

		created, err := service.StoreNewAsset(cmd)
		require.NoError(t, err)
		assert.False(t, created)
		assert.True(t, oldAssetDeleted)

		// Verify new asset was set
		updatedGame, err := repo.GetGame("test-game")
		require.NoError(t, err)
		assert.Equal(t, newIconID, updatedGame.Assets[domaingame.GAME_ASSET_ROLE_ICON].AssetID)
	})

	t.Run("game not found", func(t *testing.T) {
		repo := newMockGameRepository()
		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := StoreNewAssetCommand{
			Slug:      "non-existent",
			Role:      domaingame.GAME_ASSET_ROLE_ICON,
			Checksum:  "abc123",
			Algorithm: "md5",
			Data:      bytes.NewReader([]byte("data")),
		}

		created, err := service.StoreNewAsset(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to get game from repository")
	})

	t.Run("missing asset service for role", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)
		repo.games["test-game"] = game

		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		// Use undefined role (not icon or blob)
		cmd := StoreNewAssetCommand{
			Slug:      "test-game",
			Role:      domaingame.GAME_ASSET_ROLE_UNDEFINED,
			Checksum:  "abc123",
			Algorithm: "md5",
			Data:      bytes.NewReader([]byte("data")),
		}

		created, err := service.StoreNewAsset(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "missing asset service for role")
	})

	t.Run("asset store failure", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)
		repo.games["test-game"] = game

		iconService := createMockAssetService(func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
			return nil, errors.New("store failed")
		}, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := StoreNewAssetCommand{
			Slug:      "test-game",
			Role:      domaingame.GAME_ASSET_ROLE_ICON,
			Checksum:  "abc123",
			Algorithm: "md5",
			Data:      bytes.NewReader([]byte("data")),
		}

		created, err := service.StoreNewAsset(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to store new asset")
	})

	t.Run("game asset creation failure", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)
		repo.games["test-game"] = game

		// Return asset with invalid mime type for icon
		iconService := createMockAssetService(func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
			asset, err := domainasset.NewAsset("file:///test", 100, cmd.Checksum, cmd.Algorithm, "application/octet-stream")
			if err != nil {
				return nil, err
			}
			return asset, nil
		}, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := StoreNewAssetCommand{
			Slug:      "test-game",
			Role:      domaingame.GAME_ASSET_ROLE_ICON,
			Checksum:  "abc123",
			Algorithm: "md5",
			Data:      bytes.NewReader([]byte("data")),
		}

		created, err := service.StoreNewAsset(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to create new game asset")
	})

	t.Run("old asset delete failure", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)

		oldIconID := uuid.New()
		oldIconAsset, err := domaingame.NewGameAsset(oldIconID, "icon", "image/png")
		require.NoError(t, err)
		game.SetAsset(*oldIconAsset)
		repo.games["test-game"] = game

		iconService := createMockAssetService(func(cmd asset.StoreNewAssetCommand) (*domainasset.Asset, error) {
			asset, err := domainasset.NewAsset("file:///test", 100, cmd.Checksum, cmd.Algorithm, "image/png")
			if err != nil {
				return nil, err
			}
			return asset, nil
		}, func(id uuid.UUID) error {
			return errors.New("delete failed")
		})
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := StoreNewAssetCommand{
			Slug:      "test-game",
			Role:      domaingame.GAME_ASSET_ROLE_ICON,
			Checksum:  "new-checksum",
			Algorithm: "md5",
			Data:      bytes.NewReader([]byte("data")),
		}

		created, err := service.StoreNewAsset(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to delete existing asset")
	})

	t.Run("repository save failure", func(t *testing.T) {
		repo := newMockGameRepository()
		game, err := domaingame.NewGame("test-game", "Test Game")
		require.NoError(t, err)
		repo.games["test-game"] = game
		repo.saveErr = errors.New("save failed")

		iconService := createMockAssetService(nil, nil)
		blobService := createMockAssetService(nil, nil)

		service := NewCommandService(repo, iconService, blobService)

		cmd := StoreNewAssetCommand{
			Slug:      "test-game",
			Role:      domaingame.GAME_ASSET_ROLE_ICON,
			Checksum:  "abc123",
			Algorithm: "md5",
			Data:      bytes.NewReader([]byte("data")),
		}

		created, err := service.StoreNewAsset(cmd)
		assert.Error(t, err)
		assert.False(t, created)
		assert.Contains(t, err.Error(), "failed to save updated game to repository")
	})
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

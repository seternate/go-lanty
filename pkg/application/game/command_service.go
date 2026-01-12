package game

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/application/asset"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
	"github.com/seternate/go-lanty/pkg/domain/game"
)

type commandServiceImpl struct {
	repository    game.GameRepository
	assetServices map[game.GameAssetRole]*asset.Service
}

func NewCommandService(repository game.GameRepository, iconAssetService *asset.Service, blobAssetService *asset.Service) CommandService {
	assetServices := map[game.GameAssetRole]*asset.Service{
		game.GAME_ASSET_ROLE_ICON: iconAssetService,
		game.GAME_ASSET_ROLE_BLOB: blobAssetService,
	}

	return &commandServiceImpl{
		repository:    repository,
		assetServices: assetServices,
	}
}

func (service *commandServiceImpl) UpsertGame(cmd UpsertGameCommand) (*game.Game, bool, error) {
	seenRoles := make(map[string]bool)
	validationErrors := domainerr.ValidationErrs()

	for _, executable := range cmd.Executables {
		if seenRoles[executable.Role] {
			validationErrors.Wrap(domainerr.ValidationErr("executables", "duplicate role found").WithGot("role=%s", executable.Role))
		}
		seenRoles[executable.Role] = true
	}

	if validationErrors.HasErrors() {
		return nil, false, fmt.Errorf("failed to validate executables for slug=%s: %w", cmd.Slug, validationErrors)
	}

	executables := make([]*game.GameExec, 0, len(cmd.Executables))
	for _, executable := range cmd.Executables {
		args := make([]game.GameArg, 0, len(executable.Args))
		for _, arg := range executable.Args {
			gameArg, err := game.NewGameArg(game.GameArgInput{
				Role:           arg.Role,
				Name:           arg.Name,
				Required:       arg.Required,
				Enabled:        arg.Enabled,
				Format:         arg.Format,
				ArgSeparator:   arg.ArgumentSeparator,
				Arg:            arg.Argument,
				Description:    arg.Description,
				DefaultString:  arg.DefaultString,
				DefaultBool:    arg.DefaultBool,
				DefaultInt:     arg.DefaultInt,
				DefaultFloat:   arg.DefaultFloat,
				Enums:          arg.EnumValues,
				MinInt:         arg.MinInt,
				MaxInt:         arg.MaxInt,
				MinFloat:       arg.MinFloat,
				MaxFloat:       arg.MaxFloat,
				FloatPrecision: arg.FloatPrecision,
				OrderIndex:     arg.OrderIndex,
			})
			if err != nil {
				return nil, false, fmt.Errorf("failed to create new game arg for slug=%s: %w", cmd.Slug, err)
			}

			args = append(args, gameArg)
		}

		exec, err := game.NewGameExec(
			executable.Path,
			executable.Role,
			game.WithRequiresAdmin(executable.RequiresAdmin),
			game.WithFormat(executable.Format),
			game.WithArgSeperator(executable.ArgumentSeperator),
			game.WithArgs(args...),
		)
		if err != nil {
			return nil, false, fmt.Errorf("failed to create new game exec for slug=%s: %w", cmd.Slug, err)
		}

		executables = append(executables, exec)
	}

	existingGame, err := service.repository.GetGame(cmd.Slug)
	if err != nil {
		var notFound domainerr.NotFound
		if errors.As(err, &notFound) {
			newGame, err := game.NewGame(cmd.Slug, cmd.Name)
			if err != nil {
				return nil, false, fmt.Errorf("failed to create new game: %w", err)
			}

			for _, executable := range executables {
				newGame.SetExecutable(*executable)
			}

			err = service.repository.SaveGame(newGame)
			if err != nil {
				return nil, false, fmt.Errorf("failed to save new game to repository: %w", err)
			}

			return newGame, true, nil
		}
		return nil, false, fmt.Errorf("failed to get game from repository: %w", err)
	}

	err = existingGame.SetName(cmd.Name)
	if err != nil {
		return nil, false, fmt.Errorf("failed to update game name: %w", err)
	}

	for _, executable := range executables {
		existingGame.SetExecutable(*executable)
	}

	err = service.repository.SaveGame(existingGame)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save updated game to repository: %w", err)
	}

	return existingGame, false, nil
}

func (service *commandServiceImpl) DeleteGame(slug string) error {
	deleteGame, err := service.repository.GetGame(slug)
	if err != nil {
		return fmt.Errorf("failed to get game from repository: %w", err)
	}

	err = service.repository.DeleteGame(deleteGame.Slug)
	if err != nil {
		return fmt.Errorf("failed to delete game from repository: %w", err)
	}

	for _, asset := range deleteGame.Assets {
		assetservice, ok := service.assetServices[asset.Role]
		if !ok {
			err = fmt.Errorf("missing asset service for role=%s: %w", asset.Role.String(), err)
			continue
		}

		err = errors.Join(err, assetservice.Command.DeleteAsset(asset.AssetID))
	}

	if err != nil {
		return fmt.Errorf("failed to delete game assets for slug=%s: %w", deleteGame.Slug, err)
	}

	return nil
}

func (service *commandServiceImpl) StoreNewAsset(cmd StoreNewAssetCommand) (bool, error) {
	updateGame, err := service.repository.GetGame(cmd.Slug)
	if err != nil {
		return false, fmt.Errorf("failed to get game from repository: %w", err)
	}

	assetService, ok := service.assetServices[cmd.Role]
	if !ok {
		return false, fmt.Errorf("missing asset service for role=%s for slug=%s: %w", cmd.Role.String(), cmd.Slug, err)
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return false, fmt.Errorf("failed to generate random asset url for slug=%s: %w", cmd.Slug, err)
	}
	tempURL := "file:///" + uuid.String()
	u, err := url.Parse(tempURL)
	if err != nil {
		return false, fmt.Errorf("failed to parse asset url=%s for slug=%s: %w", tempURL, cmd.Slug, err)
	}
	assetURL := u.String()

	newAsset, err := assetService.Command.StoreNewAsset(asset.StoreNewAssetCommand{
		URL:       assetURL,
		Checksum:  cmd.Checksum,
		Algorithm: cmd.Algorithm,
		MimeType:  "",
		Data:      cmd.Data,
	})
	if err != nil {
		return false, fmt.Errorf("failed to store new asset for slug=%s: %w", cmd.Slug, err)
	}

	defer func() {
		if err != nil && newAsset != nil {
			internalErr := assetService.Command.DeleteAsset(newAsset.ID)
			if internalErr != nil {
				log.Error().Err(internalErr).Msgf("failed to clean up asset (%s) after an error while creating new game for slug=%s", newAsset.ID.String(), cmd.Slug)
			}
		}
	}()

	gameAsset, err := game.NewGameAsset(newAsset.ID, cmd.Role.String(), newAsset.MimeType)
	if err != nil {
		return false, fmt.Errorf("failed to create new game asset for slug=%s: %w", cmd.Slug, err)
	}

	oldAsset := updateGame.SetAsset(*gameAsset)
	created := true
	if oldAsset != nil {
		err = assetService.Command.DeleteAsset(oldAsset.AssetID)
		if err != nil {
			return false, fmt.Errorf("failed to delete existing asset for slug=%s: %w", cmd.Slug, err)
		}
		created = false
	}

	err = service.repository.SaveGame(updateGame)
	if err != nil {
		return false, fmt.Errorf("failed to save updated game to repository: %w", err)
	}

	return created, nil
}

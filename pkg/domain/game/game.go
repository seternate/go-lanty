package game

import (
	"fmt"
	"regexp"

	domainErrors "github.com/seternate/go-lanty/pkg/domain/error"
)

type Game struct {
	Slug   string
	Name   string
	Execs  map[GameExecRole]GameExec
	Assets map[GameAssetRole]GameAsset
}

func NewGame(slug string, name string) (*Game, error) {
	return RehydrateGame(slug, name)
}

func RehydrateGame(slug string, name string) (*Game, error) {
	validationErrors := domainErrors.ValidationErrs().WithMessage("failed to validate for game slug=%s", slug)

	err := validationErrors.Wrap(validateSlug(slug))
	if err != nil {
		return nil, fmt.Errorf("failed to validate for game slug=%s: %w", slug, err)
	}

	err = validationErrors.Wrap(validateName(name))
	if err != nil {
		return nil, fmt.Errorf("failed to validate for game slug=%s: %w", slug, err)
	}

	if len(validationErrors.Errors) > 0 {
		return nil, validationErrors
	}

	return &Game{
		Slug:   slug,
		Name:   name,
		Execs:  make(map[GameExecRole]GameExec),
		Assets: make(map[GameAssetRole]GameAsset),
	}, nil
}

func (game *Game) SetName(name string) error {
	validationErrors := domainErrors.ValidationErrs().WithMessage("failed to validate for game slug=%s", game.Slug)
	err := validationErrors.Wrap(validateName(name))
	if err != nil {
		return fmt.Errorf("failed to validate for game slug=%s: %w", game.Slug, err)
	}

	if len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	game.Name = name
	return nil
}

func (game *Game) SetExecutable(exec GameExec) {
	game.Execs[exec.Role] = exec
}

func (game *Game) SetAsset(asset GameAsset) *GameAsset {
	old, exists := game.Assets[asset.Role]
	game.Assets[asset.Role] = asset
	if exists {
		return &old
	}
	return nil
}

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validateSlug(slug string) error {
	validationErrors := domainErrors.ValidationErrs()

	if len(slug) == 0 {
		validationErrors.Wrap(domainErrors.ValidationErr("slug", "can not be empty"))
	}

	if !slugRegex.MatchString(slug) {
		validationErrors.Wrap(domainErrors.ValidationErr("slug", "does not match expected format").WithExpected("pattern=%s", slugRegex.String()).WithGot(slug))
	}

	if len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	return nil
}

func validateName(name string) error {
	if len(name) == 0 {
		return domainErrors.ValidationErr("name", "can not be empty")
	}

	return nil
}

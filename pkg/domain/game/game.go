package game

import (
	"regexp"

	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

type Game struct {
	Slug   string
	Name   string
	Execs  map[GameExecRole]GameExec
	Assets map[GameAssetRole]GameAsset
}

func NewGame(slug string, name string) (*Game, error) {
	game, err := hydrateGame(slug, name)
	if err != nil {
		return nil, domainerr.InvariantViolationErr("game", slug).WithCause(err)
	}

	return game, nil
}

func RehydrateGame(slug string, name string) (*Game, error) {
	game, err := hydrateGame(slug, name)
	if err != nil {
		return nil, domainerr.TrustedInvariantViolationErr("game", slug).WithCause(err)
	}

	return game, nil
}

func hydrateGame(slug string, name string) (*Game, error) {
	validationErrors := domainerr.ValidationErrs()

	err := validationErrors.Wrap(validateSlug(slug))
	if err != nil {
		return nil, err
	}

	err = validationErrors.Wrap(validateName(name))
	if err != nil {
		return nil, err
	}

	if validationErrors.HasErrors() {
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
	err := validateName(name)
	if err != nil {
		return domainerr.InvariantViolationErr("game", game.Slug).WithCause(err)
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
	validationErrors := domainerr.ValidationErrs()

	if len(slug) == 0 {
		validationErrors.Wrap(domainerr.ValidationErr("slug", "can not be empty"))
	}

	if !slugRegex.MatchString(slug) {
		validationErrors.Wrap(domainerr.ValidationErr("slug", "does not match expected format").WithExpected("pattern=%s", slugRegex.String()).WithGot(slug))
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func validateName(name string) error {
	if len(name) == 0 {
		return domainerr.ValidationErr("name", "can not be empty")
	}

	return nil
}

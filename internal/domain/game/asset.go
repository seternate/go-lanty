package game

import (
	"strings"

	"github.com/google/uuid"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

type GameAssetRole string

const (
	GAME_ASSET_ROLE_UNDEFINED GameAssetRole = ""
	GAME_ASSET_ROLE_ICON      GameAssetRole = "icon"
	GAME_ASSET_ROLE_BLOB      GameAssetRole = "blob"
)

var supportedGameAssetRoles = map[string]GameAssetRole{
	"icon": GAME_ASSET_ROLE_ICON,
	"blob": GAME_ASSET_ROLE_BLOB,
}

func ParseGameAssetRole(role string) (GameAssetRole, error) {
	if gameRole, ok := supportedGameAssetRoles[role]; ok {
		return gameRole, nil
	}
	available := make([]string, 0, len(supportedGameAssetRoles))
	for _, role := range supportedGameAssetRoles {
		available = append(available, role.String())
	}

	return GAME_ASSET_ROLE_UNDEFINED, domainerr.ValidationErr("game asset role", "undefined").WithExpected(strings.Join(available, ", ")).WithGot(role)
}

func (role GameAssetRole) String() string {
	return string(role)
}

type GameAsset struct {
	AssetID uuid.UUID
	Role    GameAssetRole
}

func NewGameAsset(assetID uuid.UUID, role string, mimeType string) (*GameAsset, error) {
	validationErrors := domainerr.ValidationErrs()

	gameAsset, err := hydrateGameAsset(assetID, role)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, err
	}
	if validationErrors.HasErrors() {
		return nil, domainerr.InvariantViolationErr("game asset", assetID.String()).WithCause(validationErrors)
	}

	err = validationErrors.Wrap(validateMimeTypeForRole(gameAsset.Role, mimeType))
	if err != nil {
		return nil, err
	}

	if validationErrors.HasErrors() {
		return nil, domainerr.InvariantViolationErr("game asset", assetID.String()).WithCause(validationErrors)
	}

	return gameAsset, nil
}

func RehydrateGameAsset(assetID uuid.UUID, role string) (*GameAsset, error) {
	gameAsset, err := hydrateGameAsset(assetID, role)
	if err != nil {
		return nil, domainerr.TrustedInvariantViolationErr("game asset", assetID.String()).WithCause(err)
	}
	return gameAsset, nil
}

func hydrateGameAsset(assetID uuid.UUID, role string) (*GameAsset, error) {
	validationErrors := domainerr.ValidationErrs()

	parsedRole, err := ParseGameAssetRole(role)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, err
	}

	if validationErrors.HasErrors() {
		return nil, validationErrors
	}

	return &GameAsset{
		AssetID: assetID,
		Role:    parsedRole,
	}, nil
}

func validateMimeTypeForRole(role GameAssetRole, mimeType string) error {
	if role == GAME_ASSET_ROLE_ICON {
		if !strings.HasPrefix(mimeType, "image/") {
			return domainerr.ValidationErr("mime type", "wrong for role=%s", role.String()).WithExpected("image/*").WithGot(mimeType)
		}
	}
	return nil
}

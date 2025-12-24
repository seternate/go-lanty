package game

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	domainErrors "github.com/seternate/go-lanty/pkg/domain/error"
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

	return GAME_ASSET_ROLE_UNDEFINED, domainErrors.ValidationErr("game asset role", "undefined").WithExpected(strings.Join(available, ", ")).WithGot(role)
}

func (role GameAssetRole) String() string {
	return string(role)
}

type GameAsset struct {
	AssetID uuid.UUID
	Role    GameAssetRole
}

func NewGameAsset(assetID uuid.UUID, role string, mimeType string) (*GameAsset, error) {
	validationErrors := domainErrors.ValidationErrs().WithMessage("failed to validate for asset id=%s role=%s", assetID.String(), role)

	parsedRole, err := ParseGameAssetRole(role)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, fmt.Errorf("failed to parse for asset id=%s role=%s: %w", assetID.String(), role, err)
	}

	err = validationErrors.Wrap(validateMimeTypeForRole(parsedRole, mimeType))
	if err != nil {
		return nil, fmt.Errorf("failed to validate for asset id=%s role=%s: %w", assetID.String(), role, err)
	}

	gameAsset, err := RehydrateGameAsset(assetID, role)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, fmt.Errorf("failed to rehydrate for asset id=%s role=%s: %w", assetID.String(), role, err)
	}

	if len(validationErrors.Errors) > 0 {
		return nil, validationErrors
	}

	return gameAsset, nil
}

func RehydrateGameAsset(assetID uuid.UUID, role string) (*GameAsset, error) {
	validationErrors := domainErrors.ValidationErrs().WithMessage("failed to validate for asset id=%s role=%s", assetID.String(), role)

	parsedRole, err := ParseGameAssetRole(role)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, fmt.Errorf("failed to parse for asset id=%s role=%s: %w", assetID.String(), role, err)
	}

	if len(validationErrors.Errors) > 0 {
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
			return domainErrors.ValidationErr("mime type", "wrong for role=%s", role.String()).WithExpected("image/*").WithGot(mimeType)
		}
	}
	return nil
}

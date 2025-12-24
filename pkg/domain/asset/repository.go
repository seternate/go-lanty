package asset

import (
	"github.com/google/uuid"
)

type AssetRepository interface {
	GetAsset(uuid.UUID) (*Asset, error)
	CreateAsset(*Asset) error
	DeleteAsset(uuid.UUID) error
}

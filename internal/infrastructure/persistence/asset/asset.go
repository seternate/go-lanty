package asset

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/internal/domain/asset"
)

type AssetRow struct {
	ID        uuid.UUID `db:"id"`
	URL       string    `db:"url"`
	Size      uint64    `db:"size"`
	Checksum  string    `db:"checksum"`
	Algorithm string    `db:"algorithm"`
	MimeType  string    `db:"mime_type"`
}

func (row AssetRow) Assemble() (*asset.Asset, error) {
	asset, err := asset.RehydrateAsset(row.ID, row.URL, row.Size, row.Checksum, row.Algorithm, row.MimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to rehydrate asset (database corruption): %w", err)
	}
	return asset, nil
}

func Disassemble(asset *asset.Asset) *AssetRow {
	return &AssetRow{
		ID:        asset.ID,
		URL:       asset.URL.String(),
		Size:      asset.Size,
		Checksum:  asset.Checksum,
		Algorithm: asset.Algorithm.String(),
		MimeType:  asset.MimeType,
	}
}

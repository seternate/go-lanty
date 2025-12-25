package asset

import (
	"net/url"

	"github.com/google/uuid"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

type Asset struct {
	ID        uuid.UUID
	URL       url.URL
	Size      uint64
	Checksum  string
	Algorithm ChecksumAlgorithm
	MimeType  string
}

func NewAsset(assetURL string, size uint64, checksum string, algorithm string, mimeType string) (*Asset, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, domainerr.InternalErr("failed to generate asset ID").WithCause(err)
	}

	asset, err := hydrateAsset(id, assetURL, size, checksum, algorithm, mimeType)
	if err != nil {
		return nil, domainerr.InvariantViolationErr("asset", assetURL).WithCause(err)
	}

	return asset, nil
}

func RehydrateAsset(id uuid.UUID, assetURL string, size uint64, checksum string, algorithm string, mimeType string) (*Asset, error) {
	asset, err := hydrateAsset(id, assetURL, size, checksum, algorithm, mimeType)
	if err != nil {
		return nil, domainerr.TrustedInvariantViolationErr("asset", id.String()).WithCause(err)
	}
	return asset, nil
}

func (asset *Asset) CompareChecksum(checksum string) error {
	if checksum != asset.Checksum {
		return domainerr.ValidationErr("checksum", "mismatch: asset id=%s", asset.ID).WithExpected(asset.Checksum).WithGot(checksum)
	}
	return nil
}

func hydrateAsset(id uuid.UUID, assetURL string, size uint64, checksum string, algorithm string, mimeType string) (*Asset, error) {
	validationErrors := domainerr.ValidationErrs()

	u, err := url.Parse(assetURL)
	err = validationErrors.Wrap(domainerr.ValidationErr("asset URL", "failed to parse URL").WithGot(assetURL).WithCause(err))
	if err != nil {
		return nil, err
	}

	parsedAlgo, err := ParseChecksumAlgorithm(algorithm)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, err
	}

	err = validationErrors.Wrap(validateChecksum(checksum))
	if err != nil {
		return nil, err
	}

	err = validationErrors.Wrap(validateMimeType(mimeType))
	if err != nil {
		return nil, err
	}

	err = validationErrors.Wrap(validateSize(size))
	if err != nil {
		return nil, err
	}

	if validationErrors.HasErrors() {
		return nil, validationErrors
	}

	return &Asset{
		ID:        id,
		URL:       *u,
		Size:      size,
		Checksum:  checksum,
		Algorithm: parsedAlgo,
		MimeType:  mimeType,
	}, nil
}

func validateChecksum(checksum string) error {
	if len(checksum) == 0 {
		return domainerr.ValidationErr("checksum", "can not be empty")
	}
	return nil
}

func validateMimeType(mimeType string) error {
	if len(mimeType) == 0 {
		return domainerr.ValidationErr("mime type", "can not be empty")
	}
	return nil
}

func validateSize(size uint64) error {
	if size == 0 {
		return domainerr.ValidationErr("size", "must be greater than 0")
	}
	return nil
}

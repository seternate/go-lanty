package asset

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	domainErrors "github.com/seternate/go-lanty/pkg/domain/error"
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
		return nil, domainErrors.InternalErr("failed to generate asset ID").WithCause(err)
	}

	return RehydrateAsset(id, assetURL, size, checksum, algorithm, mimeType)
}

func RehydrateAsset(id uuid.UUID, assetURL string, size uint64, checksum string, algorithm string, mimeType string) (*Asset, error) {
	validationErrors := domainErrors.ValidationErrs().WithMessage("failed to validate asset id=%s", id)

	u, err := url.Parse(assetURL)
	err = validationErrors.Wrap(domainErrors.ValidationErr("asset URL", "failed to parse URL").WithGot(assetURL).WithCause(err))
	if err != nil {
		return nil, fmt.Errorf("failed to parse for asset id=%s: %w", id, err)
	}

	parsedAlgo, err := ParseChecksumAlgorithm(algorithm)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, fmt.Errorf("failed to parse for asset id=%s: %w", id, err)
	}

	err = validationErrors.Wrap(validateChecksum(checksum))
	if err != nil {
		return nil, fmt.Errorf("failed to validate for asset id=%s: %w", id, err)
	}

	err = validationErrors.Wrap(validateMimeType(mimeType))
	if err != nil {
		return nil, fmt.Errorf("failed to validate for asset id=%s: %w", id, err)
	}

	err = validationErrors.Wrap(validateSize(size))
	if err != nil {
		return nil, fmt.Errorf("failed to validate for asset id=%s: %w", id, err)
	}

	if len(validationErrors.Errors) > 0 {
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

func (asset *Asset) CompareChecksum(checksum string) error {
	if checksum != asset.Checksum {
		return domainErrors.ValidationErr("checksum", "mismatch: asset id=%s", asset.ID).WithExpected(asset.Checksum).WithGot(checksum)
	}
	return nil
}

func validateChecksum(checksum string) error {
	if len(checksum) == 0 {
		return domainErrors.ValidationErr("checksum", "can not be empty")
	}
	return nil
}

func validateMimeType(mimeType string) error {
	if len(mimeType) == 0 {
		return domainErrors.ValidationErr("mime type", "can not be empty")
	}
	return nil
}

func validateSize(size uint64) error {
	if size == 0 {
		return domainErrors.ValidationErr("size", "must be greater than 0")
	}
	return nil
}

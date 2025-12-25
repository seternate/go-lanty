package asset

import (
	"fmt"
	"io"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/domain/asset"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

var _ CommandService = (*commandServiceImpl)(nil)

type commandServiceImpl struct {
	repository   asset.AssetRepository
	storage      StorageAdapter
	calculator   HashCalculator
	mimeDetector MimeTypeDetector
}

func NewCommandService(repository asset.AssetRepository, storage StorageAdapter, calculator HashCalculator, mimeDetector MimeTypeDetector) CommandService {
	return &commandServiceImpl{
		repository:   repository,
		storage:      storage,
		calculator:   calculator,
		mimeDetector: mimeDetector,
	}
}

func (c *commandServiceImpl) StoreNewAsset(cmd StoreNewAssetCommand) (*asset.Asset, error) {
	parsedAlgo, err := asset.ParseChecksumAlgorithm(cmd.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("failed to parse for url=%s: %w", cmd.URL, err)
	}

	hasher, err := c.calculator.GetHasher(parsedAlgo)
	if err != nil {
		return nil, fmt.Errorf("failed to get hasher for url=%s: %w", cmd.URL, err)
	}

	assetURL, err := url.Parse(cmd.URL)
	if err != nil {
		return nil, domainerr.ValidationErr("asset URL", "failed to parse").WithGot(cmd.URL).WithCause(err)
	}

	var dataReader io.Reader = cmd.Data
	mimeType := cmd.MimeType

	if mimeType == "" {
		detectedMimeType, preservedReader, err := c.mimeDetector.DetectAndPreserveReader(cmd.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to detect mime type for url=%s: %w", cmd.URL, err)
		}
		mimeType = detectedMimeType
		dataReader = preservedReader
	}

	bytesWritten, err := c.storage.Save(*assetURL, io.TeeReader(dataReader, hasher))
	if err != nil {
		return nil, fmt.Errorf("failed to save asset: %w", err)
	}
	defer func() {
		if err != nil {
			internalErr := c.storage.Delete(*assetURL)
			if internalErr != nil {
				log.Error().Err(internalErr).Msgf("failed to clean up asset binary (%s) after an error while storing new asset", cmd.URL)
			}
		}
	}()

	registeredAsset, err := asset.NewAsset(
		assetURL.String(),
		uint64(bytesWritten),
		cmd.Checksum,
		cmd.Algorithm,
		mimeType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new asset for url=%s: %w", cmd.URL, err)
	}

	calculatedChecksum, err := c.calculator.CalculateChecksum(hasher)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum for url=%s for asset id=%s: %w", cmd.URL, registeredAsset.ID.String(), err)
	}

	err = registeredAsset.CompareChecksum(calculatedChecksum)
	if err != nil {
		return nil, fmt.Errorf("failed to validate checksum: %w", err)
	}

	err = c.repository.CreateAsset(registeredAsset)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset in repository: %w", err)
	}

	return registeredAsset, nil
}

func (c *commandServiceImpl) DeleteAsset(id uuid.UUID) error {
	asset, err := c.repository.GetAsset(id)
	if err != nil {
		return fmt.Errorf("failed to get asset from repository: %w", err)
	}

	err = c.repository.DeleteAsset(id)
	if err != nil {
		return fmt.Errorf("failed to delete asset from repository: %w", err)
	}

	err = c.storage.Delete(asset.URL)
	if err != nil {
		return fmt.Errorf("failed to delete asset asset id=%s from storage: %w", id.String(), err)
	}

	return nil
}

package asset

import (
	"hash"
	"io"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/pkg/domain/asset"
)

type HashCalculator interface {
	asset.MD5Hasher
	asset.SHA256Hasher
	GetHasher(algorithm asset.ChecksumAlgorithm) (hash.Hash, error)
	CalculateChecksum(hasher hash.Hash) (string, error)
}

type MimeTypeDetector interface {
	DetectAndPreserveReader(reader io.Reader) (mimeType string, preservedReader io.Reader, err error)
}

type StorageAdapter interface {
	Fetch(url url.URL) (data io.ReadCloser, size uint64, err error)
	Save(url url.URL, data io.Reader) (written uint64, err error)
	Delete(url url.URL) error
}

type CommandService interface {
	StoreNewAsset(cmd StoreNewAssetCommand) (*asset.Asset, error)
	DeleteAsset(uuid.UUID) error
}

type StoreNewAssetCommand struct {
	URL       string
	Checksum  string
	Algorithm string
	MimeType  string
	Data      io.Reader
}

type QueryService interface {
	GetAsset(uuid.UUID) (*AssetView, error)
	GetAssetContent(uuid.UUID) (data io.ReadCloser, err error)
}

type AssetView struct {
	ID        uuid.UUID `db:"id"`
	URL       string    `db:"url"`
	Size      uint64    `db:"size"`
	Checksum  string    `db:"checksum"`
	Algorithm string    `db:"algorithm"`
	MimeType  string    `db:"mime_type"`
	CreatedAt time.Time `db:"created_at"`
}

type Service struct {
	Query   QueryService
	Command CommandService
}

func NewService(query QueryService, command CommandService) *Service {
	return &Service{
		Query:   query,
		Command: command,
	}
}

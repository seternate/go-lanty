package checksum

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"

	assetApp "github.com/seternate/go-lanty/internal/application/asset"
	"github.com/seternate/go-lanty/internal/domain/asset"
)

var _ assetApp.HashCalculator = (*calculator)(nil)

type calculator struct{}

func NewCalculator() assetApp.HashCalculator {
	return &calculator{}
}

func (c *calculator) MD5Hasher() hash.Hash {
	return md5.New()
}

func (c *calculator) SHA256Hasher() hash.Hash {
	return sha256.New()
}

func (c *calculator) GetHasher(algorithm asset.ChecksumAlgorithm) (hash.Hash, error) {
	switch algorithm {
	case asset.CHECKSUM_ALGORITHM_MD5:
		return c.MD5Hasher(), nil
	case asset.CHECKSUM_ALGORITHM_SHA256:
		return c.SHA256Hasher(), nil
	default:
		return nil, fmt.Errorf("unsupported checksum algorithm: %s", algorithm.String())
	}
}

func (c *calculator) CalculateChecksum(hasher hash.Hash) (string, error) {
	if hasher == nil {
		return "", fmt.Errorf("hasher is nil")
	}

	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum), nil
}

package asset

import (
	"hash"
	"strings"

	domainErrors "github.com/seternate/go-lanty/pkg/domain/error"
)

type ChecksumAlgorithm string

const (
	CHECKSUM_ALGORITHM_UNDEFINED ChecksumAlgorithm = ""
	CHECKSUM_ALGORITHM_MD5       ChecksumAlgorithm = "md5"
	CHECKSUM_ALGORITHM_SHA256    ChecksumAlgorithm = "sha-256"
)

type MD5Hasher interface {
	MD5Hasher() hash.Hash
}

type SHA256Hasher interface {
	SHA256Hasher() hash.Hash
}

var supportedChecksumAlgorithms = map[string]ChecksumAlgorithm{
	"md5":     CHECKSUM_ALGORITHM_MD5,
	"sha-256": CHECKSUM_ALGORITHM_SHA256,
}

func ParseChecksumAlgorithm(algorithm string) (ChecksumAlgorithm, error) {
	if checksumAlgo, ok := supportedChecksumAlgorithms[algorithm]; ok {
		return checksumAlgo, nil
	}

	available := make([]string, 0, len(supportedChecksumAlgorithms))
	for _, algo := range supportedChecksumAlgorithms {
		available = append(available, algo.String())
	}

	return CHECKSUM_ALGORITHM_UNDEFINED, domainErrors.ValidationErr("checksum algorithm", "undefined").WithExpected(strings.Join(available, ", ")).WithGot(algorithm)
}

func (algorithm ChecksumAlgorithm) String() string {
	return string(algorithm)
}

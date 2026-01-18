package content

import (
	"hash"
	"io"

	"github.com/seternate/go-lanty/internal/infrastructure/checksum"
	"github.com/seternate/go-lanty/internal/infrastructure/mimetype"
)

type HashCalculator interface {
	SHA256Hasher() hash.Hash
	CalculateChecksum(hasher hash.Hash) (string, error)
}

type MimeTypeDetector interface {
	DetectAndPreserveReader(reader io.Reader) (mimeType string, preservedReader io.Reader, err error)
}

func NewDefaultHashCalculator() HashCalculator {
	return checksum.NewCalculator()
}

func NewDefaultMimeTypeDetector() MimeTypeDetector {
	return mimetype.NewDetector()
}

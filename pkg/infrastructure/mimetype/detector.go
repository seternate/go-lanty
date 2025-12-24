package mimetype

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gabriel-vasile/mimetype"
	assetApp "github.com/seternate/go-lanty/pkg/application/asset"
)

var _ assetApp.MimeTypeDetector = (*detector)(nil)

type detector struct{}

func NewDetector() assetApp.MimeTypeDetector {
	return &detector{}
}

func (d *detector) DetectAndPreserveReader(reader io.Reader) (string, io.Reader, error) {
	// Read only the first 512 bytes for MIME type detection (sufficient for most file types)
	// This avoids loading large files entirely into memory
	const mimeDetectionBufferSize = 512
	mimeDetectionBuffer := make([]byte, mimeDetectionBufferSize)
	bytesRead, err := io.ReadFull(reader, mimeDetectionBuffer)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", nil, fmt.Errorf("failed to read first 512 bytes into buffer: %w", err)
	}
	// Trim buffer to actual bytes read
	mimeDetectionBuffer = mimeDetectionBuffer[:bytesRead]

	// Detect MIME type from the actual file content (first 512 bytes are sufficient)
	mtype, err := mimetype.DetectReader(bytes.NewReader(mimeDetectionBuffer))
	if err != nil {
		return "", nil, fmt.Errorf("error reading from sniffing bytes: %w", err)
	}

	// Combine the read bytes with the remaining stream for full file access
	// This allows the caller to read the complete file without loading it all into memory
	preservedReader := io.MultiReader(
		bytes.NewReader(mimeDetectionBuffer),
		reader,
	)

	return mtype.String(), preservedReader, nil
}

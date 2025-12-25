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
	const mimeDetectionBufferSize = 512
	mimeDetectionBuffer := make([]byte, mimeDetectionBufferSize)
	bytesRead, err := io.ReadFull(reader, mimeDetectionBuffer)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", nil, fmt.Errorf("failed to read first 512 bytes into buffer: %w", err)
	}
	mimeDetectionBuffer = mimeDetectionBuffer[:bytesRead]

	mtype, err := mimetype.DetectReader(bytes.NewReader(mimeDetectionBuffer))
	if err != nil {
		return "", nil, fmt.Errorf("error reading from sniffing bytes: %w", err)
	}

	preservedReader := io.MultiReader(
		bytes.NewReader(mimeDetectionBuffer),
		reader,
	)

	return mtype.String(), preservedReader, nil
}

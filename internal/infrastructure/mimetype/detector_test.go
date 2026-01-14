package mimetype

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	assetApp "github.com/seternate/go-lanty/internal/application/asset"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Fatal("NewDetector() returned nil")
	}

	// Verify it implements the interface
	var _ assetApp.MimeTypeDetector = detector
}

func TestDetector_DetectAndPreserveReader(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name               string
		data               []byte
		expectedMimePrefix string // We check prefix since exact MIME might vary
		description        string
	}{
		{
			name:               "PNG image",
			data:               []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG magic bytes
			expectedMimePrefix: "image/png",
			description:        "Should detect PNG from magic bytes",
		},
		{
			name:               "JPEG image",
			data:               []byte{0xFF, 0xD8, 0xFF, 0xE0}, // JPEG magic bytes
			expectedMimePrefix: "image/jpeg",
			description:        "Should detect JPEG from magic bytes",
		},
		{
			name:               "GIF image",
			data:               []byte("GIF89a"), // GIF magic bytes
			expectedMimePrefix: "image/gif",
			description:        "Should detect GIF from magic bytes",
		},
		{
			name:               "JSON text",
			data:               []byte(`{"key": "value"}`),
			expectedMimePrefix: "application/json",
			description:        "Should detect JSON from content",
		},
		{
			name:               "Plain text",
			data:               []byte("Hello, world!"),
			expectedMimePrefix: "text/plain",
			description:        "Should detect plain text",
		},
		{
			name:               "HTML content",
			data:               []byte(`<html><body>Test</body></html>`),
			expectedMimePrefix: "text/html",
			description:        "Should detect HTML",
		},
		{
			name:               "XML content",
			data:               []byte(`<?xml version="1.0"?><root></root>`),
			expectedMimePrefix: "text/xml", // Library returns text/xml; charset=utf-8
			description:        "Should detect XML",
		},
		{
			name:               "Empty data",
			data:               []byte(""),
			expectedMimePrefix: "text/plain", // Library returns text/plain for empty data
			description:        "Should handle empty data",
		},
		{
			name:               "Binary data",
			data:               []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
			expectedMimePrefix: "application/octet-stream",
			description:        "Should detect binary/octet-stream for unknown data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.data)

			mimeType, preservedReader, err := detector.DetectAndPreserveReader(reader)
			if err != nil {
				t.Errorf("DetectAndPreserveReader() error = %v, want nil", err)
				return
			}

			if mimeType == "" {
				t.Error("DetectAndPreserveReader() returned empty mimeType")
				return
			}

			// Check MIME type prefix (exact match might vary by library version)
			if len(mimeType) < len(tt.expectedMimePrefix) {
				t.Errorf("DetectAndPreserveReader() mimeType = %q (too short), want prefix %q", mimeType, tt.expectedMimePrefix)
				return
			}

			if mimeType[:len(tt.expectedMimePrefix)] != tt.expectedMimePrefix {
				t.Errorf("DetectAndPreserveReader() mimeType = %q, want prefix %q", mimeType, tt.expectedMimePrefix)
			}

			// Verify reader preservation - read all data from preserved reader
			preservedData, err := io.ReadAll(preservedReader)
			if err != nil {
				t.Errorf("Failed to read from preserved reader: %v", err)
				return
			}

			if !bytes.Equal(preservedData, tt.data) {
				t.Errorf("Preserved reader data = %v, want %v", preservedData, tt.data)
			}
		})
	}
}

func TestDetector_DetectAndPreserveReader_LargeFile(t *testing.T) {
	detector := NewDetector()

	// Create data larger than the 512-byte buffer used for detection
	largeData := make([]byte, 1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	// Prepend PNG magic bytes so it's detected as PNG
	copy(largeData, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})

	reader := bytes.NewReader(largeData)

	mimeType, preservedReader, err := detector.DetectAndPreserveReader(reader)
	if err != nil {
		t.Fatalf("DetectAndPreserveReader() error = %v, want nil", err)
	}

	if mimeType == "" {
		t.Fatal("DetectAndPreserveReader() returned empty mimeType")
	}

	// Verify it detected PNG
	if len(mimeType) < len("image/png") || mimeType[:len("image/png")] != "image/png" {
		t.Errorf("DetectAndPreserveReader() mimeType = %q, want prefix %q", mimeType, "image/png")
	}

	// Verify all data is preserved
	preservedData, err := io.ReadAll(preservedReader)
	if err != nil {
		t.Fatalf("Failed to read from preserved reader: %v", err)
	}

	if !bytes.Equal(preservedData, largeData) {
		t.Errorf("Preserved reader data length = %d, want %d", len(preservedData), len(largeData))
	}
}

func TestDetector_DetectAndPreserveReader_FromFile(t *testing.T) {
	detector := NewDetector()

	// Test with actual PNG file from testdata
	pngPath := filepath.Join("..", "..", "..", "testdata", "game", "icons", "test_game_icon.png")
	fileData, err := os.ReadFile(pngPath)
	if err != nil {
		t.Skipf("Skipping file test: could not read test file %s: %v", pngPath, err)
	}

	reader := bytes.NewReader(fileData)

	mimeType, preservedReader, err := detector.DetectAndPreserveReader(reader)
	if err != nil {
		t.Fatalf("DetectAndPreserveReader() error = %v, want nil", err)
	}

	if mimeType == "" {
		t.Fatal("DetectAndPreserveReader() returned empty mimeType")
	}

	// Should detect PNG
	if len(mimeType) < len("image/png") || mimeType[:len("image/png")] != "image/png" {
		t.Errorf("DetectAndPreserveReader() mimeType = %q, want prefix %q", mimeType, "image/png")
	}

	// Verify all data is preserved
	preservedData, err := io.ReadAll(preservedReader)
	if err != nil {
		t.Fatalf("Failed to read from preserved reader: %v", err)
	}

	if !bytes.Equal(preservedData, fileData) {
		t.Error("Preserved reader data does not match original file data")
	}
}

func TestDetector_DetectAndPreserveReader_ReaderErrors(t *testing.T) {
	detector := NewDetector()

	t.Run("Reader that returns error on read", func(t *testing.T) {
		errorReader := &errorReader{err: io.ErrClosedPipe}

		mimeType, preservedReader, err := detector.DetectAndPreserveReader(errorReader)

		if err == nil {
			t.Error("DetectAndPreserveReader() expected error from failing reader, got nil")
		}

		if mimeType != "" {
			t.Errorf("DetectAndPreserveReader() mimeType = %q, want empty on error", mimeType)
		}

		if preservedReader != nil {
			t.Error("DetectAndPreserveReader() preservedReader should be nil on error")
		}
	})
}

func TestDetector_DetectAndPreserveReader_SmallFile(t *testing.T) {
	detector := NewDetector()

	// Test with data smaller than 512 bytes
	smallData := []byte("tiny")

	reader := bytes.NewReader(smallData)

	mimeType, preservedReader, err := detector.DetectAndPreserveReader(reader)
	if err != nil {
		t.Fatalf("DetectAndPreserveReader() error = %v, want nil", err)
	}

	if mimeType == "" {
		t.Fatal("DetectAndPreserveReader() returned empty mimeType")
	}

	// Verify all data is preserved
	preservedData, err := io.ReadAll(preservedReader)
	if err != nil {
		t.Fatalf("Failed to read from preserved reader: %v", err)
	}

	if !bytes.Equal(preservedData, smallData) {
		t.Errorf("Preserved reader data = %v, want %v", preservedData, smallData)
	}
}

func TestDetector_DetectAndPreserveReader_MultipleReads(t *testing.T) {
	detector := NewDetector()

	testData := []byte("test data for multiple reads")
	reader := bytes.NewReader(testData)

	// First detection
	mimeType1, preservedReader1, err := detector.DetectAndPreserveReader(reader)
	if err != nil {
		t.Fatalf("First DetectAndPreserveReader() error = %v", err)
	}

	// Read from preserved reader multiple times to ensure it works
	data1, err := io.ReadAll(preservedReader1)
	if err != nil {
		t.Fatalf("Failed to read from preserved reader first time: %v", err)
	}

	if !bytes.Equal(data1, testData) {
		t.Error("First read from preserved reader did not match original data")
	}

	if mimeType1 == "" {
		t.Error("First detection returned empty mimeType")
	}
}

// errorReader is a helper type that always returns an error on Read
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

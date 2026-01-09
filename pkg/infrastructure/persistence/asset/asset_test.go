package asset

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/pkg/domain/asset"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

func TestAssetRow_Assemble(t *testing.T) {
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("Valid asset row with MD5", func(t *testing.T) {
		row := AssetRow{
			ID:        testID,
			URL:       "https://example.com/file.bin",
			Size:      1024,
			Checksum:  "abc123def456",
			Algorithm: "md5",
			MimeType:  "application/octet-stream",
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled == nil {
			t.Fatal("expected asset, got nil")
		}
		if assembled.ID != testID {
			t.Errorf("expected ID %v, got %v", testID, assembled.ID)
		}
		if assembled.URL.String() != "https://example.com/file.bin" {
			t.Errorf("expected URL %q, got %q", "https://example.com/file.bin", assembled.URL.String())
		}
		if assembled.Size != 1024 {
			t.Errorf("expected Size 1024, got %d", assembled.Size)
		}
		if assembled.Checksum != "abc123def456" {
			t.Errorf("expected Checksum %q, got %q", "abc123def456", assembled.Checksum)
		}
		if assembled.Algorithm.String() != "md5" {
			t.Errorf("expected Algorithm %q, got %q", "md5", assembled.Algorithm.String())
		}
		if assembled.MimeType != "application/octet-stream" {
			t.Errorf("expected MimeType %q, got %q", "application/octet-stream", assembled.MimeType)
		}
	})

	t.Run("Valid asset row with SHA-256", func(t *testing.T) {
		row := AssetRow{
			ID:        testID,
			URL:       "https://example.com/file.bin",
			Size:      2048,
			Checksum:  "sha256hash",
			Algorithm: "sha-256",
			MimeType:  "application/pdf",
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled == nil {
			t.Fatal("expected asset, got nil")
		}
		if assembled.Algorithm.String() != "sha-256" {
			t.Errorf("expected Algorithm %q, got %q", "sha-256", assembled.Algorithm.String())
		}
		if assembled.MimeType != "application/pdf" {
			t.Errorf("expected MimeType %q, got %q", "application/pdf", assembled.MimeType)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		row := AssetRow{
			ID:        testID,
			URL:       "http://[invalid",
			Size:      1024,
			Checksum:  "abc123",
			Algorithm: "md5",
			MimeType:  "application/octet-stream",
		}

		assembled, err := row.Assemble()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})

	t.Run("Invalid checksum algorithm", func(t *testing.T) {
		row := AssetRow{
			ID:        testID,
			URL:       "https://example.com/file.bin",
			Size:      1024,
			Checksum:  "abc123",
			Algorithm: "invalid-algorithm",
			MimeType:  "application/octet-stream",
		}

		assembled, err := row.Assemble()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})

	t.Run("Empty checksum", func(t *testing.T) {
		row := AssetRow{
			ID:        testID,
			URL:       "https://example.com/file.bin",
			Size:      1024,
			Checksum:  "",
			Algorithm: "md5",
			MimeType:  "application/octet-stream",
		}

		assembled, err := row.Assemble()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})

	t.Run("Empty mime type", func(t *testing.T) {
		row := AssetRow{
			ID:        testID,
			URL:       "https://example.com/file.bin",
			Size:      1024,
			Checksum:  "abc123",
			Algorithm: "md5",
			MimeType:  "",
		}

		assembled, err := row.Assemble()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})

	t.Run("Zero size", func(t *testing.T) {
		row := AssetRow{
			ID:        testID,
			URL:       "https://example.com/file.bin",
			Size:      0,
			Checksum:  "abc123",
			Algorithm: "md5",
			MimeType:  "application/octet-stream",
		}

		assembled, err := row.Assemble()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})
}

func TestDisassemble(t *testing.T) {
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("Valid asset with MD5", func(t *testing.T) {
		asset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		row := Disassemble(asset)
		if row == nil {
			t.Fatal("expected row, got nil")
		}
		if row.ID != testID {
			t.Errorf("expected ID %v, got %v", testID, row.ID)
		}
		if row.URL != "https://example.com/file.bin" {
			t.Errorf("expected URL %q, got %q", "https://example.com/file.bin", row.URL)
		}
		if row.Size != 1024 {
			t.Errorf("expected Size 1024, got %d", row.Size)
		}
		if row.Checksum != "abc123" {
			t.Errorf("expected Checksum %q, got %q", "abc123", row.Checksum)
		}
		if row.Algorithm != "md5" {
			t.Errorf("expected Algorithm %q, got %q", "md5", row.Algorithm)
		}
		if row.MimeType != "application/octet-stream" {
			t.Errorf("expected MimeType %q, got %q", "application/octet-stream", row.MimeType)
		}
	})

	t.Run("Valid asset with SHA-256", func(t *testing.T) {
		asset, err := asset.RehydrateAsset(testID, "https://example.com/file.pdf", 2048, "sha256hash", "sha-256", "application/pdf")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		row := Disassemble(asset)
		if row == nil {
			t.Fatal("expected row, got nil")
		}
		if row.Algorithm != "sha-256" {
			t.Errorf("expected Algorithm %q, got %q", "sha-256", row.Algorithm)
		}
		if row.MimeType != "application/pdf" {
			t.Errorf("expected MimeType %q, got %q", "application/pdf", row.MimeType)
		}
		if row.Size != 2048 {
			t.Errorf("expected Size 2048, got %d", row.Size)
		}
	})

	t.Run("Round trip: Disassemble then Assemble", func(t *testing.T) {
		originalAsset, err := asset.RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		row := Disassemble(originalAsset)
		reassembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("failed to reassemble: %v", err)
		}

		if reassembled.ID != originalAsset.ID {
			t.Errorf("expected ID %v, got %v", originalAsset.ID, reassembled.ID)
		}
		if reassembled.URL.String() != originalAsset.URL.String() {
			t.Errorf("expected URL %q, got %q", originalAsset.URL.String(), reassembled.URL.String())
		}
		if reassembled.Size != originalAsset.Size {
			t.Errorf("expected Size %d, got %d", originalAsset.Size, reassembled.Size)
		}
		if reassembled.Checksum != originalAsset.Checksum {
			t.Errorf("expected Checksum %q, got %q", originalAsset.Checksum, reassembled.Checksum)
		}
		if reassembled.Algorithm != originalAsset.Algorithm {
			t.Errorf("expected Algorithm %v, got %v", originalAsset.Algorithm, reassembled.Algorithm)
		}
		if reassembled.MimeType != originalAsset.MimeType {
			t.Errorf("expected MimeType %q, got %q", originalAsset.MimeType, reassembled.MimeType)
		}
	})

	t.Run("Complex URL with query parameters", func(t *testing.T) {
		url := "https://example.com/file.bin?version=1&token=abc123"
		asset, err := asset.RehydrateAsset(testID, url, 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("failed to create test asset: %v", err)
		}

		row := Disassemble(asset)
		if row.URL != url {
			t.Errorf("expected URL %q, got %q", url, row.URL)
		}
	})
}

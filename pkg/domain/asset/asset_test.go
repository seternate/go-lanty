package asset

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

func TestNewAsset(t *testing.T) {
	t.Run("Valid asset", func(t *testing.T) {
		asset, err := NewAsset("https://example.com/file.bin", 1024, "abc123", "md5", "application/octet-stream")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if asset == nil {
			t.Fatal("expected asset, got nil")
		}
		if asset.ID == uuid.Nil {
			t.Error("expected non-nil UUID")
		}
		if asset.URL.String() != "https://example.com/file.bin" {
			t.Errorf("expected URL %q, got %q", "https://example.com/file.bin", asset.URL.String())
		}
		if asset.Size != 1024 {
			t.Errorf("expected Size 1024, got %d", asset.Size)
		}
		if asset.Checksum != "abc123" {
			t.Errorf("expected Checksum %q, got %q", "abc123", asset.Checksum)
		}
		if asset.Algorithm != CHECKSUM_ALGORITHM_MD5 {
			t.Errorf("expected Algorithm %v, got %v", CHECKSUM_ALGORITHM_MD5, asset.Algorithm)
		}
		if asset.MimeType != "application/octet-stream" {
			t.Errorf("expected MimeType %q, got %q", "application/octet-stream", asset.MimeType)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		asset, err := NewAsset("http://[invalid", 1024, "abc123", "md5", "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid checksum algorithm", func(t *testing.T) {
		asset, err := NewAsset("https://example.com/file.bin", 1024, "abc123", "invalid", "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var validationErrs *domainerr.ValidationErrors
		if !errors.As(err, &validationErrs) {
			t.Errorf("expected ValidationErrors, got %T", err)
		}
	})

	t.Run("Empty checksum", func(t *testing.T) {
		asset, err := NewAsset("https://example.com/file.bin", 1024, "", "md5", "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var validationErrs *domainerr.ValidationErrors
		if !errors.As(err, &validationErrs) {
			t.Errorf("expected ValidationErrors, got %T", err)
		}
	})

	t.Run("Empty mime type", func(t *testing.T) {
		asset, err := NewAsset("https://example.com/file.bin", 1024, "abc123", "md5", "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var validationErrs *domainerr.ValidationErrors
		if !errors.As(err, &validationErrs) {
			t.Errorf("expected ValidationErrors, got %T", err)
		}
	})

	t.Run("Zero size", func(t *testing.T) {
		asset, err := NewAsset("https://example.com/file.bin", 0, "abc123", "md5", "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var validationErrs *domainerr.ValidationErrors
		if !errors.As(err, &validationErrs) {
			t.Errorf("expected ValidationErrors, got %T", err)
		}
	})

	t.Run("Multiple validation errors", func(t *testing.T) {
		asset, err := NewAsset("http://[invalid", 0, "", "invalid", "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})
}

func TestRehydrateAsset(t *testing.T) {
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("Valid asset", func(t *testing.T) {
		asset, err := RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "sha-256", "application/octet-stream")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if asset == nil {
			t.Fatal("expected asset, got nil")
		}
		if asset.ID != testID {
			t.Errorf("expected ID %v, got %v", testID, asset.ID)
		}
		if asset.URL.String() != "https://example.com/file.bin" {
			t.Errorf("expected URL %q, got %q", "https://example.com/file.bin", asset.URL.String())
		}
		if asset.Size != 1024 {
			t.Errorf("expected Size 1024, got %d", asset.Size)
		}
		if asset.Checksum != "abc123" {
			t.Errorf("expected Checksum %q, got %q", "abc123", asset.Checksum)
		}
		if asset.Algorithm != CHECKSUM_ALGORITHM_SHA256 {
			t.Errorf("expected Algorithm %v, got %v", CHECKSUM_ALGORITHM_SHA256, asset.Algorithm)
		}
		if asset.MimeType != "application/octet-stream" {
			t.Errorf("expected MimeType %q, got %q", "application/octet-stream", asset.MimeType)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		asset, err := RehydrateAsset(testID, "http://[invalid", 1024, "abc123", "md5", "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid checksum algorithm", func(t *testing.T) {
		asset, err := RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "invalid", "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})

	t.Run("Empty checksum", func(t *testing.T) {
		asset, err := RehydrateAsset(testID, "https://example.com/file.bin", 1024, "", "md5", "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})

	t.Run("Empty mime type", func(t *testing.T) {
		asset, err := RehydrateAsset(testID, "https://example.com/file.bin", 1024, "abc123", "md5", "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})

	t.Run("Zero size", func(t *testing.T) {
		asset, err := RehydrateAsset(testID, "https://example.com/file.bin", 0, "abc123", "md5", "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if asset != nil {
			t.Error("expected nil asset")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})
}

func TestAsset_CompareChecksum(t *testing.T) {
	asset, err := NewAsset("https://example.com/file.bin", 1024, "expected-checksum", "md5", "application/octet-stream")
	if err != nil {
		t.Fatalf("failed to create test asset: %v", err)
	}

	t.Run("Matching checksum", func(t *testing.T) {
		err := asset.CompareChecksum("expected-checksum")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Non-matching checksum", func(t *testing.T) {
		err := asset.CompareChecksum("different-checksum")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil {
			if validationErr.Expected != "expected-checksum" {
				t.Errorf("expected Expected %q, got %q", "expected-checksum", validationErr.Expected)
			}
			if validationErr.Got != "different-checksum" {
				t.Errorf("expected Got %q, got %q", "different-checksum", validationErr.Got)
			}
		}
	})

	t.Run("Empty checksum", func(t *testing.T) {
		err := asset.CompareChecksum("")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})
}

func TestValidateChecksum(t *testing.T) {
	t.Run("Valid checksum", func(t *testing.T) {
		err := validateChecksum("abc123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Empty checksum", func(t *testing.T) {
		err := validateChecksum("")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil && validationErr.Field != "checksum" {
			t.Errorf("expected Field %q, got %q", "checksum", validationErr.Field)
		}
	})
}

func TestValidateMimeType(t *testing.T) {
	t.Run("Valid mime type", func(t *testing.T) {
		err := validateMimeType("application/octet-stream")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Empty mime type", func(t *testing.T) {
		err := validateMimeType("")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil && validationErr.Field != "mime type" {
			t.Errorf("expected Field %q, got %q", "mime type", validationErr.Field)
		}
	})
}

func TestValidateSize(t *testing.T) {
	t.Run("Valid size", func(t *testing.T) {
		err := validateSize(1024)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Zero size", func(t *testing.T) {
		err := validateSize(0)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil && validationErr.Field != "size" {
			t.Errorf("expected Field %q, got %q", "size", validationErr.Field)
		}
	})

	t.Run("Large size", func(t *testing.T) {
		err := validateSize(18446744073709551615) // max uint64
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

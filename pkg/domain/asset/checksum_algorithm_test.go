package asset

import (
	"errors"
	"testing"

	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

func TestChecksumAlgorithm_String(t *testing.T) {
	tests := []struct {
		name     string
		algorithm ChecksumAlgorithm
		expected string
	}{
		{"MD5", CHECKSUM_ALGORITHM_MD5, "md5"},
		{"SHA256", CHECKSUM_ALGORITHM_SHA256, "sha-256"},
		{"Undefined", CHECKSUM_ALGORITHM_UNDEFINED, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.algorithm.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseChecksumAlgorithm(t *testing.T) {
	tests := []struct {
		name        string
		algorithm   string
		expected    ChecksumAlgorithm
		expectError bool
	}{
		{"Valid MD5", "md5", CHECKSUM_ALGORITHM_MD5, false},
		{"Valid SHA256", "sha-256", CHECKSUM_ALGORITHM_SHA256, false},
		{"Invalid algorithm", "sha1", CHECKSUM_ALGORITHM_UNDEFINED, true},
		{"Empty string", "", CHECKSUM_ALGORITHM_UNDEFINED, true},
		{"Case sensitive - uppercase", "MD5", CHECKSUM_ALGORITHM_UNDEFINED, true},
		{"Case sensitive - mixed", "Sha-256", CHECKSUM_ALGORITHM_UNDEFINED, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseChecksumAlgorithm(tt.algorithm)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				var validationErr *domainerr.ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("expected ValidationError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if got != tt.expected {
					t.Errorf("ParseChecksumAlgorithm() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

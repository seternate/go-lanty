package checksum

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"testing"

	assetApp "github.com/seternate/go-lanty/internal/application/asset"
	"github.com/seternate/go-lanty/internal/domain/asset"
)

func TestNewCalculator(t *testing.T) {
	calc := NewCalculator()
	if calc == nil {
		t.Fatal("NewCalculator() returned nil")
	}

	// Verify it implements the interface
	var _ assetApp.HashCalculator = calc
}

func TestCalculator_MD5Hasher(t *testing.T) {
	calc := NewCalculator()

	t.Run("Returns valid MD5 hasher", func(t *testing.T) {
		hasher := calc.MD5Hasher()
		if hasher == nil {
			t.Fatal("MD5Hasher() returned nil")
		}

		// Test that it's actually an MD5 hasher by hashing known data
		testData := []byte("hello world")
		hasher.Write(testData)
		sum := hasher.Sum(nil)

		// Expected MD5 hash of "hello world"
		expectedHash, _ := hex.DecodeString("5eb63bbbe01eeed093cb22bb8f5acdc3")
		if hex.EncodeToString(sum) != hex.EncodeToString(expectedHash) {
			t.Errorf("MD5Hasher() produced incorrect hash. Got %x, expected %x", sum, expectedHash)
		}
	})

	t.Run("Returns new instance each time", func(t *testing.T) {
		hasher1 := calc.MD5Hasher()
		hasher2 := calc.MD5Hasher()

		// Write different data to each
		hasher1.Write([]byte("test1"))
		hasher2.Write([]byte("test2"))

		sum1 := hasher1.Sum(nil)
		sum2 := hasher2.Sum(nil)

		if hex.EncodeToString(sum1) == hex.EncodeToString(sum2) {
			t.Error("MD5Hasher() returned the same instance, expected new instances")
		}
	})
}

func TestCalculator_SHA256Hasher(t *testing.T) {
	calc := NewCalculator()

	t.Run("Returns valid SHA256 hasher", func(t *testing.T) {
		hasher := calc.SHA256Hasher()
		if hasher == nil {
			t.Fatal("SHA256Hasher() returned nil")
		}

		// Test that it's actually a SHA256 hasher by hashing known data
		testData := []byte("hello world")
		hasher.Write(testData)
		sum := hasher.Sum(nil)

		// Expected SHA256 hash of "hello world"
		expectedHash, _ := hex.DecodeString("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9")
		if hex.EncodeToString(sum) != hex.EncodeToString(expectedHash) {
			t.Errorf("SHA256Hasher() produced incorrect hash. Got %x, expected %x", sum, expectedHash)
		}
	})

	t.Run("Returns new instance each time", func(t *testing.T) {
		hasher1 := calc.SHA256Hasher()
		hasher2 := calc.SHA256Hasher()

		// Write different data to each
		hasher1.Write([]byte("test1"))
		hasher2.Write([]byte("test2"))

		sum1 := hasher1.Sum(nil)
		sum2 := hasher2.Sum(nil)

		if hex.EncodeToString(sum1) == hex.EncodeToString(sum2) {
			t.Error("SHA256Hasher() returned the same instance, expected new instances")
		}
	})
}

func TestCalculator_GetHasher(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name        string
		algorithm   asset.ChecksumAlgorithm
		expectError bool
		validate    func(t *testing.T, h interface{})
	}{
		{
			name:        "MD5 algorithm",
			algorithm:   asset.CHECKSUM_ALGORITHM_MD5,
			expectError: false,
			validate: func(t *testing.T, h interface{}) {
				hasher, ok := h.(hash.Hash)
				if !ok {
					t.Errorf("expected hash.Hash, got %T", h)
					return
				}
				// Verify it works by hashing data
				hasher.Write([]byte("test"))
				sum := hasher.Sum(nil)
				if len(sum) != md5.Size {
					t.Errorf("expected hash size %d, got %d", md5.Size, len(sum))
				}
			},
		},
		{
			name:        "SHA256 algorithm",
			algorithm:   asset.CHECKSUM_ALGORITHM_SHA256,
			expectError: false,
			validate: func(t *testing.T, h interface{}) {
				hasher, ok := h.(hash.Hash)
				if !ok {
					t.Errorf("expected hash.Hash, got %T", h)
					return
				}
				// Verify it works by hashing data
				hasher.Write([]byte("test"))
				sum := hasher.Sum(nil)
				if len(sum) != sha256.Size {
					t.Errorf("expected hash size %d, got %d", sha256.Size, len(sum))
				}
			},
		},
		{
			name:        "Undefined algorithm",
			algorithm:   asset.CHECKSUM_ALGORITHM_UNDEFINED,
			expectError: true,
			validate:    nil,
		},
		{
			name:        "Invalid algorithm",
			algorithm:   asset.ChecksumAlgorithm("invalid"),
			expectError: true,
			validate:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher, err := calc.GetHasher(tt.algorithm)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if hasher != nil {
					t.Error("expected nil hasher on error")
				}
				// Verify error message mentions the algorithm
				if err.Error() == "" {
					t.Error("error message is empty")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if hasher == nil {
					t.Fatal("expected hasher, got nil")
				}
				if tt.validate != nil {
					tt.validate(t, hasher)
				}
			}
		})
	}
}

func TestCalculator_CalculateChecksum(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name        string
		setupHasher func() hash.Hash
		data        []byte
		expectedHex string
		expectError bool
	}{
		{
			name: "MD5 checksum calculation",
			setupHasher: func() hash.Hash {
				return calc.MD5Hasher()
			},
			data:        []byte("hello world"),
			expectedHex: "5eb63bbbe01eeed093cb22bb8f5acdc3",
			expectError: false,
		},
		{
			name: "SHA256 checksum calculation",
			setupHasher: func() hash.Hash {
				return calc.SHA256Hasher()
			},
			data:        []byte("hello world"),
			expectedHex: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
			expectError: false,
		},
		{
			name: "Empty data MD5",
			setupHasher: func() hash.Hash {
				return calc.MD5Hasher()
			},
			data:        []byte(""),
			expectedHex: "d41d8cd98f00b204e9800998ecf8427e",
			expectError: false,
		},
		{
			name: "Empty data SHA256",
			setupHasher: func() hash.Hash {
				return calc.SHA256Hasher()
			},
			data:        []byte(""),
			expectedHex: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			expectError: false,
		},
		{
			name: "Nil hasher",
			setupHasher: func() hash.Hash {
				return nil
			},
			data:        []byte("test"),
			expectedHex: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := tt.setupHasher()

			if hasher != nil {
				// Write data to hasher
				_, err := hasher.Write(tt.data)
				if err != nil {
					t.Fatalf("failed to write to hasher: %v", err)
				}
			}

			result, err := calc.CalculateChecksum(hasher)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if result != "" {
					t.Errorf("expected empty result on error, got %q", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if result != tt.expectedHex {
					t.Errorf("CalculateChecksum() = %q, want %q", result, tt.expectedHex)
				}
				// Verify it's valid hex
				if _, err := hex.DecodeString(result); err != nil {
					t.Errorf("result is not valid hex: %v", err)
				}
			}
		})
	}
}

func TestCalculator_Integration(t *testing.T) {
	calc := NewCalculator()

	t.Run("Full workflow with MD5", func(t *testing.T) {
		hasher, err := calc.GetHasher(asset.CHECKSUM_ALGORITHM_MD5)
		if err != nil {
			t.Fatalf("unexpected error getting hasher: %v", err)
		}

		testData := []byte("integration test data")
		hasher.Write(testData)

		checksum, err := calc.CalculateChecksum(hasher)
		if err != nil {
			t.Fatalf("unexpected error calculating checksum: %v", err)
		}

		// Verify checksum format (should be hex string)
		if len(checksum) != 32 { // MD5 produces 32 hex characters
			t.Errorf("expected checksum length 32, got %d", len(checksum))
		}

		// Verify it's valid hex
		if _, err := hex.DecodeString(checksum); err != nil {
			t.Errorf("checksum is not valid hex: %v", err)
		}
	})

	t.Run("Full workflow with SHA256", func(t *testing.T) {
		hasher, err := calc.GetHasher(asset.CHECKSUM_ALGORITHM_SHA256)
		if err != nil {
			t.Fatalf("unexpected error getting hasher: %v", err)
		}

		testData := []byte("integration test data")
		hasher.Write(testData)

		checksum, err := calc.CalculateChecksum(hasher)
		if err != nil {
			t.Fatalf("unexpected error calculating checksum: %v", err)
		}

		// Verify checksum format (should be hex string)
		if len(checksum) != 64 { // SHA256 produces 64 hex characters
			t.Errorf("expected checksum length 64, got %d", len(checksum))
		}

		// Verify it's valid hex
		if _, err := hex.DecodeString(checksum); err != nil {
			t.Errorf("checksum is not valid hex: %v", err)
		}
	})
}

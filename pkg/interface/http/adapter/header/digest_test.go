package header

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestDecodeContentDigestHeader(t *testing.T) {
	t.Run("Valid MD5 digest", func(t *testing.T) {
		// MD5 hash of "test" = 098f6bcd4621d373cade4e832627b4f6
		checksumHex := "098f6bcd4621d373cade4e832627b4f6"
		checksumBytes, _ := hex.DecodeString(checksumHex)
		checksumBase64 := base64.StdEncoding.EncodeToString(checksumBytes)
		digestValue := "md5=" + checksumBase64

		algorithm, checksum, err := DecodeContentDigestHeader(digestValue)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if algorithm != "md5" {
			t.Errorf("expected algorithm %q, got %q", "md5", algorithm)
		}
		if checksum != checksumHex {
			t.Errorf("expected checksum %q, got %q", checksumHex, checksum)
		}
	})

	t.Run("Valid SHA256 digest", func(t *testing.T) {
		// SHA256 hash of "test" = 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08
		checksumHex := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
		checksumBytes, _ := hex.DecodeString(checksumHex)
		checksumBase64 := base64.StdEncoding.EncodeToString(checksumBytes)
		digestValue := "sha256=" + checksumBase64

		algorithm, checksum, err := DecodeContentDigestHeader(digestValue)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if algorithm != "sha256" {
			t.Errorf("expected algorithm %q, got %q", "sha256", algorithm)
		}
		if checksum != checksumHex {
			t.Errorf("expected checksum %q, got %q", checksumHex, checksum)
		}
	})

	t.Run("Algorithm case normalization", func(t *testing.T) {
		checksumHex := "098f6bcd4621d373cade4e832627b4f6"
		checksumBytes, _ := hex.DecodeString(checksumHex)
		checksumBase64 := base64.StdEncoding.EncodeToString(checksumBytes)

		testCases := []struct {
			name      string
			algorithm string
			expected  string
		}{
			{"Uppercase", "MD5", "md5"},
			{"Mixed case", "Sha256", "sha256"},
			{"Lowercase", "md5", "md5"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				digestValue := tc.algorithm + "=" + checksumBase64
				algorithm, _, err := DecodeContentDigestHeader(digestValue)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if algorithm != tc.expected {
					t.Errorf("expected algorithm %q, got %q", tc.expected, algorithm)
				}
			})
		}
	})

	t.Run("Whitespace handling", func(t *testing.T) {
		checksumHex := "098f6bcd4621d373cade4e832627b4f6"
		checksumBytes, _ := hex.DecodeString(checksumHex)
		checksumBase64 := base64.StdEncoding.EncodeToString(checksumBytes)

		testCases := []struct {
			name        string
			digestValue string
		}{
			{"Spaces around equals", " md5 = " + checksumBase64},
			{"Tabs around equals", "\tmd5\t=\t" + checksumBase64},
			{"Mixed whitespace", "  md5  =  " + checksumBase64 + "  "},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				algorithm, checksum, err := DecodeContentDigestHeader(tc.digestValue)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if algorithm != "md5" {
					t.Errorf("expected algorithm %q, got %q", "md5", algorithm)
				}
				if checksum != checksumHex {
					t.Errorf("expected checksum %q, got %q", checksumHex, checksum)
				}
			})
		}
	})

	t.Run("Empty digest value", func(t *testing.T) {
		algorithm, checksum, err := DecodeContentDigestHeader("")
		if err == nil {
			t.Fatal("expected error for empty digest value")
		}
		if algorithm != "" {
			t.Errorf("expected empty algorithm, got %q", algorithm)
		}
		if checksum != "" {
			t.Errorf("expected empty checksum, got %q", checksum)
		}
		if !strings.Contains(err.Error(), "empty") {
			t.Errorf("expected error message to contain 'empty', got %q", err.Error())
		}
	})

	t.Run("Missing equals sign", func(t *testing.T) {
		testCases := []struct {
			name        string
			digestValue string
		}{
			{"No equals", "md5abc123"},
			{"Only algorithm", "md5"},
			{"Only checksum", "abc123"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				algorithm, checksum, err := DecodeContentDigestHeader(tc.digestValue)
				if err == nil {
					t.Fatal("expected error for invalid format")
				}
				if algorithm != "" {
					t.Errorf("expected empty algorithm, got %q", algorithm)
				}
				if checksum != "" {
					t.Errorf("expected empty checksum, got %q", checksum)
				}
				if !strings.Contains(err.Error(), "invalid digest value format") {
					t.Errorf("expected error message to contain 'invalid digest value format', got %q", err.Error())
				}
			})
		}
	})

	t.Run("Invalid base64 checksum", func(t *testing.T) {
		testCases := []struct {
			name        string
			digestValue string
		}{
			{"Invalid base64 characters", "md5=abc!@#"},
			{"Incomplete base64", "md5=abc"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				algorithm, checksum, err := DecodeContentDigestHeader(tc.digestValue)
				if err == nil {
					t.Fatal("expected error for invalid base64")
				}
				if algorithm != "" {
					t.Errorf("expected empty algorithm, got %q", algorithm)
				}
				if checksum != "" {
					t.Errorf("expected empty checksum, got %q", checksum)
				}
				if !strings.Contains(err.Error(), "failed to decode base64") {
					t.Errorf("expected error message to contain 'failed to decode base64', got %q", err.Error())
				}
			})
		}
	})
}

func TestEncodeContentDigestHeader(t *testing.T) {
	t.Run("Valid MD5 digest", func(t *testing.T) {
		algorithm := "md5"
		checksum := "098f6bcd4621d373cade4e832627b4f6"

		result, err := EncodeContentDigestHeader(algorithm, checksum)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify format
		if !strings.Contains(result, "=") {
			t.Errorf("expected result to contain '=', got %q", result)
		}

		parts := strings.SplitN(result, "=", 2)
		if len(parts) != 2 {
			t.Fatalf("expected result to have format algorithm=checksum, got %q", result)
		}

		if parts[0] != "md5" {
			t.Errorf("expected algorithm %q, got %q", "md5", parts[0])
		}

		// Decode base64 and verify it matches original hex
		decodedBytes, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			t.Fatalf("failed to decode base64: %v", err)
		}
		decodedHex := hex.EncodeToString(decodedBytes)
		if decodedHex != checksum {
			t.Errorf("expected checksum %q, got %q", checksum, decodedHex)
		}
	})

	t.Run("Valid SHA256 digest", func(t *testing.T) {
		algorithm := "sha256"
		checksum := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"

		result, err := EncodeContentDigestHeader(algorithm, checksum)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		parts := strings.SplitN(result, "=", 2)
		if len(parts) != 2 {
			t.Fatalf("expected result to have format algorithm=checksum, got %q", result)
		}

		if parts[0] != "sha256" {
			t.Errorf("expected algorithm %q, got %q", "sha256", parts[0])
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			t.Fatalf("failed to decode base64: %v", err)
		}
		decodedHex := hex.EncodeToString(decodedBytes)
		if decodedHex != checksum {
			t.Errorf("expected checksum %q, got %q", checksum, decodedHex)
		}
	})

	t.Run("Algorithm case normalization", func(t *testing.T) {
		checksum := "098f6bcd4621d373cade4e832627b4f6"

		testCases := []struct {
			name      string
			algorithm string
			expected  string
		}{
			{"Uppercase", "MD5", "md5"},
			{"Mixed case", "Sha256", "sha256"},
			{"Lowercase", "md5", "md5"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := EncodeContentDigestHeader(tc.algorithm, checksum)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				parts := strings.SplitN(result, "=", 2)
				if len(parts) != 2 {
					t.Fatalf("expected result to have format algorithm=checksum, got %q", result)
				}

				if parts[0] != tc.expected {
					t.Errorf("expected algorithm %q, got %q", tc.expected, parts[0])
				}
			})
		}
	})

	t.Run("Empty algorithm", func(t *testing.T) {
		checksum := "098f6bcd4621d373cade4e832627b4f6"

		result, err := EncodeContentDigestHeader("", checksum)
		if err == nil {
			t.Fatal("expected error for empty algorithm")
		}
		if result != "" {
			t.Errorf("expected empty result, got %q", result)
		}
		if !strings.Contains(err.Error(), "empty but required") {
			t.Errorf("expected error message to contain 'empty but required', got %q", err.Error())
		}
	})

	t.Run("Empty checksum", func(t *testing.T) {
		algorithm := "md5"

		result, err := EncodeContentDigestHeader(algorithm, "")
		if err == nil {
			t.Fatal("expected error for empty checksum")
		}
		if result != "" {
			t.Errorf("expected empty result, got %q", result)
		}
		if !strings.Contains(err.Error(), "empty but required") {
			t.Errorf("expected error message to contain 'empty but required', got %q", err.Error())
		}
	})

	t.Run("Both empty", func(t *testing.T) {
		result, err := EncodeContentDigestHeader("", "")
		if err == nil {
			t.Fatal("expected error for empty inputs")
		}
		if result != "" {
			t.Errorf("expected empty result, got %q", result)
		}
		if !strings.Contains(err.Error(), "empty but required") {
			t.Errorf("expected error message to contain 'empty but required', got %q", err.Error())
		}
	})

	t.Run("Invalid hex checksum", func(t *testing.T) {
		testCases := []struct {
			name      string
			algorithm string
			checksum  string
		}{
			{"Invalid hex characters", "md5", "ghijklmnop"},
			{"Odd length hex", "md5", "abc"},
			{"Non-hex characters", "md5", "xyz123"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := EncodeContentDigestHeader(tc.algorithm, tc.checksum)
				if err == nil {
					t.Fatal("expected error for invalid hex checksum")
				}
				if result != "" {
					t.Errorf("expected empty result, got %q", result)
				}
				if !strings.Contains(err.Error(), "failed to decode hex") {
					t.Errorf("expected error message to contain 'failed to decode hex', got %q", err.Error())
				}
			})
		}
	})
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	testCases := []struct {
		name      string
		algorithm string
		checksum  string
	}{
		{"MD5", "md5", "098f6bcd4621d373cade4e832627b4f6"},
		{"SHA256", "sha256", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
		{"SHA1", "sha1", "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"},
		{"Uppercase algorithm", "MD5", "098f6bcd4621d373cade4e832627b4f6"},
		{"Mixed case algorithm", "Sha256", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode
			encoded, err := EncodeContentDigestHeader(tc.algorithm, tc.checksum)
			if err != nil {
				t.Fatalf("unexpected error during encoding: %v", err)
			}

			// Decode
			decodedAlgorithm, decodedChecksum, err := DecodeContentDigestHeader(encoded)
			if err != nil {
				t.Fatalf("unexpected error during decoding: %v", err)
			}

			// Verify algorithm (should be lowercased)
			expectedAlgorithm := strings.ToLower(tc.algorithm)
			if decodedAlgorithm != expectedAlgorithm {
				t.Errorf("expected algorithm %q, got %q", expectedAlgorithm, decodedAlgorithm)
			}

			// Verify checksum
			if decodedChecksum != tc.checksum {
				t.Errorf("expected checksum %q, got %q", tc.checksum, decodedChecksum)
			}
		})
	}
}

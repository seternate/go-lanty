package header

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func DecodeContentDigestHeader(digestValue string) (algorithm string, checksum string, err error) {
	if digestValue == "" {
		return "", "", fmt.Errorf("digest value is empty")
	}

	parts := strings.SplitN(digestValue, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid digest value format: expected algorithm=checksum, got %s", digestValue)
	}

	algorithm = strings.ToLower(strings.TrimSpace(parts[0]))
	checksumValue := strings.TrimSpace(parts[1])

	// Decode base64 checksum (RFC 9530 standard format)
	checksumBytes, err := base64.StdEncoding.DecodeString(checksumValue)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode base64 checksum=%s: %w", checksumValue, err)
	}

	// Convert bytes to hex string
	checksum = fmt.Sprintf("%x", checksumBytes)
	return algorithm, checksum, nil
}

func EncodeContentDigestHeader(algorithm string, checksum string) (string, error) {
	if algorithm == "" || checksum == "" {
		return "", fmt.Errorf("algorithm and checksum are empty but required")
	}

	// Convert hex string to bytes
	checksumBytes, err := hex.DecodeString(checksum)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex checksum=%s: %w", checksum, err)
	}

	// Encode to base64
	checksumBase64 := base64.StdEncoding.EncodeToString(checksumBytes)

	// Format: algorithm=base64_checksum
	algorithmLower := strings.ToLower(algorithm)
	return fmt.Sprintf("%s=%s", algorithmLower, checksumBase64), nil
}

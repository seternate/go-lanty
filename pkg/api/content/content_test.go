package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultHashCalculator(t *testing.T) {
	calculator := NewDefaultHashCalculator()
	assert.NotNil(t, calculator)
}

func TestNewDefaultMimeTypeDetector(t *testing.T) {
	detector := NewDefaultMimeTypeDetector()
	assert.NotNil(t, detector)
}

func TestNewDefaultHashCalculator_ImplementsInterface(t *testing.T) {
	calculator := NewDefaultHashCalculator()
	assert.NotNil(t, calculator)
	
	// Verify it implements the HashCalculator interface
	var _ HashCalculator = calculator
}

func TestNewDefaultMimeTypeDetector_ImplementsInterface(t *testing.T) {
	detector := NewDefaultMimeTypeDetector()
	assert.NotNil(t, detector)
	
	// Verify it implements the MimeTypeDetector interface
	var _ MimeTypeDetector = detector
}

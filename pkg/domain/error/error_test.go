package domainerr

import "testing"

func TestErrorCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected string
	}{
		{"Validation", ErrorCodeValidation, "validation"},
		{"InvariantViolation", ErrorCodeInvariantViolation, "invariant_violation"},
		{"TrustedInvariantViolation", ErrorCodeTrustedInvariantViolation, "trusted_invariant_violation"},
		{"Conflict", ErrorCodeConflict, "conflict"},
		{"NotFound", ErrorCodeNotFound, "not_found"},
		{"Internal", ErrorCodeInternal, "internal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.code) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.code))
			}
		})
	}
}

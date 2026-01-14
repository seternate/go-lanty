package game

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

func TestGameAssetRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     GameAssetRole
		expected string
	}{
		{"Icon", GAME_ASSET_ROLE_ICON, "icon"},
		{"Blob", GAME_ASSET_ROLE_BLOB, "blob"},
		{"Undefined", GAME_ASSET_ROLE_UNDEFINED, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseGameAssetRole(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		expected    GameAssetRole
		expectError bool
	}{
		{"Valid icon", "icon", GAME_ASSET_ROLE_ICON, false},
		{"Valid blob", "blob", GAME_ASSET_ROLE_BLOB, false},
		{"Invalid role", "invalid", GAME_ASSET_ROLE_UNDEFINED, true},
		{"Empty string", "", GAME_ASSET_ROLE_UNDEFINED, true},
		{"Case sensitive - uppercase", "ICON", GAME_ASSET_ROLE_UNDEFINED, true},
		{"Case sensitive - mixed", "Icon", GAME_ASSET_ROLE_UNDEFINED, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGameAssetRole(tt.role)
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
					t.Errorf("ParseGameAssetRole() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestNewGameAsset(t *testing.T) {
	testID := uuid.New()

	t.Run("Valid icon asset", func(t *testing.T) {
		asset, err := NewGameAsset(testID, "icon", "image/png")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if asset == nil {
			t.Fatal("expected asset, got nil")
		}
		if asset.AssetID != testID {
			t.Errorf("expected AssetID %v, got %v", testID, asset.AssetID)
		}
		if asset.Role != GAME_ASSET_ROLE_ICON {
			t.Errorf("expected Role %v, got %v", GAME_ASSET_ROLE_ICON, asset.Role)
		}
	})

	t.Run("Valid blob asset", func(t *testing.T) {
		asset, err := NewGameAsset(testID, "blob", "application/octet-stream")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if asset == nil {
			t.Fatal("expected asset, got nil")
		}
		if asset.Role != GAME_ASSET_ROLE_BLOB {
			t.Errorf("expected Role %v, got %v", GAME_ASSET_ROLE_BLOB, asset.Role)
		}
	})

	t.Run("Invalid role", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		asset, err := NewGameAsset(testID, "invalid", "image/png")
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

	t.Run("Icon with non-image mime type", func(t *testing.T) {
		asset, err := NewGameAsset(testID, "icon", "application/octet-stream")
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

	t.Run("Icon with valid image mime types", func(t *testing.T) {
		validMimeTypes := []string{"image/png", "image/jpeg", "image/gif", "image/webp", "image/svg+xml"}
		for _, mimeType := range validMimeTypes {
			t.Run(mimeType, func(t *testing.T) {
				asset, err := NewGameAsset(testID, "icon", mimeType)
				if err != nil {
					t.Errorf("unexpected error for %q: %v", mimeType, err)
				}
				if asset == nil {
					t.Errorf("expected asset for %q, got nil", mimeType)
				}
			})
		}
	})

	t.Run("Blob with any mime type", func(t *testing.T) {
		mimeTypes := []string{"application/octet-stream", "text/plain", "image/png"}
		for _, mimeType := range mimeTypes {
			t.Run(mimeType, func(t *testing.T) {
				asset, err := NewGameAsset(testID, "blob", mimeType)
				if err != nil {
					t.Errorf("unexpected error for %q: %v", mimeType, err)
				}
				if asset == nil {
					t.Errorf("expected asset for %q, got nil", mimeType)
				}
			})
		}
	})
}

func TestRehydrateGameAsset(t *testing.T) {
	testID := uuid.New()

	t.Run("Valid asset", func(t *testing.T) {
		asset, err := RehydrateGameAsset(testID, "icon")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if asset == nil {
			t.Fatal("expected asset, got nil")
		}
		if asset.AssetID != testID {
			t.Errorf("expected AssetID %v, got %v", testID, asset.AssetID)
		}
		if asset.Role != GAME_ASSET_ROLE_ICON {
			t.Errorf("expected Role %v, got %v", GAME_ASSET_ROLE_ICON, asset.Role)
		}
	})

	t.Run("Invalid role returns TrustedInvariantViolationError", func(t *testing.T) {
		asset, err := RehydrateGameAsset(testID, "invalid")
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

func TestValidateMimeTypeForRole(t *testing.T) {
	t.Run("Icon with image mime type", func(t *testing.T) {
		err := validateMimeTypeForRole(GAME_ASSET_ROLE_ICON, "image/png")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Icon with non-image mime type", func(t *testing.T) {
		err := validateMimeTypeForRole(GAME_ASSET_ROLE_ICON, "application/octet-stream")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})

	t.Run("Blob with any mime type", func(t *testing.T) {
		mimeTypes := []string{"application/octet-stream", "text/plain", "image/png"}
		for _, mimeType := range mimeTypes {
			err := validateMimeTypeForRole(GAME_ASSET_ROLE_BLOB, mimeType)
			if err != nil {
				t.Errorf("unexpected error for %q: %v", mimeType, err)
			}
		}
	})

	t.Run("Undefined role with any mime type", func(t *testing.T) {
		err := validateMimeTypeForRole(GAME_ASSET_ROLE_UNDEFINED, "application/octet-stream")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

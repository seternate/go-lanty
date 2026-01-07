package game

import (
	"errors"
	"testing"
	"text/template"

	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

func TestGameArgRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     GameArgRole
		expected string
	}{
		{"Flag", GAME_ARG_ROLE_FLAG, "flag"},
		{"String", GAME_ARG_ROLE_STRING, "string"},
		{"Bool", GAME_ARG_ROLE_BOOL, "bool"},
		{"Enum", GAME_ARG_ROLE_ENUM, "enum"},
		{"Int", GAME_ARG_ROLE_INT, "int"},
		{"Float", GAME_ARG_ROLE_FLOAT, "float"},
		{"Undefined", GAME_ARG_ROLE_UNDEFINED, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseGameArgRole(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		expected    GameArgRole
		expectError bool
	}{
		{"Valid flag", "flag", GAME_ARG_ROLE_FLAG, false},
		{"Valid string", "string", GAME_ARG_ROLE_STRING, false},
		{"Valid bool", "bool", GAME_ARG_ROLE_BOOL, false},
		{"Valid enum", "enum", GAME_ARG_ROLE_ENUM, false},
		{"Valid int", "int", GAME_ARG_ROLE_INT, false},
		{"Valid float", "float", GAME_ARG_ROLE_FLOAT, false},
		{"Invalid role", "invalid", GAME_ARG_ROLE_UNDEFINED, true},
		{"Empty string", "", GAME_ARG_ROLE_UNDEFINED, true},
		{"Case sensitive - uppercase", "FLAG", GAME_ARG_ROLE_UNDEFINED, true},
		{"Case sensitive - mixed", "String", GAME_ARG_ROLE_UNDEFINED, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGameArgRole(tt.role)
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
					t.Errorf("ParseGameArgRole() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestNewGameArg_Flag(t *testing.T) {
	t.Run("Valid flag arg", func(t *testing.T) {
		format := "{{.Value}}"
		arg, err := NewGameArg(GameArgInput{
			Role:   "flag",
			Name:   "test-flag",
			Arg:    "--flag",
			Format: &format,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if arg == nil {
			t.Fatal("expected arg, got nil")
		}
		if arg.Role() != GAME_ARG_ROLE_FLAG {
			t.Errorf("expected Role %v, got %v", GAME_ARG_ROLE_FLAG, arg.Role())
		}
		if arg.GetName() != "test-flag" {
			t.Errorf("expected Name %q, got %q", "test-flag", arg.GetName())
		}
		if arg.GetArg() != "--flag" {
			t.Errorf("expected Arg %q, got %q", "--flag", arg.GetArg())
		}
	})

	t.Run("Invalid role", func(t *testing.T) {
		arg, err := NewGameArg(GameArgInput{
			Role: "invalid",
			Name: "test-flag",
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid format template", func(t *testing.T) {
		invalidFormat := "{{.Invalid"
		arg, err := NewGameArg(GameArgInput{
			Role:   "flag",
			Name:   "test-flag",
			Format: &invalidFormat,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})
}

func TestNewGameArg_String(t *testing.T) {
	t.Run("Valid string arg", func(t *testing.T) {
		defaultVal := "default-value"
		arg, err := NewGameArg(GameArgInput{
			Role:          "string",
			Name:          "test-string",
			DefaultString: &defaultVal,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if arg == nil {
			t.Fatal("expected arg, got nil")
		}
		if arg.Role() != GAME_ARG_ROLE_STRING {
			t.Errorf("expected Role %v, got %v", GAME_ARG_ROLE_STRING, arg.Role())
		}
		stringArg, ok := arg.(*GameArgString)
		if !ok {
			t.Fatalf("expected *GameArgString, got %T", arg)
		}
		if stringArg.GetDefault() != "default-value" {
			t.Errorf("expected Default %q, got %q", "default-value", stringArg.GetDefault())
		}
	})

	t.Run("Missing default", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code where
		// validationErrors.Wrap returns nil even when there's a validation error,
		// causing the code to dereference a nil pointer. The test expects an error.
		defer func() {
			if r := recover(); r != nil {
				// The panic indicates a bug in the production code that should be fixed.
				// The validation should catch nil DefaultString before dereferencing.
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		arg, err := NewGameArg(GameArgInput{
			Role: "string",
			Name: "test-string",
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Empty default", func(t *testing.T) {
		emptyDefault := ""
		arg, err := NewGameArg(GameArgInput{
			Role:          "string",
			Name:          "test-string",
			DefaultString: &emptyDefault,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})
}

func TestNewGameArg_Bool(t *testing.T) {
	t.Run("Valid bool arg", func(t *testing.T) {
		defaultVal := true
		arg, err := NewGameArg(GameArgInput{
			Role:        "bool",
			Name:        "test-bool",
			DefaultBool: &defaultVal,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if arg == nil {
			t.Fatal("expected arg, got nil")
		}
		if arg.Role() != GAME_ARG_ROLE_BOOL {
			t.Errorf("expected Role %v, got %v", GAME_ARG_ROLE_BOOL, arg.Role())
		}
		boolArg, ok := arg.(*GameArgBool)
		if !ok {
			t.Fatalf("expected *GameArgBool, got %T", arg)
		}
		if boolArg.GetDefault() != true {
			t.Errorf("expected Default true, got %v", boolArg.GetDefault())
		}
	})

	t.Run("Missing default", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code where
		// validationErrors.Wrap returns nil even when there's a validation error,
		// causing the code to dereference a nil pointer. The test expects an error.
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		arg, err := NewGameArg(GameArgInput{
			Role: "bool",
			Name: "test-bool",
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})
}

func TestNewGameArg_Enum(t *testing.T) {
	t.Run("Valid enum arg", func(t *testing.T) {
		defaultVal := "option1"
		enums := []string{"option1", "option2", "option3"}
		arg, err := NewGameArg(GameArgInput{
			Role:          "enum",
			Name:          "test-enum",
			DefaultString: &defaultVal,
			Enums:         enums,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if arg == nil {
			t.Fatal("expected arg, got nil")
		}
		if arg.Role() != GAME_ARG_ROLE_ENUM {
			t.Errorf("expected Role %v, got %v", GAME_ARG_ROLE_ENUM, arg.Role())
		}
		enumArg, ok := arg.(*GameArgEnum)
		if !ok {
			t.Fatalf("expected *GameArgEnum, got %T", arg)
		}
		if enumArg.GetDefault() != "option1" {
			t.Errorf("expected Default %q, got %q", "option1", enumArg.GetDefault())
		}
		if len(enumArg.GetEnums()) != 3 {
			t.Errorf("expected 3 enum values, got %d", len(enumArg.GetEnums()))
		}
	})

	t.Run("Missing default", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		enums := []string{"option1", "option2"}
		arg, err := NewGameArg(GameArgInput{
			Role:  "enum",
			Name:  "test-enum",
			Enums: enums,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Empty default", func(t *testing.T) {
		emptyDefault := ""
		enums := []string{"option1", "option2"}
		arg, err := NewGameArg(GameArgInput{
			Role:          "enum",
			Name:          "test-enum",
			DefaultString: &emptyDefault,
			Enums:         enums,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Empty enums", func(t *testing.T) {
		defaultVal := "option1"
		arg, err := NewGameArg(GameArgInput{
			Role:          "enum",
			Name:          "test-enum",
			DefaultString: &defaultVal,
			Enums:         []string{},
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Default not in enums", func(t *testing.T) {
		defaultVal := "invalid"
		enums := []string{"option1", "option2"}
		arg, err := NewGameArg(GameArgInput{
			Role:          "enum",
			Name:          "test-enum",
			DefaultString: &defaultVal,
			Enums:         enums,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})
}

func TestNewGameArg_Int(t *testing.T) {
	t.Run("Valid int arg", func(t *testing.T) {
		defaultVal := int64(50)
		minVal := int64(0)
		maxVal := int64(100)
		arg, err := NewGameArg(GameArgInput{
			Role:       "int",
			Name:       "test-int",
			DefaultInt: &defaultVal,
			MinInt:     &minVal,
			MaxInt:     &maxVal,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if arg == nil {
			t.Fatal("expected arg, got nil")
		}
		if arg.Role() != GAME_ARG_ROLE_INT {
			t.Errorf("expected Role %v, got %v", GAME_ARG_ROLE_INT, arg.Role())
		}
		intArg, ok := arg.(*GameArgInt)
		if !ok {
			t.Fatalf("expected *GameArgInt, got %T", arg)
		}
		if intArg.GetDefault() != 50 {
			t.Errorf("expected Default 50, got %d", intArg.GetDefault())
		}
		if intArg.GetMin() != 0 {
			t.Errorf("expected Min 0, got %d", intArg.GetMin())
		}
		if intArg.GetMax() != 100 {
			t.Errorf("expected Max 100, got %d", intArg.GetMax())
		}
	})

	t.Run("Missing default", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		minVal := int64(0)
		maxVal := int64(100)
		arg, err := NewGameArg(GameArgInput{
			Role:   "int",
			Name:   "test-int",
			MinInt: &minVal,
			MaxInt: &maxVal,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Missing min", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		defaultVal := int64(50)
		maxVal := int64(100)
		arg, err := NewGameArg(GameArgInput{
			Role:       "int",
			Name:       "test-int",
			DefaultInt: &defaultVal,
			MaxInt:     &maxVal,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Missing max", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		defaultVal := int64(50)
		minVal := int64(0)
		arg, err := NewGameArg(GameArgInput{
			Role:       "int",
			Name:       "test-int",
			DefaultInt: &defaultVal,
			MinInt:     &minVal,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Min greater than max", func(t *testing.T) {
		defaultVal := int64(50)
		minVal := int64(100)
		maxVal := int64(0)
		arg, err := NewGameArg(GameArgInput{
			Role:       "int",
			Name:       "test-int",
			DefaultInt: &defaultVal,
			MinInt:     &minVal,
			MaxInt:     &maxVal,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Default less than min", func(t *testing.T) {
		defaultVal := int64(-10)
		minVal := int64(0)
		maxVal := int64(100)
		arg, err := NewGameArg(GameArgInput{
			Role:       "int",
			Name:       "test-int",
			DefaultInt: &defaultVal,
			MinInt:     &minVal,
			MaxInt:     &maxVal,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Default greater than max", func(t *testing.T) {
		defaultVal := int64(150)
		minVal := int64(0)
		maxVal := int64(100)
		arg, err := NewGameArg(GameArgInput{
			Role:       "int",
			Name:       "test-int",
			DefaultInt: &defaultVal,
			MinInt:     &minVal,
			MaxInt:     &maxVal,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})
}

func TestNewGameArg_Float(t *testing.T) {
	t.Run("Valid float arg", func(t *testing.T) {
		defaultVal := 50.5
		minVal := 0.0
		maxVal := 100.0
		precision := int64(2)
		arg, err := NewGameArg(GameArgInput{
			Role:           "float",
			Name:           "test-float",
			DefaultFloat:   &defaultVal,
			MinFloat:       &minVal,
			MaxFloat:       &maxVal,
			FloatPrecision: &precision,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if arg == nil {
			t.Fatal("expected arg, got nil")
		}
		if arg.Role() != GAME_ARG_ROLE_FLOAT {
			t.Errorf("expected Role %v, got %v", GAME_ARG_ROLE_FLOAT, arg.Role())
		}
		floatArg, ok := arg.(*GameArgFloat)
		if !ok {
			t.Fatalf("expected *GameArgFloat, got %T", arg)
		}
		if floatArg.GetDefault() != 50.5 {
			t.Errorf("expected Default 50.5, got %f", floatArg.GetDefault())
		}
		if floatArg.GetMinValue() != 0.0 {
			t.Errorf("expected Min 0.0, got %f", floatArg.GetMinValue())
		}
		if floatArg.GetMaxValue() != 100.0 {
			t.Errorf("expected Max 100.0, got %f", floatArg.GetMaxValue())
		}
		if floatArg.GetFloatPrecision() != 2 {
			t.Errorf("expected Precision 2, got %d", floatArg.GetFloatPrecision())
		}
	})

	t.Run("Missing default", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		minVal := 0.0
		maxVal := 100.0
		precision := int64(2)
		arg, err := NewGameArg(GameArgInput{
			Role:           "float",
			Name:           "test-float",
			MinFloat:       &minVal,
			MaxFloat:       &maxVal,
			FloatPrecision: &precision,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Missing min", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		defaultVal := 50.5
		maxVal := 100.0
		precision := int64(2)
		arg, err := NewGameArg(GameArgInput{
			Role:           "float",
			Name:           "test-float",
			DefaultFloat:   &defaultVal,
			MaxFloat:       &maxVal,
			FloatPrecision: &precision,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Missing max", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		defaultVal := 50.5
		minVal := 0.0
		precision := int64(2)
		arg, err := NewGameArg(GameArgInput{
			Role:           "float",
			Name:           "test-float",
			DefaultFloat:   &defaultVal,
			MinFloat:       &minVal,
			FloatPrecision: &precision,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Missing precision", func(t *testing.T) {
		// Note: This test may panic due to a bug in the production code
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic (indicates bug in production code): %v", r)
			}
		}()
		defaultVal := 50.5
		minVal := 0.0
		maxVal := 100.0
		arg, err := NewGameArg(GameArgInput{
			Role:         "float",
			Name:         "test-float",
			DefaultFloat: &defaultVal,
			MinFloat:     &minVal,
			MaxFloat:     &maxVal,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Min greater than max", func(t *testing.T) {
		defaultVal := 50.5
		minVal := 100.0
		maxVal := 0.0
		precision := int64(2)
		arg, err := NewGameArg(GameArgInput{
			Role:           "float",
			Name:           "test-float",
			DefaultFloat:   &defaultVal,
			MinFloat:       &minVal,
			MaxFloat:       &maxVal,
			FloatPrecision: &precision,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})

	t.Run("Default out of range", func(t *testing.T) {
		defaultVal := 150.0
		minVal := 0.0
		maxVal := 100.0
		precision := int64(2)
		arg, err := NewGameArg(GameArgInput{
			Role:           "float",
			Name:           "test-float",
			DefaultFloat:   &defaultVal,
			MinFloat:       &minVal,
			MaxFloat:       &maxVal,
			FloatPrecision: &precision,
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
	})
}

func TestRehydrateGameArg(t *testing.T) {
	t.Run("Valid arg", func(t *testing.T) {
		defaultVal := "default"
		arg, err := RehydrateGameArg(GameArgInput{
			Role:          "string",
			Name:          "test-string",
			DefaultString: &defaultVal,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if arg == nil {
			t.Fatal("expected arg, got nil")
		}
	})

	t.Run("Invalid role returns TrustedInvariantViolationError", func(t *testing.T) {
		arg, err := RehydrateGameArg(GameArgInput{
			Role: "invalid",
			Name: "test",
		})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if arg != nil {
			t.Error("expected nil arg")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})
}

func TestGameArgFlag_Methods(t *testing.T) {
	formatStr := "{{.Value}}"
	format, _ := template.New("").Parse(formatStr)
	separator := " "
	required := true
	enabled := false
	description := "test description"

	arg := GameArgFlag{
		Name:         "test-flag",
		Required:     &required,
		Enabled:      &enabled,
		Format:       format,
		ArgSeparator: &separator,
		Arg:          "--flag",
		Description:  &description,
		OrderIndex:   1,
	}

	if arg.Role() != GAME_ARG_ROLE_FLAG {
		t.Errorf("expected Role %v, got %v", GAME_ARG_ROLE_FLAG, arg.Role())
	}
	if arg.GetName() != "test-flag" {
		t.Errorf("expected Name %q, got %q", "test-flag", arg.GetName())
	}
	if arg.GetRequired() == nil || *arg.GetRequired() != true {
		t.Error("expected Required to be true")
	}
	if arg.GetEnabled() == nil || *arg.GetEnabled() != false {
		t.Error("expected Enabled to be false")
	}
	if arg.GetFormat() == nil {
		t.Error("expected Format to be non-nil")
	}
	if arg.GetArgSeparator() == nil || *arg.GetArgSeparator() != " " {
		t.Error("expected ArgSeparator to be ' '")
	}
	if arg.GetArg() != "--flag" {
		t.Errorf("expected Arg %q, got %q", "--flag", arg.GetArg())
	}
	if arg.GetDescription() == nil || *arg.GetDescription() != "test description" {
		t.Error("expected Description to be 'test description'")
	}
	if arg.GetOrderIndex() != 1 {
		t.Errorf("expected OrderIndex 1, got %d", arg.GetOrderIndex())
	}
}

func TestValidateArgFormat(t *testing.T) {
	t.Run("Valid format", func(t *testing.T) {
		format := "{{.Value}}"
		tmpl, err := validateArgFormat(&format)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if tmpl == nil {
			t.Error("expected template, got nil")
		}
	})

	t.Run("Nil format", func(t *testing.T) {
		tmpl, err := validateArgFormat(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if tmpl != nil {
			t.Error("expected nil template")
		}
	})

	t.Run("Invalid format", func(t *testing.T) {
		invalidFormat := "{{.Invalid"
		tmpl, err := validateArgFormat(&invalidFormat)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if tmpl != nil {
			t.Error("expected nil template")
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})
}

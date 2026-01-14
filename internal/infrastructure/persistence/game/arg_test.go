package game

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/internal/domain/game"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

func TestGameArgRow_Assemble(t *testing.T) {
	execID := uuid.New()

	t.Run("Valid flag arg", func(t *testing.T) {
		row := GameArgRow{
			GameExecID:     execID,
			Role:           "flag",
			Name:           "test-flag",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--flag",
			Description:    nil,
			DefaultString:  nil,
			DefaultBool:    nil,
			DefaultInt:     nil,
			DefaultFloat:   nil,
			Enums:          nil,
			MinInt:         nil,
			MaxInt:         nil,
			MinFloat:       nil,
			MaxFloat:       nil,
			FloatPrecision: nil,
			OrderIndex:     0,
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled == nil {
			t.Fatal("expected arg, got nil")
		}
		if assembled.Role() != game.GAME_ARG_ROLE_FLAG {
			t.Errorf("expected Role %v, got %v", game.GAME_ARG_ROLE_FLAG, assembled.Role())
		}
		if assembled.GetName() != "test-flag" {
			t.Errorf("expected Name %q, got %q", "test-flag", assembled.GetName())
		}
		if assembled.GetArg() != "--flag" {
			t.Errorf("expected Arg %q, got %q", "--flag", assembled.GetArg())
		}
	})

	t.Run("Valid string arg", func(t *testing.T) {
		defaultStr := "default-value"
		row := GameArgRow{
			GameExecID:     execID,
			Role:           "string",
			Name:           "test-string",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--string",
			Description:    nil,
			DefaultString:  &defaultStr,
			DefaultBool:    nil,
			DefaultInt:     nil,
			DefaultFloat:   nil,
			Enums:          nil,
			MinInt:         nil,
			MaxInt:         nil,
			MinFloat:       nil,
			MaxFloat:       nil,
			FloatPrecision: nil,
			OrderIndex:     1,
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled.Role() != game.GAME_ARG_ROLE_STRING {
			t.Errorf("expected Role %v, got %v", game.GAME_ARG_ROLE_STRING, assembled.Role())
		}
		if strArg, ok := assembled.(*game.GameArgString); ok {
			if strArg.GetDefault() != defaultStr {
				t.Errorf("expected Default %q, got %q", defaultStr, strArg.GetDefault())
			}
		} else {
			t.Error("expected GameArgString type")
		}
	})

	t.Run("Valid bool arg", func(t *testing.T) {
		defaultBool := true
		row := GameArgRow{
			GameExecID:     execID,
			Role:           "bool",
			Name:           "test-bool",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--bool",
			Description:    nil,
			DefaultString:  nil,
			DefaultBool:    &defaultBool,
			DefaultInt:     nil,
			DefaultFloat:   nil,
			Enums:          nil,
			MinInt:         nil,
			MaxInt:         nil,
			MinFloat:       nil,
			MaxFloat:       nil,
			FloatPrecision: nil,
			OrderIndex:     2,
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled.Role() != game.GAME_ARG_ROLE_BOOL {
			t.Errorf("expected Role %v, got %v", game.GAME_ARG_ROLE_BOOL, assembled.Role())
		}
		if boolArg, ok := assembled.(*game.GameArgBool); ok {
			if boolArg.Default != defaultBool {
				t.Errorf("expected Default %v, got %v", defaultBool, boolArg.Default)
			}
		} else {
			t.Error("expected GameArgBool type")
		}
	})

	t.Run("Valid enum arg", func(t *testing.T) {
		defaultStr := "option1"
		enums := []string{"option1", "option2", "option3"}
		row := GameArgRow{
			GameExecID:     execID,
			Role:           "enum",
			Name:           "test-enum",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--enum",
			Description:    nil,
			DefaultString:  &defaultStr,
			DefaultBool:    nil,
			DefaultInt:     nil,
			DefaultFloat:   nil,
			Enums:          enums,
			MinInt:         nil,
			MaxInt:         nil,
			MinFloat:       nil,
			MaxFloat:       nil,
			FloatPrecision: nil,
			OrderIndex:     3,
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled.Role() != game.GAME_ARG_ROLE_ENUM {
			t.Errorf("expected Role %v, got %v", game.GAME_ARG_ROLE_ENUM, assembled.Role())
		}
		if enumArg, ok := assembled.(*game.GameArgEnum); ok {
			if enumArg.GetDefault() != defaultStr {
				t.Errorf("expected Default %q, got %q", defaultStr, enumArg.GetDefault())
			}
			if len(enumArg.Values) != len(enums) {
				t.Errorf("expected %d enum values, got %d", len(enums), len(enumArg.Values))
			}
		} else {
			t.Error("expected GameArgEnum type")
		}
	})

	t.Run("Valid int arg", func(t *testing.T) {
		defaultInt := int64(42)
		minInt := int64(0)
		maxInt := int64(100)
		row := GameArgRow{
			GameExecID:     execID,
			Role:           "int",
			Name:           "test-int",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--int",
			Description:    nil,
			DefaultString:  nil,
			DefaultBool:    nil,
			DefaultInt:     &defaultInt,
			DefaultFloat:   nil,
			Enums:          nil,
			MinInt:         &minInt,
			MaxInt:         &maxInt,
			MinFloat:       nil,
			MaxFloat:       nil,
			FloatPrecision: nil,
			OrderIndex:     4,
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled.Role() != game.GAME_ARG_ROLE_INT {
			t.Errorf("expected Role %v, got %v", game.GAME_ARG_ROLE_INT, assembled.Role())
		}
		if intArg, ok := assembled.(*game.GameArgInt); ok {
			if intArg.Default != defaultInt {
				t.Errorf("expected Default %d, got %d", defaultInt, intArg.Default)
			}
			if intArg.Min != minInt {
				t.Errorf("expected Min %d, got %d", minInt, intArg.Min)
			}
			if intArg.Max != maxInt {
				t.Errorf("expected Max %d, got %d", maxInt, intArg.Max)
			}
		} else {
			t.Error("expected GameArgInt type")
		}
	})

	t.Run("Valid float arg", func(t *testing.T) {
		defaultFloat := 3.14
		minFloat := 0.0
		maxFloat := 10.0
		floatPrecision := int64(2)
		row := GameArgRow{
			GameExecID:     execID,
			Role:           "float",
			Name:           "test-float",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--float",
			Description:    nil,
			DefaultString:  nil,
			DefaultBool:    nil,
			DefaultInt:     nil,
			DefaultFloat:   &defaultFloat,
			Enums:          nil,
			MinInt:         nil,
			MaxInt:         nil,
			MinFloat:       &minFloat,
			MaxFloat:       &maxFloat,
			FloatPrecision: &floatPrecision,
			OrderIndex:     5,
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled.Role() != game.GAME_ARG_ROLE_FLOAT {
			t.Errorf("expected Role %v, got %v", game.GAME_ARG_ROLE_FLOAT, assembled.Role())
		}
		if floatArg, ok := assembled.(*game.GameArgFloat); ok {
			if floatArg.Default != defaultFloat {
				t.Errorf("expected Default %f, got %f", defaultFloat, floatArg.Default)
			}
			if floatArg.Min != minFloat {
				t.Errorf("expected Min %f, got %f", minFloat, floatArg.Min)
			}
			if floatArg.Max != maxFloat {
				t.Errorf("expected Max %f, got %f", maxFloat, floatArg.Max)
			}
			if floatArg.FloatPrecision != floatPrecision {
				t.Errorf("expected FloatPrecision %d, got %d", floatPrecision, floatArg.FloatPrecision)
			}
		} else {
			t.Error("expected GameArgFloat type")
		}
	})

	t.Run("Invalid role", func(t *testing.T) {
		row := GameArgRow{
			GameExecID:     execID,
			Role:           "invalid-role",
			Name:           "test-arg",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--arg",
			Description:    nil,
			DefaultString:  nil,
			DefaultBool:    nil,
			DefaultInt:     nil,
			DefaultFloat:   nil,
			Enums:          nil,
			MinInt:         nil,
			MaxInt:         nil,
			MinFloat:       nil,
			MaxFloat:       nil,
			FloatPrecision: nil,
			OrderIndex:     0,
		}

		assembled, err := row.Assemble()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil arg")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})
}

func TestDisassembleGameArg(t *testing.T) {
	execID := uuid.New()

	t.Run("Disassemble flag arg", func(t *testing.T) {
		arg, err := game.NewGameArg(game.GameArgInput{
			Role:       "flag",
			Name:       "test-flag",
			Arg:        "--flag",
			OrderIndex: 0,
		})
		if err != nil {
			t.Fatalf("failed to create test arg: %v", err)
		}

		row := disassembleGameArg(execID, arg)
		if row == nil {
			t.Fatal("expected row, got nil")
		}
		if row.GameExecID != execID {
			t.Errorf("expected GameExecID %v, got %v", execID, row.GameExecID)
		}
		if row.Role != "flag" {
			t.Errorf("expected Role %q, got %q", "flag", row.Role)
		}
		if row.Name != "test-flag" {
			t.Errorf("expected Name %q, got %q", "test-flag", row.Name)
		}
	})

	t.Run("Disassemble string arg", func(t *testing.T) {
		defaultStr := "default"
		arg, err := game.NewGameArg(game.GameArgInput{
			Role:          "string",
			Name:          "test-string",
			Arg:           "--string",
			DefaultString: &defaultStr,
			OrderIndex:    1,
		})
		if err != nil {
			t.Fatalf("failed to create test arg: %v", err)
		}

		row := disassembleGameArg(execID, arg)
		if row.DefaultString == nil || *row.DefaultString != defaultStr {
			t.Errorf("expected DefaultString %q, got %v", defaultStr, row.DefaultString)
		}
	})

	t.Run("Disassemble bool arg", func(t *testing.T) {
		defaultBool := true
		arg, err := game.NewGameArg(game.GameArgInput{
			Role:        "bool",
			Name:        "test-bool",
			Arg:         "--bool",
			DefaultBool: &defaultBool,
			OrderIndex:  2,
		})
		if err != nil {
			t.Fatalf("failed to create test arg: %v", err)
		}

		row := disassembleGameArg(execID, arg)
		if row.DefaultBool == nil || *row.DefaultBool != defaultBool {
			t.Errorf("expected DefaultBool %v, got %v", defaultBool, row.DefaultBool)
		}
	})

	t.Run("Round trip: Disassemble then Assemble", func(t *testing.T) {
		originalArg, err := game.NewGameArg(game.GameArgInput{
			Role:       "flag",
			Name:       "test-flag",
			Arg:        "--flag",
			OrderIndex: 0,
		})
		if err != nil {
			t.Fatalf("failed to create test arg: %v", err)
		}

		row := disassembleGameArg(execID, originalArg)
		reassembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("failed to reassemble: %v", err)
		}

		if reassembled.Role() != originalArg.Role() {
			t.Errorf("expected Role %v, got %v", originalArg.Role(), reassembled.Role())
		}
		if reassembled.GetName() != originalArg.GetName() {
			t.Errorf("expected Name %q, got %q", originalArg.GetName(), reassembled.GetName())
		}
		if reassembled.GetArg() != originalArg.GetArg() {
			t.Errorf("expected Arg %q, got %q", originalArg.GetArg(), reassembled.GetArg())
		}
	})
}

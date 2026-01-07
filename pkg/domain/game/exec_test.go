package game

import (
	"errors"
	"testing"

	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

func TestGameExecRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     GameExecRole
		expected string
	}{
		{"Client", GAME_EXEC_ROLE_CLIENT, "client"},
		{"Server", GAME_EXEC_ROLE_SERVER, "server"},
		{"Undefined", GAME_EXEC_ROLE_UNDEFINED, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseGameExecRole(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		expected    GameExecRole
		expectError bool
	}{
		{"Valid client", "client", GAME_EXEC_ROLE_CLIENT, false},
		{"Valid server", "server", GAME_EXEC_ROLE_SERVER, false},
		{"Invalid role", "invalid", GAME_EXEC_ROLE_UNDEFINED, true},
		{"Empty string", "", GAME_EXEC_ROLE_UNDEFINED, true},
		{"Case sensitive - uppercase", "CLIENT", GAME_EXEC_ROLE_UNDEFINED, true},
		{"Case sensitive - mixed", "Client", GAME_EXEC_ROLE_UNDEFINED, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGameExecRole(tt.role)
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
					t.Errorf("ParseGameExecRole() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestNewGameExec(t *testing.T) {
	t.Run("Valid exec with minimal options", func(t *testing.T) {
		exec, err := NewGameExec("/path/to/executable", "client")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec == nil {
			t.Fatal("expected exec, got nil")
		}
		if exec.Role != GAME_EXEC_ROLE_CLIENT {
			t.Errorf("expected Role %v, got %v", GAME_EXEC_ROLE_CLIENT, exec.Role)
		}
		if exec.Path != "/path/to/executable" {
			t.Errorf("expected Path %q, got %q", "/path/to/executable", exec.Path)
		}
		if exec.RequiresAdmin != nil {
			t.Error("expected RequiresAdmin to be nil")
		}
		if exec.Format != nil {
			t.Error("expected Format to be nil")
		}
		if exec.ArgSeperator != nil {
			t.Error("expected ArgSeperator to be nil")
		}
		if len(exec.Args) != 0 {
			t.Errorf("expected empty Args, got %d", len(exec.Args))
		}
	})

	t.Run("Empty path", func(t *testing.T) {
		exec, err := NewGameExec("", "client")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if exec != nil {
			t.Error("expected nil exec")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid role", func(t *testing.T) {
		exec, err := NewGameExec("/path/to/executable", "invalid")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if exec != nil {
			t.Error("expected nil exec")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("With RequiresAdmin", func(t *testing.T) {
		requiresAdmin := true
		exec, err := NewGameExec("/path/to/executable", "client", WithRequiresAdmin(&requiresAdmin))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.RequiresAdmin == nil {
			t.Fatal("expected RequiresAdmin to be non-nil")
		}
		if *exec.RequiresAdmin != true {
			t.Error("expected RequiresAdmin to be true")
		}
	})

	t.Run("With Format", func(t *testing.T) {
		format := "{{.Value}}"
		exec, err := NewGameExec("/path/to/executable", "client", WithFormat(&format))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.Format == nil {
			t.Fatal("expected Format to be non-nil")
		}
	})

	t.Run("With invalid Format", func(t *testing.T) {
		invalidFormat := "{{.Invalid"
		exec, err := NewGameExec("/path/to/executable", "client", WithFormat(&invalidFormat))
		if err == nil {
			t.Error("expected error, got nil")
		}
		if exec != nil {
			t.Error("expected nil exec")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("With ArgSeperator", func(t *testing.T) {
		separator := " "
		exec, err := NewGameExec("/path/to/executable", "client", WithArgSeperator(&separator))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.ArgSeperator == nil {
			t.Fatal("expected ArgSeperator to be non-nil")
		}
		if *exec.ArgSeperator != " " {
			t.Errorf("expected ArgSeperator %q, got %q", " ", *exec.ArgSeperator)
		}
	})

	t.Run("With Args", func(t *testing.T) {
		defaultVal := "default"
		arg1, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg1",
			DefaultString: &defaultVal,
			OrderIndex:   1,
		})
		arg2, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg2",
			DefaultString: &defaultVal,
			OrderIndex:   2,
		})
		exec, err := NewGameExec("/path/to/executable", "client", WithArgs(arg1, arg2))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(exec.Args) != 2 {
			t.Errorf("expected 2 args, got %d", len(exec.Args))
		}
	})

	t.Run("With duplicate OrderIndex", func(t *testing.T) {
		defaultVal := "default"
		arg1, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg1",
			DefaultString: &defaultVal,
			OrderIndex:   1,
		})
		arg2, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg2",
			DefaultString: &defaultVal,
			OrderIndex:   1, // duplicate
		})
		exec, err := NewGameExec("/path/to/executable", "client", WithArgs(arg1, arg2))
		if err == nil {
			t.Error("expected error, got nil")
		}
		if exec != nil {
			t.Error("expected nil exec")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("With multiple options", func(t *testing.T) {
		requiresAdmin := true
		format := "{{.Value}}"
		separator := " "
		defaultVal := "default"
		arg, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg1",
			DefaultString: &defaultVal,
			OrderIndex:   1,
		})
		exec, err := NewGameExec("/path/to/executable", "server",
			WithRequiresAdmin(&requiresAdmin),
			WithFormat(&format),
			WithArgSeperator(&separator),
			WithArgs(arg),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec.Role != GAME_EXEC_ROLE_SERVER {
			t.Errorf("expected Role %v, got %v", GAME_EXEC_ROLE_SERVER, exec.Role)
		}
		if exec.RequiresAdmin == nil || *exec.RequiresAdmin != true {
			t.Error("expected RequiresAdmin to be true")
		}
		if exec.Format == nil {
			t.Error("expected Format to be non-nil")
		}
		if exec.ArgSeperator == nil || *exec.ArgSeperator != " " {
			t.Error("expected ArgSeperator to be ' '")
		}
		if len(exec.Args) != 1 {
			t.Errorf("expected 1 arg, got %d", len(exec.Args))
		}
	})
}

func TestRehydrateGameExec(t *testing.T) {
	t.Run("Valid exec", func(t *testing.T) {
		exec, err := RehydrateGameExec("/path/to/executable", "client")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exec == nil {
			t.Fatal("expected exec, got nil")
		}
	})

	t.Run("Invalid role returns TrustedInvariantViolationError", func(t *testing.T) {
		exec, err := RehydrateGameExec("/path/to/executable", "invalid")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if exec != nil {
			t.Error("expected nil exec")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})

	t.Run("Empty path returns TrustedInvariantViolationError", func(t *testing.T) {
		exec, err := RehydrateGameExec("", "client")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if exec != nil {
			t.Error("expected nil exec")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})
}

func TestValidatePath(t *testing.T) {
	t.Run("Valid path", func(t *testing.T) {
		err := validatePath("/path/to/executable")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Empty path", func(t *testing.T) {
		err := validatePath("")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil && validationErr.Field != "path" {
			t.Errorf("expected Field %q, got %q", "path", validationErr.Field)
		}
	})
}

func TestValidateUniqueOrderIndexArgs(t *testing.T) {
	t.Run("Unique order indices", func(t *testing.T) {
		defaultVal := "default"
		arg1, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg1",
			DefaultString: &defaultVal,
			OrderIndex:   1,
		})
		arg2, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg2",
			DefaultString: &defaultVal,
			OrderIndex:   2,
		})
		err := validateUniqueOrderIndexArgs([]GameArg{arg1, arg2})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Duplicate order indices", func(t *testing.T) {
		defaultVal := "default"
		arg1, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg1",
			DefaultString: &defaultVal,
			OrderIndex:   1,
		})
		arg2, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg2",
			DefaultString: &defaultVal,
			OrderIndex:   1, // duplicate
		})
		err := validateUniqueOrderIndexArgs([]GameArg{arg1, arg2})
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})

	t.Run("Empty args", func(t *testing.T) {
		err := validateUniqueOrderIndexArgs([]GameArg{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestWithRequiresAdmin(t *testing.T) {
	t.Run("Set to true", func(t *testing.T) {
		exec := &GameExec{}
		requiresAdmin := true
		err := WithRequiresAdmin(&requiresAdmin)(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if exec.RequiresAdmin == nil {
			t.Fatal("expected RequiresAdmin to be non-nil")
		}
		if *exec.RequiresAdmin != true {
			t.Error("expected RequiresAdmin to be true")
		}
	})

	t.Run("Set to false", func(t *testing.T) {
		exec := &GameExec{}
		requiresAdmin := false
		err := WithRequiresAdmin(&requiresAdmin)(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if exec.RequiresAdmin == nil {
			t.Fatal("expected RequiresAdmin to be non-nil")
		}
		if *exec.RequiresAdmin != false {
			t.Error("expected RequiresAdmin to be false")
		}
	})

	t.Run("Nil value", func(t *testing.T) {
		exec := &GameExec{}
		err := WithRequiresAdmin(nil)(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if exec.RequiresAdmin != nil {
			t.Error("expected RequiresAdmin to be nil")
		}
	})
}

func TestWithFormat(t *testing.T) {
	t.Run("Valid format", func(t *testing.T) {
		exec := &GameExec{}
		format := "{{.Value}}"
		err := WithFormat(&format)(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if exec.Format == nil {
			t.Fatal("expected Format to be non-nil")
		}
		// Verify it's a valid template
		if exec.Format == nil {
			t.Error("expected Format to be a valid template")
		}
	})

	t.Run("Invalid format", func(t *testing.T) {
		exec := &GameExec{}
		invalidFormat := "{{.Invalid"
		err := WithFormat(&invalidFormat)(exec)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})

	t.Run("Nil format", func(t *testing.T) {
		exec := &GameExec{}
		err := WithFormat(nil)(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if exec.Format != nil {
			t.Error("expected Format to be nil")
		}
	})
}

func TestWithArgSeperator(t *testing.T) {
	t.Run("Set separator", func(t *testing.T) {
		exec := &GameExec{}
		separator := " "
		err := WithArgSeperator(&separator)(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if exec.ArgSeperator == nil {
			t.Fatal("expected ArgSeperator to be non-nil")
		}
		if *exec.ArgSeperator != " " {
			t.Errorf("expected ArgSeperator %q, got %q", " ", *exec.ArgSeperator)
		}
	})

	t.Run("Nil separator", func(t *testing.T) {
		exec := &GameExec{}
		err := WithArgSeperator(nil)(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if exec.ArgSeperator != nil {
			t.Error("expected ArgSeperator to be nil")
		}
	})
}

func TestWithArgs(t *testing.T) {
	t.Run("Valid args", func(t *testing.T) {
		exec := &GameExec{}
		defaultVal := "default"
		arg1, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg1",
			DefaultString: &defaultVal,
			OrderIndex:   1,
		})
		arg2, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg2",
			DefaultString: &defaultVal,
			OrderIndex:   2,
		})
		err := WithArgs(arg1, arg2)(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(exec.Args) != 2 {
			t.Errorf("expected 2 args, got %d", len(exec.Args))
		}
	})

	t.Run("Duplicate order indices", func(t *testing.T) {
		exec := &GameExec{}
		defaultVal := "default"
		arg1, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg1",
			DefaultString: &defaultVal,
			OrderIndex:   1,
		})
		arg2, _ := NewGameArg(GameArgInput{
			Role:         "string",
			Name:         "arg2",
			DefaultString: &defaultVal,
			OrderIndex:   1, // duplicate
		})
		err := WithArgs(arg1, arg2)(exec)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})

	t.Run("Empty args", func(t *testing.T) {
		exec := &GameExec{}
		err := WithArgs()(exec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(exec.Args) != 0 {
			t.Errorf("expected 0 args, got %d", len(exec.Args))
		}
	})
}

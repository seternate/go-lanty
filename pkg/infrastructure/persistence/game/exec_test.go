package game

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/pkg/domain/game"
	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

func TestGameExecRow_Assemble(t *testing.T) {
	t.Run("Valid exec row with client role", func(t *testing.T) {
		row := GameExecRow{
			ID:            uuid.New(),
			GameSlug:      "test-game",
			Role:          "client",
			Path:          "/path/to/client",
			RequiresAdmin: nil,
			Format:        nil,
			ArgSeperator:  nil,
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled == nil {
			t.Fatal("expected exec, got nil")
		}
		if assembled.Role.String() != "client" {
			t.Errorf("expected Role %q, got %q", "client", assembled.Role.String())
		}
		if assembled.Path != "/path/to/client" {
			t.Errorf("expected Path %q, got %q", "/path/to/client", assembled.Path)
		}
		if assembled.RequiresAdmin != nil {
			t.Error("expected RequiresAdmin to be nil")
		}
		if assembled.Format != nil {
			t.Error("expected Format to be nil")
		}
		if assembled.ArgSeperator != nil {
			t.Error("expected ArgSeperator to be nil")
		}
	})

	t.Run("Valid exec row with server role and options", func(t *testing.T) {
		requiresAdmin := true
		format := "{{.Path}}"
		argSeparator := " "
		row := GameExecRow{
			ID:            uuid.New(),
			GameSlug:      "test-game",
			Role:          "server",
			Path:          "/path/to/server",
			RequiresAdmin: &requiresAdmin,
			Format:        &format,
			ArgSeperator:  &argSeparator,
		}

		assembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if assembled.Role.String() != "server" {
			t.Errorf("expected Role %q, got %q", "server", assembled.Role.String())
		}
		if assembled.RequiresAdmin == nil || !*assembled.RequiresAdmin {
			t.Error("expected RequiresAdmin to be true")
		}
		if assembled.Format == nil {
			t.Error("expected Format to be set")
		}
		if assembled.ArgSeperator == nil || *assembled.ArgSeperator != " " {
			t.Error("expected ArgSeperator to be ' '")
		}
	})

	t.Run("Valid exec row with args", func(t *testing.T) {
		row := GameExecRow{
			ID:            uuid.New(),
			GameSlug:      "test-game",
			Role:          "client",
			Path:          "/path/to/client",
			RequiresAdmin: nil,
			Format:        nil,
			ArgSeperator:  nil,
		}

		argRow := GameArgRow{
			GameExecID:     row.ID,
			Role:           "flag",
			Name:           "test-arg",
			Required:       nil,
			Enabled:        nil,
			Format:         nil,
			ArgSeparator:   nil,
			Arg:            "--test",
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

		assembled, err := row.Assemble(argRow)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(assembled.Args) != 1 {
			t.Errorf("expected 1 arg, got %d", len(assembled.Args))
		}
		if assembled.Args[0].GetName() != "test-arg" {
			t.Errorf("expected arg name %q, got %q", "test-arg", assembled.Args[0].GetName())
		}
	})

	t.Run("Invalid role", func(t *testing.T) {
		row := GameExecRow{
			ID:            uuid.New(),
			GameSlug:      "test-game",
			Role:          "invalid-role",
			Path:          "/path/to/client",
			RequiresAdmin: nil,
			Format:        nil,
			ArgSeperator:  nil,
		}

		assembled, err := row.Assemble()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil exec")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})

	t.Run("Invalid path", func(t *testing.T) {
		row := GameExecRow{
			ID:            uuid.New(),
			GameSlug:      "test-game",
			Role:          "client",
			Path:          "",
			RequiresAdmin: nil,
			Format:        nil,
			ArgSeperator:  nil,
		}

		assembled, err := row.Assemble()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if assembled != nil {
			t.Error("expected nil exec")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T: %v", err, err)
		}
	})
}

func TestDisassembleGameExec(t *testing.T) {
	t.Run("Valid exec with client role", func(t *testing.T) {
		exec, err := game.RehydrateGameExec("/path/to/client", "client")
		if err != nil {
			t.Fatalf("failed to create test exec: %v", err)
		}

		row := disassembleGameExec("test-game", *exec)
		if row == nil {
			t.Fatal("expected row, got nil")
		}
		if row.GameSlug != "test-game" {
			t.Errorf("expected GameSlug %q, got %q", "test-game", row.GameSlug)
		}
		if row.Role != "client" {
			t.Errorf("expected Role %q, got %q", "client", row.Role)
		}
		if row.Path != "/path/to/client" {
			t.Errorf("expected Path %q, got %q", "/path/to/client", row.Path)
		}
		if row.ID == uuid.Nil {
			t.Error("expected non-nil UUID")
		}
	})

	t.Run("Valid exec with server role and options", func(t *testing.T) {
		requiresAdmin := true
		format := "{{.Path}}"
		argSeparator := " "
		exec, err := game.RehydrateGameExec("/path/to/server", "server",
			game.WithRequiresAdmin(&requiresAdmin),
			game.WithFormat(&format),
			game.WithArgSeperator(&argSeparator),
		)
		if err != nil {
			t.Fatalf("failed to create test exec: %v", err)
		}

		row := disassembleGameExec("test-game", *exec)
		if row.RequiresAdmin == nil || !*row.RequiresAdmin {
			t.Error("expected RequiresAdmin to be true")
		}
		if row.Format == nil || *row.Format != format {
			t.Errorf("expected Format %q, got %v", format, row.Format)
		}
		if row.ArgSeperator == nil || *row.ArgSeperator != argSeparator {
			t.Errorf("expected ArgSeperator %q, got %v", argSeparator, row.ArgSeperator)
		}
	})

	t.Run("Round trip: Disassemble then Assemble", func(t *testing.T) {
		originalExec, err := game.RehydrateGameExec("/path/to/client", "client")
		if err != nil {
			t.Fatalf("failed to create test exec: %v", err)
		}

		row := disassembleGameExec("test-game", *originalExec)
		reassembled, err := row.Assemble()
		if err != nil {
			t.Fatalf("failed to reassemble: %v", err)
		}

		if reassembled.Role != originalExec.Role {
			t.Errorf("expected Role %v, got %v", originalExec.Role, reassembled.Role)
		}
		if reassembled.Path != originalExec.Path {
			t.Errorf("expected Path %q, got %q", originalExec.Path, reassembled.Path)
		}
	})
}

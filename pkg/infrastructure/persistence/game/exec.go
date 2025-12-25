package game

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/pkg/domain/game"
)

type GameExecRow struct {
	ID            uuid.UUID `db:"id"`
	GameSlug      string    `db:"game_slug"`
	Role          string    `db:"role"`
	Path          string    `db:"path"`
	RequiresAdmin *bool     `db:"requires_admin"`
	Format        *string   `db:"format"`
	ArgSeperator  *string   `db:"arg_seperator"`
}

func (row GameExecRow) Assemble(argrows ...GameArgRow) (*game.GameExec, error) {
	args := make([]game.GameArg, 0, len(argrows))

	for _, argrow := range argrows {
		arg, err := argrow.Assemble()
		if err != nil {
			return nil, fmt.Errorf("failed to assemble game arg for game slug=%s: %w", row.GameSlug, err)
		}
		args = append(args, arg)
	}

	exec, err := game.RehydrateGameExec(row.Path, row.Role, game.WithRequiresAdmin(row.RequiresAdmin), game.WithFormat(row.Format), game.WithArgSeperator(row.ArgSeperator), game.WithArgs(args...))
	if err != nil {
		return nil, fmt.Errorf("failed to rehydrate game exec for game slug=%s (database corruption): %w", row.GameSlug, err)
	}
	return exec, nil
}

func disassembleGameExec(slug string, exec game.GameExec) *GameExecRow {
	var format *string
	if tmpl := exec.Format; tmpl != nil {
		raw := tmpl.Tree.Root.String()
		format = &raw
	}

	return &GameExecRow{
		ID:            uuid.New(),
		GameSlug:      slug,
		Role:          exec.Role.String(),
		Path:          exec.Path,
		RequiresAdmin: exec.RequiresAdmin,
		Format:        format,
		ArgSeperator:  exec.ArgSeperator,
	}
}

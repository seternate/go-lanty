package game

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/seternate/go-lanty/internal/domain/game"
)

type GameRow struct {
	Slug string `db:"slug"`
	Name string `db:"name"`
}

func (row GameRow) Assemble(execrows []GameExecRow, argrows []GameArgRow, assetrows []GameAssetRow) (*game.Game, error) {
	var err error

	argrowsByExec := make(map[uuid.UUID][]GameArgRow, len(argrows))
	for _, argrow := range argrows {
		argrowsByExec[argrow.GameExecID] = append(argrowsByExec[argrow.GameExecID], argrow)
	}

	execs := make([]game.GameExec, 0, len(execrows))
	for _, execrow := range execrows {
		exec, err := execrow.Assemble(argrowsByExec[execrow.ID]...)
		if err != nil {
			return nil, fmt.Errorf("failed to assemble executable: %w", err)
		}
		execs = append(execs, *exec)
	}

	g, err := game.RehydrateGame(row.Slug, row.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to rehydrate game (database corruption): %w", err)
	}

	for _, exec := range execs {
		g.SetExecutable(exec)
	}

	for _, assetrow := range assetrows {
		asset, err := game.RehydrateGameAsset(assetrow.AssetID, assetrow.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to rehydrate game asset for slug=%s (database corruption): %w", row.Slug, err)
		}
		g.SetAsset(*asset)
	}

	return g, nil
}

type DisassembledGame struct {
	GameRow   GameRow
	ExecRows  []GameExecRow
	ArgRows   []GameArgRow
	AssetRows []GameAssetRow
}

func Disassemble(game *game.Game) *DisassembledGame {
	disassembledGame := &DisassembledGame{}

	disassembledGame.GameRow = GameRow{
		Slug: game.Slug,
		Name: game.Name,
	}

	for _, exec := range game.Execs {
		execrow := disassembleGameExec(game.Slug, exec)
		disassembledGame.ExecRows = append(disassembledGame.ExecRows, *execrow)

		for _, arg := range exec.Args {
			disassembledGame.ArgRows = append(disassembledGame.ArgRows, *disassembleGameArg(execrow.ID, arg))
		}
	}

	for _, asset := range game.Assets {
		disassembledGame.AssetRows = append(disassembledGame.AssetRows, *disassembleGameAsset(game.Slug, asset))
	}

	return disassembledGame
}

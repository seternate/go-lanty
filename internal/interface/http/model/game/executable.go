package model

import (
	"time"

	appGameSrv "github.com/seternate/go-lanty/internal/application/game"
)

type GameExecutableResponse struct {
	Role              string                      `json:"role"`
	Path              string                      `json:"path"`
	RequiresAdmin     *bool                       `json:"requiresadmin"`
	Format            *string                     `json:"format"`
	ArgumentSeperator *string                     `json:"argumentseperator"`
	Args              []GameExecutableArgResponse `json:"args"`
	CreatedAt         time.Time                   `json:"createdat"`
}

func NewGameExecutableResponse(execView appGameSrv.GameExecView) *GameExecutableResponse {
	args := make([]GameExecutableArgResponse, 0, len(execView.Args))
	for _, arg := range execView.Args {
		args = append(args, *NewGameExecutableArgResponse(arg))
	}

	return &GameExecutableResponse{
		Role:              execView.Role,
		Path:              execView.Path,
		RequiresAdmin:     execView.RequiresAdmin,
		Format:            execView.Format,
		ArgumentSeperator: execView.ArgSeperator,
		Args:              args,
		CreatedAt:         execView.CreatedAt,
	}
}

type UpsertGameExecutableRequest struct {
	Role              string                           `json:"role"`
	Path              string                           `json:"path"`
	RequiresAdmin     *bool                            `json:"requiresadmin"`
	Format            *string                          `json:"format"`
	ArgumentSeperator *string                          `json:"argumentseperator"`
	Args              []UpsertGameExecutableArgRequest `json:"args"`
}

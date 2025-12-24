package model

import (
	"time"

	appGameSrv "github.com/seternate/go-lanty/pkg/application/game"
)

type GameResponse struct {
	Slug        string                   `json:"slug"`
	Name        string                   `json:"name"`
	Executables []GameExecutableResponse `json:"executables"`
	CreatedAt   time.Time                `json:"createdat"`
}

func NewGameResponse(gameView *appGameSrv.GameView) *GameResponse {
	executables := make([]GameExecutableResponse, 0, len(gameView.Execs))
	for _, exec := range gameView.Execs {
		executables = append(executables, *NewGameExecutableResponse(exec))
	}

	return &GameResponse{
		Slug:        gameView.Slug,
		Name:        gameView.Name,
		Executables: executables,
		CreatedAt:   gameView.CreatedAt,
	}
}

type UpsertGameRequest struct {
	Name        string                        `json:"name"`
	Executables []UpsertGameExecutableRequest `json:"executables"`
}

func (req *UpsertGameRequest) ToCommand(slug string) appGameSrv.UpsertGameCommand {
	executables := make([]appGameSrv.UpsertGameExecutable, 0, len(req.Executables))
	for _, execReq := range req.Executables {
		args := make([]appGameSrv.UpsertGameExecutableArg, 0, len(execReq.Args))
		for _, argReq := range execReq.Args {
			args = append(args, appGameSrv.UpsertGameExecutableArg{
				Role:              argReq.Role,
				Name:              argReq.Name,
				Required:          argReq.Required,
				Enabled:           argReq.Enabled,
				Format:            argReq.Format,
				ArgumentSeparator: argReq.ArgumentSeparator,
				Argument:          argReq.Argument,
				Description:       argReq.Description,
				DefaultString:     argReq.DefaultString,
				DefaultBool:       argReq.DefaultBool,
				DefaultInt:        argReq.DefaultInt,
				DefaultFloat:      argReq.DefaultFloat,
				EnumValues:        argReq.EnumValues,
				MinInt:            argReq.MinInt,
				MaxInt:            argReq.MaxInt,
				MinFloat:          argReq.MinFloat,
				MaxFloat:          argReq.MaxFloat,
				FloatPrecision:    argReq.FloatPrecision,
				OrderIndex:        argReq.OrderIndex,
			})
		}
		executables = append(executables, appGameSrv.UpsertGameExecutable{
			Role:              execReq.Role,
			Path:              execReq.Path,
			RequiresAdmin:     execReq.RequiresAdmin,
			Format:            execReq.Format,
			ArgumentSeperator: execReq.ArgumentSeperator,
			Args:              args,
		})
	}

	return appGameSrv.UpsertGameCommand{
		Slug:        slug,
		Name:        req.Name,
		Executables: executables,
	}
}

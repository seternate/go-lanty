package game

import (
	appGameSrv "github.com/seternate/go-lanty/internal/application/game"
	apimodel "github.com/seternate/go-lanty/pkg/api/models/game"
)

func NewGame(gameView *appGameSrv.GameView) *apimodel.Game {
	executables := make([]apimodel.Executable, 0, len(gameView.Execs))
	for _, exec := range gameView.Execs {
		executables = append(executables, *NewExecutable(exec))
	}

	return &apimodel.Game{
		Slug:        gameView.Slug,
		Name:        gameView.Name,
		Executables: executables,
		CreatedAt:   gameView.CreatedAt,
	}
}

func ToCommand(req *apimodel.UpsertGameRequest, slug string) appGameSrv.UpsertGameCommand {
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

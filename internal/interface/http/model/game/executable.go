package game

import (
	apimodel "github.com/seternate/go-lanty/pkg/api/models/game"
	appGameSrv "github.com/seternate/go-lanty/internal/application/game"
)

func NewExecutable(execView appGameSrv.GameExecView) *apimodel.Executable {
	args := make([]apimodel.Arg, 0, len(execView.Args))
	for _, arg := range execView.Args {
		args = append(args, *NewArg(arg))
	}

	return &apimodel.Executable{
		Role:              execView.Role,
		Path:              execView.Path,
		RequiresAdmin:     execView.RequiresAdmin,
		Format:            execView.Format,
		ArgumentSeperator: execView.ArgSeperator,
		Args:              args,
		CreatedAt:         execView.CreatedAt,
	}
}

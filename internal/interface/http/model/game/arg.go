package game

import (
	apimodel "github.com/seternate/go-lanty/pkg/api/models/game"
	appGameSrv "github.com/seternate/go-lanty/internal/application/game"
)

func NewArg(argView appGameSrv.GameArgView) *apimodel.Arg {
	return &apimodel.Arg{
		Role:              argView.Role,
		Required:          argView.Required,
		Enabled:           argView.Enabled,
		Format:            argView.Format,
		ArgumentSeparator: argView.ArgSeparator,
		Argument:          argView.Arg,
		Description:       argView.Description,
		DefaultString:     argView.DefaultString,
		DefaultBool:       argView.DefaultBool,
		DefaultInt:        argView.DefaultInt,
		DefaultFloat:      argView.DefaultFloat,
		EnumValues:        argView.Enums,
		MinInt:            argView.MinInt,
		MaxInt:            argView.MaxInt,
		MinFloat:          argView.MinFloat,
		MaxFloat:          argView.MaxFloat,
		FloatPrecision:    argView.FloatPrecision,
		OrderIndex:        argView.OrderIndex,
		CreatedAt:         argView.CreatedAt,
	}
}

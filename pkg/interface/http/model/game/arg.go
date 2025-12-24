package model

import (
	"time"

	appGameSrv "github.com/seternate/go-lanty/pkg/application/game"
)

type GameExecutableArgResponse struct {
	Role              string    `json:"role"`
	Required          *bool     `json:"required"`
	Enabled           *bool     `json:"enabled"`
	Format            *string   `json:"format"`
	ArgumentSeparator *string   `json:"argumentseperator"`
	Argument          string    `json:"argument"`
	Description       *string   `json:"description"`
	DefaultString     *string   `json:"defaultstring"`
	DefaultBool       *bool     `json:"defaultbool"`
	DefaultInt        *int64    `json:"defaultint"`
	DefaultFloat      *float64  `json:"defaultfloat"`
	EnumValues        []string  `json:"enumvalues"`
	MinInt            *int64    `json:"minint"`
	MaxInt            *int64    `json:"maxint"`
	MinFloat          *float64  `json:"minfloat"`
	MaxFloat          *float64  `json:"maxfloat"`
	FloatPrecision    *int64    `json:"floatprecision"`
	OrderIndex        int64     `json:"orderindex"`
	CreatedAt         time.Time `json:"createdat"`
}

func NewGameExecutableArgResponse(argView appGameSrv.GameArgView) *GameExecutableArgResponse {
	return &GameExecutableArgResponse{
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

type UpsertGameExecutableArgRequest struct {
	Role              string   `json:"role"`
	Name              string   `json:"name"`
	Required          *bool    `json:"required"`
	Enabled           *bool    `json:"enabled"`
	Format            *string  `json:"format"`
	ArgumentSeparator *string  `json:"argumentseperator"`
	Argument          string   `json:"argument"`
	Description       *string  `json:"description"`
	DefaultString     *string  `json:"defaultstring"`
	DefaultBool       *bool    `json:"defaultbool"`
	DefaultInt        *int64   `json:"defaultint"`
	DefaultFloat      *float64 `json:"defaultfloat"`
	EnumValues        []string `json:"enumvalues"`
	MinInt            *int64   `json:"minint"`
	MaxInt            *int64   `json:"maxint"`
	MinFloat          *float64 `json:"minfloat"`
	MaxFloat          *float64 `json:"maxfloat"`
	FloatPrecision    *int64   `json:"floatprecision"`
	OrderIndex        int64    `json:"orderindex"`
}

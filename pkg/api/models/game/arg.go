package gamemodel

import (
	"time"
)

type Arg struct {
	Role              string    `json:"role"`
	Required          *bool     `json:"required,omitempty"`
	Enabled           *bool     `json:"enabled,omitempty"`
	Format            *string   `json:"format,omitempty"`
	ArgumentSeparator *string   `json:"argumentseperator,omitempty"`
	Argument          string    `json:"argument"`
	Description       *string   `json:"description,omitempty"`
	DefaultString     *string   `json:"defaultstring,omitempty"`
	DefaultBool       *bool     `json:"defaultbool,omitempty"`
	DefaultInt        *int64    `json:"defaultint,omitempty"`
	DefaultFloat      *float64  `json:"defaultfloat,omitempty"`
	EnumValues        []string  `json:"enumvalues"`
	MinInt            *int64    `json:"minint,omitempty"`
	MaxInt            *int64    `json:"maxint,omitempty"`
	MinFloat          *float64  `json:"minfloat,omitempty"`
	MaxFloat          *float64  `json:"maxfloat,omitempty"`
	FloatPrecision    *int64    `json:"floatprecision,omitempty"`
	OrderIndex        int64     `json:"orderindex"`
	CreatedAt         time.Time `json:"createdat"`
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

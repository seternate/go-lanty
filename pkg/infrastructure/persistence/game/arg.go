package game

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/seternate/go-lanty/pkg/domain/game"
)

type GameArgRow struct {
	GameExecID     uuid.UUID      `db:"game_exec_id"`
	Role           string         `db:"role"`
	Name           string         `db:"name"`
	Required       *bool          `db:"required"`
	Enabled        *bool          `db:"enabled"`
	Format         *string        `db:"format"`
	ArgSeparator   *string        `db:"arg_separator"`
	Arg            string         `db:"arg"`
	Description    *string        `db:"description"`
	DefaultString  *string        `db:"default_string"`
	DefaultBool    *bool          `db:"default_bool"`
	DefaultInt     *int64         `db:"default_int"`
	DefaultFloat   *float64       `db:"default_float"`
	Enums          pq.StringArray `db:"enums"`
	MinInt         *int64         `db:"min_int"`
	MaxInt         *int64         `db:"max_int"`
	MinFloat       *float64       `db:"min_float"`
	MaxFloat       *float64       `db:"max_float"`
	FloatPrecision *int64         `db:"float_precision"`
	OrderIndex     int64          `db:"order_index"`
}

func (row GameArgRow) Assemble() (game.GameArg, error) {
	arg, err := game.RehydrateGameArg(game.GameArgInput{
		Role:           row.Role,
		Name:           row.Name,
		Required:       row.Required,
		Enabled:        row.Enabled,
		Format:         row.Format,
		ArgSeparator:   row.ArgSeparator,
		Arg:            row.Arg,
		Description:    row.Description,
		OrderIndex:     row.OrderIndex,
		DefaultString:  row.DefaultString,
		DefaultBool:    row.DefaultBool,
		DefaultInt:     row.DefaultInt,
		DefaultFloat:   row.DefaultFloat,
		Enums:          row.Enums,
		MinInt:         row.MinInt,
		MaxInt:         row.MaxInt,
		MinFloat:       row.MinFloat,
		MaxFloat:       row.MaxFloat,
		FloatPrecision: row.FloatPrecision,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to rehydrate game arg for game exec id=%s: %w", row.GameExecID, err)
	}
	return arg, nil
}

func disassembleGameArg(execid uuid.UUID, arg game.GameArg) *GameArgRow {
	var format *string
	if tmpl := arg.GetFormat(); tmpl != nil {
		raw := tmpl.Tree.Root.String()
		format = &raw
	}

	var defaultString *string
	var defaultBool *bool
	var defaultInt *int64
	var defaultFloat *float64
	var enums []string
	var minInt *int64
	var maxInt *int64
	var minFloat *float64
	var maxFloat *float64
	var floatPrecision *int64

	switch v := arg.(type) {
	case *game.GameArgString:
		defaultString = &v.Default
	case *game.GameArgBool:
		defaultBool = &v.Default
	case *game.GameArgEnum:
		defaultString = &v.Default
		enums = v.Values
	case *game.GameArgInt:
		defaultInt = &v.Default
		minInt = &v.Min
		maxInt = &v.Max
	case *game.GameArgFloat:
		defaultFloat = &v.Default
		minFloat = &v.Min
		maxFloat = &v.Max
		floatPrecision = &v.FloatPrecision
	}

	return &GameArgRow{
		GameExecID:     execid,
		Role:           arg.Role().String(),
		Name:           arg.GetName(),
		Required:       arg.GetRequired(),
		Enabled:        arg.GetEnabled(),
		Format:         format,
		ArgSeparator:   arg.GetArgSeparator(),
		Arg:            arg.GetArg(),
		Description:    arg.GetDescription(),
		DefaultString:  defaultString,
		DefaultBool:    defaultBool,
		DefaultInt:     defaultInt,
		DefaultFloat:   defaultFloat,
		Enums:          enums,
		MinInt:         minInt,
		MaxInt:         maxInt,
		MinFloat:       minFloat,
		MaxFloat:       maxFloat,
		FloatPrecision: floatPrecision,
		OrderIndex:     arg.GetOrderIndex(),
	}
}

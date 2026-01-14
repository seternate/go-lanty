package game

import (
	"fmt"
	"slices"
	"strings"
	"text/template"

	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

type GameArgRole string

const (
	GAME_ARG_ROLE_UNDEFINED GameArgRole = ""
	GAME_ARG_ROLE_FLAG      GameArgRole = "flag"
	GAME_ARG_ROLE_STRING    GameArgRole = "string"
	GAME_ARG_ROLE_BOOL      GameArgRole = "bool"
	GAME_ARG_ROLE_ENUM      GameArgRole = "enum"
	GAME_ARG_ROLE_INT       GameArgRole = "int"
	GAME_ARG_ROLE_FLOAT     GameArgRole = "float"
)

var supportedGameArgRoles = map[string]GameArgRole{
	"flag":   GAME_ARG_ROLE_FLAG,
	"string": GAME_ARG_ROLE_STRING,
	"bool":   GAME_ARG_ROLE_BOOL,
	"enum":   GAME_ARG_ROLE_ENUM,
	"int":    GAME_ARG_ROLE_INT,
	"float":  GAME_ARG_ROLE_FLOAT,
}

func ParseGameArgRole(role string) (GameArgRole, error) {
	if gameRole, ok := supportedGameArgRoles[role]; ok {
		return gameRole, nil
	}
	available := make([]string, 0, len(supportedGameArgRoles))
	for _, role := range supportedGameArgRoles {
		available = append(available, role.String())
	}

	return GAME_ARG_ROLE_UNDEFINED, domainerr.ValidationErr("game arg role", "undefined").WithExpected(strings.Join(available, ", ")).WithGot(role)
}

func (role GameArgRole) String() string {
	return string(role)
}

type GameArgInput struct {
	Role string

	Name         string
	Required     *bool
	Enabled      *bool
	Format       *string
	ArgSeparator *string
	Arg          string
	Description  *string
	OrderIndex   int64

	DefaultString *string
	DefaultBool   *bool
	DefaultInt    *int64
	DefaultFloat  *float64

	Enums          []string
	MinInt         *int64
	MaxInt         *int64
	MinFloat       *float64
	MaxFloat       *float64
	FloatPrecision *int64
}

type GameArg interface {
	Role() GameArgRole
	GetName() string
	GetRequired() *bool
	GetEnabled() *bool
	GetFormat() *template.Template
	GetArgSeparator() *string
	GetArg() string
	GetDescription() *string
	GetOrderIndex() int64
}

func NewGameArg(input GameArgInput) (GameArg, error) {
	gameArg, err := hydrateGameArg(input)
	if err != nil {
		return nil, domainerr.InvariantViolationErr("game arg", input.Name).WithCause(err)
	}

	return gameArg, nil
}

func RehydrateGameArg(input GameArgInput) (GameArg, error) {
	gameArg, err := hydrateGameArg(input)
	if err != nil {
		return nil, domainerr.TrustedInvariantViolationErr("game arg", input.Name).WithCause(err)
	}

	return gameArg, nil
}

func hydrateGameArg(input GameArgInput) (GameArg, error) {
	validationErrors := domainerr.ValidationErrs()

	role, err := ParseGameArgRole(input.Role)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, err
	}

	format, err := validateArgFormat(input.Format)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, err
	}

	err = validationErrors.Wrap(validateOrderIndex(input.OrderIndex))
	if err != nil {
		return nil, err
	}

	base := GameArgFlag{
		Name:         input.Name,
		Required:     input.Required,
		Enabled:      input.Enabled,
		Format:       format,
		ArgSeparator: input.ArgSeparator,
		Arg:          input.Arg,
		Description:  input.Description,
		OrderIndex:   input.OrderIndex,
	}

	var arg GameArg
	switch role {
	case GAME_ARG_ROLE_FLAG:
		arg = &base
	case GAME_ARG_ROLE_STRING:
		err = validationErrors.Wrap(validateStringArg(input))
		if err != nil {
			return nil, err
		}
		if validationErrors.HasErrors() {
			return nil, validationErrors
		}
		arg = &GameArgString{
			GameArgFlag: base,
			Default:     *input.DefaultString,
		}
	case GAME_ARG_ROLE_BOOL:
		err = validationErrors.Wrap(validateBoolArg(input))
		if err != nil {
			return nil, err
		}
		if validationErrors.HasErrors() {
			return nil, validationErrors
		}
		arg = &GameArgBool{
			GameArgFlag: base,
			Default:     *input.DefaultBool,
		}
	case GAME_ARG_ROLE_ENUM:
		err = validationErrors.Wrap(validateEnumArg(input))
		if err != nil {
			return nil, err
		}
		if validationErrors.HasErrors() {
			return nil, validationErrors
		}
		arg = &GameArgEnum{
			GameArgFlag: base,
			Default:     *input.DefaultString,
			Values:      input.Enums,
		}
	case GAME_ARG_ROLE_INT:
		err = validationErrors.Wrap(validateIntArg(input))
		if err != nil {
			return nil, err
		}
		if validationErrors.HasErrors() {
			return nil, validationErrors
		}
		arg = &GameArgInt{
			GameArgFlag: base,
			Default:     *input.DefaultInt,
			Min:         *input.MinInt,
			Max:         *input.MaxInt,
		}
	case GAME_ARG_ROLE_FLOAT:
		err = validationErrors.Wrap(validateFloatArg(input))
		if err != nil {
			return nil, err
		}
		if validationErrors.HasErrors() {
			return nil, validationErrors
		}
		arg = &GameArgFloat{
			GameArgFlag:    base,
			Default:        *input.DefaultFloat,
			Min:            *input.MinFloat,
			Max:            *input.MaxFloat,
			FloatPrecision: *input.FloatPrecision,
		}
	default:
		return nil, domainerr.InternalError{Reason: fmt.Sprintf("undefined game argument role=%s for name=%s", input.Role, input.Name)}
	}

	if validationErrors.HasErrors() {
		return nil, validationErrors
	}

	return arg, nil
}

var _ GameArg = (*GameArgFlag)(nil)

type GameArgFlag struct {
	Name         string
	Required     *bool
	Enabled      *bool
	Format       *template.Template
	ArgSeparator *string
	Arg          string
	Description  *string
	OrderIndex   int64
}

func (arg GameArgFlag) Role() GameArgRole             { return GAME_ARG_ROLE_FLAG }
func (arg GameArgFlag) GetName() string               { return arg.Name }
func (arg GameArgFlag) GetRequired() *bool            { return arg.Required }
func (arg GameArgFlag) GetEnabled() *bool             { return arg.Enabled }
func (arg GameArgFlag) GetFormat() *template.Template { return arg.Format }
func (arg GameArgFlag) GetArgSeparator() *string      { return arg.ArgSeparator }
func (arg GameArgFlag) GetArg() string                { return arg.Arg }
func (arg GameArgFlag) GetDescription() *string       { return arg.Description }
func (arg GameArgFlag) GetOrderIndex() int64          { return arg.OrderIndex }

var _ GameArg = (*GameArgString)(nil)

type GameArgString struct {
	GameArgFlag
	Default string
}

func (arg GameArgString) Role() GameArgRole  { return GAME_ARG_ROLE_STRING }
func (arg GameArgString) GetDefault() string { return arg.Default }

var _ GameArg = (*GameArgBool)(nil)

type GameArgBool struct {
	GameArgFlag
	Default bool
}

func (a GameArgBool) Role() GameArgRole { return GAME_ARG_ROLE_BOOL }
func (a GameArgBool) GetDefault() bool  { return a.Default }

var _ GameArg = (*GameArgEnum)(nil)

type GameArgEnum struct {
	GameArgFlag
	Default string
	Values  []string
}

func (arg GameArgEnum) Role() GameArgRole  { return GAME_ARG_ROLE_ENUM }
func (arg GameArgEnum) GetDefault() string { return arg.Default }
func (arg GameArgEnum) GetEnums() []string { return arg.Values }

var _ GameArg = (*GameArgInt)(nil)

type GameArgInt struct {
	GameArgFlag
	Default int64
	Min     int64
	Max     int64
}

func (arg GameArgInt) Role() GameArgRole { return GAME_ARG_ROLE_INT }
func (arg GameArgInt) GetDefault() int64 { return arg.Default }
func (arg GameArgInt) GetMin() int64     { return arg.Min }
func (arg GameArgInt) GetMax() int64     { return arg.Max }

var _ GameArg = (*GameArgFloat)(nil)

type GameArgFloat struct {
	GameArgFlag
	Default        float64
	Min            float64
	Max            float64
	FloatPrecision int64
}

func (arg GameArgFloat) Role() GameArgRole        { return GAME_ARG_ROLE_FLOAT }
func (arg GameArgFloat) GetDefault() float64      { return arg.Default }
func (arg GameArgFloat) GetMinValue() float64     { return arg.Min }
func (arg GameArgFloat) GetMaxValue() float64     { return arg.Max }
func (arg GameArgFloat) GetFloatPrecision() int64 { return arg.FloatPrecision }

func validateArgFormat(format *string) (*template.Template, error) {
	if format == nil {
		return nil, nil
	}

	tmpl, err := template.New("").Parse(*format)
	if err != nil {
		return nil, domainerr.ValidationErr("format", "failed to parse").WithGot(*format).WithCause(err)
	}

	return tmpl, nil
}

func validateOrderIndex(orderIndex int64) error {
	if orderIndex < 0 {
		return domainerr.ValidationErr("order index", "must be non-negative").WithGot("order_index=%d", orderIndex)
	}

	return nil
}

func validateStringArg(input GameArgInput) error {
	if input.DefaultString == nil || len(*input.DefaultString) == 0 {
		return domainerr.ValidationErr("default", "can not be empty")
	}

	return nil
}

func validateBoolArg(input GameArgInput) error {
	if input.DefaultBool == nil {
		return domainerr.ValidationErr("default", "can not be empty")
	}

	return nil
}

func validateEnumArg(input GameArgInput) error {
	validationErrors := domainerr.ValidationErrs()

	if input.DefaultString == nil || len(*input.DefaultString) == 0 {
		validationErrors.Wrap(domainerr.ValidationErr("default", "can not be empty"))
	}

	if len(input.Enums) == 0 {
		validationErrors.Wrap(domainerr.ValidationErr("enum values", "must be provided"))
	}

	if input.DefaultString != nil && len(*input.DefaultString) > 0 && !slices.Contains(input.Enums, *input.DefaultString) {
		validationErrors.Wrap(domainerr.ValidationErr("default", "must be part of values").WithExpected(strings.Join(input.Enums, ", ")).WithGot(*input.DefaultString))
	}

	return validationErrors.OrNil()
}

func validateIntArg(input GameArgInput) error {
	validationErrors := domainerr.ValidationErrs()

	if input.DefaultInt == nil {
		validationErrors.Wrap(domainerr.ValidationErr("default", "can not be empty"))
	}

	if input.MinInt == nil {
		validationErrors.Wrap(domainerr.ValidationErr("min value", "can not be empty"))
	}
	if input.MaxInt == nil {
		validationErrors.Wrap(domainerr.ValidationErr("max value", "can not be empty"))
	}

	if input.MinInt != nil && input.MaxInt != nil && *input.MinInt > *input.MaxInt {
		validationErrors.Wrap(domainerr.ValidationErr("min value", "cannot be greater than max value").WithExpected("less than or equal to %d", *input.MaxInt).WithGot("%d", *input.MinInt))
	}

	if input.DefaultInt != nil && input.MinInt != nil && *input.DefaultInt < *input.MinInt {
		validationErrors.Wrap(domainerr.ValidationErr("default", "can not be less than min value").WithExpected("more than or equal to %d", *input.MinInt).WithGot("%d", *input.DefaultInt))
	}

	if input.DefaultInt != nil && input.MaxInt != nil && *input.DefaultInt > *input.MaxInt {
		validationErrors.Wrap(domainerr.ValidationErr("default", "can not be greater than max value").WithExpected("less than or equal to %d", *input.MaxInt).WithGot("%d", *input.DefaultInt))
	}

	return validationErrors.OrNil()
}

func validateFloatArg(input GameArgInput) error {
	validationErrors := domainerr.ValidationErrs()

	if input.DefaultFloat == nil {
		validationErrors.Wrap(domainerr.ValidationErr("default", "can not be empty"))
	}

	if input.MinFloat == nil {
		validationErrors.Wrap(domainerr.ValidationErr("min value", "can not be empty"))
	}
	if input.MaxFloat == nil {
		validationErrors.Wrap(domainerr.ValidationErr("max value", "can not be empty"))
	}

	if input.MinFloat != nil && input.MaxFloat != nil && *input.MinFloat > *input.MaxFloat {
		validationErrors.Wrap(domainerr.ValidationErr("min value", "cannot be greater than max value").WithExpected("less than or equal to %f", *input.MaxFloat).WithGot("%f", *input.MinFloat))
	}

	if input.FloatPrecision == nil {
		validationErrors.Wrap(domainerr.ValidationErr("float precision", "can not be empty"))
	}

	if input.DefaultFloat != nil && input.MinFloat != nil && *input.DefaultFloat < *input.MinFloat {
		validationErrors.Wrap(domainerr.ValidationErr("default", "can not be less than min value").WithExpected("more than or equal to %f", *input.MinFloat).WithGot("%f", *input.DefaultFloat))
	}

	if input.DefaultFloat != nil && input.MaxFloat != nil && *input.DefaultFloat > *input.MaxFloat {
		validationErrors.Wrap(domainerr.ValidationErr("default", "can not be greater than max value").WithExpected("less than or equal to %f", *input.MaxFloat).WithGot("%f", *input.DefaultFloat))
	}

	return validationErrors.OrNil()
}

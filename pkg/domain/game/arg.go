package game

import (
	"fmt"
	"slices"
	"strings"
	"text/template"

	domainErrors "github.com/seternate/go-lanty/pkg/domain/error"
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

	return GAME_ARG_ROLE_UNDEFINED, domainErrors.ValidationErr("game arg role", "undefined").WithExpected(strings.Join(available, ", ")).WithGot(role)
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
	return RehydrateGameArg(input)
}

func RehydrateGameArg(input GameArgInput) (GameArg, error) {
	validationErrors := domainErrors.ValidationErrs().WithMessage("failed to validate arg name=%s", input.Name)

	role, err := ParseGameArgRole(input.Role)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, fmt.Errorf("failed to parse for arg name=%q: %w", input.Name, err)
	}

	format, err := validateArgFormat(input.Format)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, fmt.Errorf("failed to validate for arg name=%q: %w", input.Name, err)
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
		err = validationErrors.Wrap(validateStringArg(input), "failed for string argument")
		if err != nil {
			return nil, fmt.Errorf("failed to validate string argument for arg name=%q: %w", input.Name, err)
		}
		arg = &GameArgString{
			GameArgFlag: base,
			Default:     *input.DefaultString,
		}
	case GAME_ARG_ROLE_BOOL:
		err = validationErrors.Wrap(validateBoolArg(input), "failed for boolean argument")
		if err != nil {
			return nil, fmt.Errorf("failed to validate boolean argument for arg name=%q: %w", input.Name, err)
		}
		arg = &GameArgBool{
			GameArgFlag: base,
			Default:     *input.DefaultBool,
		}
	case GAME_ARG_ROLE_ENUM:
		err = validationErrors.Wrap(validateEnumArg(input), "failed for enum argument")
		if err != nil {
			return nil, fmt.Errorf("failed to validate enum argument for arg name=%q: %w", input.Name, err)
		}
		arg = &GameArgEnum{
			GameArgFlag: base,
			Default:     *input.DefaultString,
			Values:      input.Enums,
		}
	case GAME_ARG_ROLE_INT:
		err = validationErrors.Wrap(validateIntArg(input), "failed for integer argument")
		if err != nil {
			return nil, fmt.Errorf("failed to validate integer argument for arg name=%q: %w", input.Name, err)
		}
		arg = &GameArgInt{
			GameArgFlag: base,
			Default:     *input.DefaultInt,
			Min:         *input.MinInt,
			Max:         *input.MaxInt,
		}
	case GAME_ARG_ROLE_FLOAT:
		err = validationErrors.Wrap(validateFloatArg(input), "failed for float argument")
		if err != nil {
			return nil, fmt.Errorf("failed to validate float argument for arg name=%q: %w", input.Name, err)
		}
		arg = &GameArgFloat{
			GameArgFlag:    base,
			Default:        *input.DefaultFloat,
			Min:            *input.MinFloat,
			Max:            *input.MaxFloat,
			FloatPrecision: *input.FloatPrecision,
		}
	default:
		return nil, domainErrors.InternalError{Reason: fmt.Sprintf("undefined game argument role=%s for name=%s", input.Role, input.Name)}
	}

	if len(validationErrors.Errors) > 0 {
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
		return nil, domainErrors.ValidationErr("format", "failed to parse").WithGot(*format).WithCause(err)
	}

	return tmpl, nil
}

func validateStringArg(input GameArgInput) error {
	if input.DefaultString == nil || len(*input.DefaultString) == 0 {
		return domainErrors.ValidationErr("default", "can not be empty")
	}

	return nil
}

func validateBoolArg(input GameArgInput) error {
	if input.DefaultBool == nil {
		return domainErrors.ValidationErr("default", "can not be empty")
	}

	return nil
}

func validateEnumArg(input GameArgInput) error {
	validationErrors := domainErrors.ValidationErrs()

	if input.DefaultString == nil || len(*input.DefaultString) == 0 {
		validationErrors.Wrap(domainErrors.ValidationErr("default", "can not be empty"))
	}

	if len(input.Enums) == 0 {
		validationErrors.Wrap(domainErrors.ValidationErr("enum values", "must be provided"))
	}

	if !slices.Contains(input.Enums, *input.DefaultString) {
		validationErrors.Wrap(domainErrors.ValidationErr("default", "must be part of values").WithExpected(strings.Join(input.Enums, ", ")).WithGot(*input.DefaultString))

	}

	return validationErrors
}

func validateIntArg(input GameArgInput) error {
	validationErrors := domainErrors.ValidationErrs()

	if input.DefaultInt == nil {
		validationErrors.Wrap(domainErrors.ValidationErr("default", "can not be empty"))
	}

	if input.MinInt == nil {
		validationErrors.Wrap(domainErrors.ValidationErr("min value", "can not be empty"))
	}
	if input.MaxInt == nil {
		validationErrors.Wrap(domainErrors.ValidationErr("max value", "can not be empty"))
	}

	if input.MinInt != nil && input.MaxInt != nil && *input.MinInt > *input.MaxInt {
		validationErrors.Wrap(domainErrors.ValidationErr("min value", "cannot be greater than max value").WithExpected("less than or equal to %d", *input.MaxInt).WithGot("%d", *input.MinInt))
	}

	if *input.DefaultInt < *input.MinInt {
		validationErrors.Wrap(domainErrors.ValidationErr("default", "can not be less than min value").WithExpected("more than or equal to %d", *input.MinInt).WithGot("%d", *input.DefaultInt))
	}

	if *input.DefaultInt > *input.MaxInt {
		validationErrors.Wrap(domainErrors.ValidationErr("default", "can not be greater than max value").WithExpected("less than or equal to %d", *input.MaxInt).WithGot("%d", *input.DefaultInt))
	}

	return validationErrors
}

func validateFloatArg(input GameArgInput) error {
	validationErrors := domainErrors.ValidationErrs()

	if input.DefaultFloat == nil {
		validationErrors.Wrap(domainErrors.ValidationErr("default", "can not be empty"))
	}

	if input.MinFloat == nil {
		validationErrors.Wrap(domainErrors.ValidationErr("min value", "can not be empty"))
	}
	if input.MaxFloat == nil {
		validationErrors.Wrap(domainErrors.ValidationErr("max value", "can not be empty"))
	}

	if input.MinFloat != nil && input.MaxFloat != nil && *input.MinFloat > *input.MaxFloat {
		validationErrors.Wrap(domainErrors.ValidationErr("min value", "cannot be greater than max value").WithExpected("less than or equal to %f", *input.MaxFloat).WithGot("%f", *input.MinFloat))
	}

	if input.FloatPrecision == nil {
		validationErrors.Wrap(domainErrors.ValidationErr("float precision", "can not be empty"))
	}

	if *input.DefaultFloat < *input.MinFloat {
		validationErrors.Wrap(domainErrors.ValidationErr("default", "can not be less than min value").WithExpected("more than or equal to %f", *input.MinFloat).WithGot("%f", *input.DefaultFloat))
	}

	if *input.DefaultFloat > *input.MaxFloat {
		validationErrors.Wrap(domainErrors.ValidationErr("default", "can not be greater than max value").WithExpected("less than or equal to %f", *input.MaxFloat).WithGot("%f", *input.DefaultFloat))
	}

	return validationErrors
}

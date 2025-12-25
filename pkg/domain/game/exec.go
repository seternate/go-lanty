package game

import (
	"strings"
	"text/template"

	domainerr "github.com/seternate/go-lanty/pkg/domain/error"
)

type GameExecRole string

const (
	GAME_EXEC_ROLE_UNDEFINED GameExecRole = ""
	GAME_EXEC_ROLE_CLIENT    GameExecRole = "client"
	GAME_EXEC_ROLE_SERVER    GameExecRole = "server"
)

var supportedGameExecRoles = map[string]GameExecRole{
	"client": GAME_EXEC_ROLE_CLIENT,
	"server": GAME_EXEC_ROLE_SERVER,
}

func ParseGameExecRole(role string) (GameExecRole, error) {
	if gameRole, ok := supportedGameExecRoles[role]; ok {
		return gameRole, nil
	}
	available := make([]string, 0, len(supportedGameExecRoles))
	for _, role := range supportedGameExecRoles {
		available = append(available, role.String())
	}

	return GAME_EXEC_ROLE_UNDEFINED, domainerr.ValidationErr("game exec role", "undefined").WithExpected(strings.Join(available, ", ")).WithGot(role)
}

func (role GameExecRole) String() string {
	return string(role)
}

type GameExec struct {
	Role          GameExecRole
	Path          string
	RequiresAdmin *bool
	Format        *template.Template
	ArgSeperator  *string
	Args          []GameArg
}

func NewGameExec(path string, role string, opts ...GameExecOpts) (*GameExec, error) {
	gameExec, err := hydrateGameExec(path, role, opts...)
	if err != nil {
		return nil, domainerr.InvariantViolationErr("game exec", path).WithCause(err)
	}

	return gameExec, nil
}

func RehydrateGameExec(path string, role string, opts ...GameExecOpts) (*GameExec, error) {
	gameExec, err := hydrateGameExec(path, role, opts...)
	if err != nil {
		return nil, domainerr.TrustedInvariantViolationErr("game exec", path).WithCause(err)
	}

	return gameExec, nil
}

func hydrateGameExec(path string, role string, opts ...GameExecOpts) (*GameExec, error) {
	validationErrors := domainerr.ValidationErrs()

	parsedRole, err := ParseGameExecRole(role)
	err = validationErrors.Wrap(err)
	if err != nil {
		return nil, err
	}

	err = validationErrors.Wrap(validatePath(path))
	if err != nil {
		return nil, err
	}

	exec := &GameExec{
		Role: parsedRole,
		Path: path,
	}
	for _, opt := range opts {
		err = validationErrors.Wrap(opt(exec))
		if err != nil {
			return nil, err
		}
	}

	if validationErrors.HasErrors() {
		return nil, validationErrors
	}

	return exec, nil
}

func validatePath(path string) error {
	if len(path) == 0 {
		return domainerr.ValidationErr("path", "can not be empty")
	}

	return nil
}

func validateUniqueOrderIndexArgs(args []GameArg) error {
	seen := make(map[int64]string)
	for _, arg := range args {
		if duplicate, exists := seen[arg.GetOrderIndex()]; exists {
			return domainerr.ValidationErr("arg order index", "must be unique").WithGot("order_index=%d, duplicate_arg_names=%s, %s", arg.GetOrderIndex(), arg.GetName(), duplicate)
		}
		seen[arg.GetOrderIndex()] = arg.GetName()
	}

	return nil
}

type GameExecOpts func(*GameExec) error

func WithRequiresAdmin(requiresAdmin *bool) GameExecOpts {
	return func(exec *GameExec) error {
		if requiresAdmin != nil {
			val := *requiresAdmin
			exec.RequiresAdmin = &val
		}
		return nil
	}
}

func WithFormat(format *string) GameExecOpts {
	return func(exec *GameExec) error {
		if format != nil {
			val := *format
			tmpl, err := template.New("").Parse(val)
			if err != nil {
				return domainerr.ValidationErr("format template", "failed to parse").WithGot(val).WithCause(err)
			}
			exec.Format = tmpl
		}
		return nil
	}
}

func WithArgSeperator(seperator *string) GameExecOpts {
	return func(exec *GameExec) error {
		if seperator != nil {
			val := *seperator
			exec.ArgSeperator = &val
		}
		return nil
	}
}

func WithArgs(args ...GameArg) GameExecOpts {
	return func(exec *GameExec) error {
		if len(args) > 0 {
			err := validateUniqueOrderIndexArgs(args)
			if err != nil {
				return err
			}
			exec.Args = args
		}
		return nil
	}
}

package game

import (
	"errors"
	"fmt"

	"github.com/seternate/go-lanty/pkg/game/argument"
)

type Game struct {
	Slug   string      `json:"slug" yaml:"slug"`
	Name   string      `json:"name" yaml:"name"`
	Client Executable  `json:"client" yaml:"client"`
	Server *Executable `json:"server,omitempty" yaml:"server,omitempty"`
}

func (game *Game) ValidateLazy() (err error) {
	if len(game.Slug) == 0 {
		err = errors.Join(err, errors.New("missing slug"))
	}
	if len(game.Name) == 0 {
		err = errors.Join(err, errors.New("missing name"))
	}
	validateClientErr := game.Client.ValidateLazy()
	if validateClientErr != nil {
		validateClientErr = fmt.Errorf("%s: %w", "Client", validateClientErr)
	}
	err = errors.Join(err, validateClientErr)
	if game.Server != nil {
		validateServerErr := game.Server.ValidateLazy()
		if validateServerErr != nil {
			validateServerErr = fmt.Errorf("%s: %w", "Server", validateServerErr)
		}
		err = errors.Join(err, validateServerErr)
	}
	return
}

func (game *Game) NormalizeState() {
	game.Client.NormalizeState()
	if game.Server != nil {
		game.Server.NormalizeState()
	}
}

func (game *Game) CanConnectToServer() bool {
	return game.Client.CanConnect()
}

func (game *Game) CanStartServer() bool {
	return game.Server != nil
}

func (left Game) Equal(right Game) bool {
	equal := left.Slug == right.Slug && left.Name == right.Name && left.Client.Equal(right.Client)
	if left.Server == nil && right.Server == nil {
		return equal
	} else if left.Server != nil && right.Server != nil {
		return equal && left.Server.Equal(*right.Server)
	}
	return false
}

type Executable struct {
	Executable string              `json:"executable" yaml:"executable"`
	Arguments  *argument.Arguments `json:"arguments,omitempty" yaml:"arguments,omitempty"`
}

func (executable *Executable) ValidateLazy() (err error) {
	if len(executable.Executable) == 0 {
		err = errors.Join(err, errors.New("missing executable"))
	}
	if executable.Arguments != nil {
		err = errors.Join(err, executable.Arguments.ValidateLazy())
	}
	return
}

func (executable *Executable) NormalizeState() {
	if executable.Arguments != nil {
		executable.Arguments.NormalizeState()
	}
}

func (executable *Executable) CanConnect() bool {
	if executable.Arguments == nil {
		return false
	}
	for _, arg := range executable.Arguments.Arguments {
		if arg.GetType() == argument.TYPE_CONNECT {
			return true
		}
	}
	return false
}

func (executable *Executable) Args() (args []string, err error) {
	if executable.Arguments != nil {
		args, err = executable.Arguments.Parse()
	}
	return
}

func (executable *Executable) ParseConnectArg(ip string) (connectArg []string, err error) {
	for _, arg := range executable.Arguments.Arguments {
		if arg.GetType() == argument.TYPE_CONNECT {
			connectArg, err = arg.(*argument.Connect).ParseWithIP(ip)
			return
		}
	}
	err = errors.New("executable can not connect")
	return
}

func (left *Executable) Equal(right Executable) bool {
	equal := left.Executable == right.Executable
	if left.Arguments == nil && right.Arguments == nil {
		return equal
	} else if left.Arguments != nil && right.Arguments != nil {
		return equal && left.Arguments.Equal(*right.Arguments)
	}
	return false
}

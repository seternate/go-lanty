package handler

import (
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/setting"
	"github.com/seternate/go-lanty/pkg/user"
)

type Handler struct {
	Setting     *setting.Settings
	Gamehandler *Gamehandler
	Userhandler *Userhandler
}

func NewHandler(settings *setting.Settings) *Handler {
	return &Handler{
		Setting:     settings,
		Gamehandler: nil,
		Userhandler: nil,
	}
}

func (handler *Handler) WithGamehandler(games game.Games) *Handler {
	gamehandler := &Gamehandler{
		parent: handler,
		Games:  games,
	}
	handler.Gamehandler = gamehandler

	return handler
}

func (handler *Handler) WithUserhandler() *Handler {
	userhandler := &Userhandler{
		parent: handler,
		Users:  make(map[string]user.User),
	}
	handler.Userhandler = userhandler

	return handler
}

package handler

import (
	"github.com/rs/zerolog/log"
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

func (handler *Handler) WithGamehandler() *Handler {
	games, err := game.LoadFromDirectory(handler.Setting.GameConfigDirectory)
	if err != nil {
		log.Fatal().Err(err).Str("directory", handler.Setting.GameConfigDirectory).Msg("failed to load game configuration files")
	}
	log.Debug().Int("size", games.Size()).Msg("successfully loaded games from configuration files")

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

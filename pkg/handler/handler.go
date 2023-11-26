package handler

import (
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-server/pkg/game"
	"github.com/seternate/go-lanty-server/pkg/setting"
	"github.com/seternate/go-lanty-server/pkg/user"
)

type Handler struct {
	Setting     *setting.Settings
	GameHandler *GameHandler
	UserHandler *UserHandler
}

func NewHandler(settings *setting.Settings) *Handler {
	return &Handler{Setting: settings, GameHandler: nil}
}

func (h *Handler) WithGameHandler(games game.Games) *Handler {
	gameHandler := &GameHandler{
		handler: h,
		Games:   make(map[string]game.Game),
	}
	h.GameHandler = gameHandler

	for _, game := range games {
		gameHandler.Games[game.Slug] = game
		log.Debug().Msgf("Added '%s' to GamesHandler", game.Slug)
	}
	return h
}

func (h *Handler) WithUserHandler() *Handler {
	userHandler := &UserHandler{
		handler: h,
		Users:   make(map[string]user.User),
	}
	h.UserHandler = userHandler

	return h
}

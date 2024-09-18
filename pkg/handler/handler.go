package handler

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/setting"
	"golang.org/x/sync/errgroup"
)

type Handler struct {
	Setting         *setting.Settings
	Healthhandler   *HealthHandler
	Gamehandler     *Gamehandler
	Userhandler     *Userhandler
	Downloadhandler *Downloadhandler
	Chathandler     *ChatHandler
	Filehandler     *Filehandler
}

func NewHandler(settings *setting.Settings) *Handler {
	return &Handler{
		Setting: settings,
	}
}

func (handler *Handler) WithHealthHandler() *Handler {
	handler.Healthhandler = NewHealthHandler()
	return handler
}

func (handler *Handler) WithGamehandler() *Handler {
	games, err := game.LoadFromDirectory(handler.Setting.GameConfigDirectory)
	if err != nil {
		log.Fatal().Err(err).Str("directory", handler.Setting.GameConfigDirectory).Msg("failed to load game configuration files")
	}
	log.Debug().Int("size", games.Size()).Msg("successfully loaded games from configuration files")
	handler.Gamehandler = NewGameHandler(handler, games)
	return handler
}

func (handler *Handler) WithUserhandler(ctx context.Context, errgrp *errgroup.Group) *Handler {
	handler.Userhandler = NewUserHandler(ctx, errgrp, handler)
	return handler
}

func (handler *Handler) WithDownloadHandler() *Handler {
	handler.Downloadhandler = NewDownloadHandler(handler)
	return handler
}

func (handler *Handler) WithChatHandler() *Handler {
	handler.Chathandler = NewChatHandler(handler)
	return handler
}

func (handler *Handler) WithFileHandler() *Handler {
	handler.Filehandler = NewFileHandler(handler)
	return handler
}

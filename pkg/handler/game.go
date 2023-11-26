package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/network"
)

type GameHandler struct {
	handler *Handler
	Games   map[string]game.Game
}

func (h *GameHandler) GetGames(w http.ResponseWriter, req *http.Request) {
	keys := make([]string, 0, len(h.Games))
	for k := range h.Games {
		keys = append(keys, k)
	}

	response, err := json.Marshal(keys)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to encode keys of Games map")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *GameHandler) GetGame(w http.ResponseWriter, req *http.Request) {
	game, err := h.getGame(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn().Err(err).Send()
		return
	}

	response, err := json.Marshal(game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msgf("Failed to encode game '%s'", game.Slug)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *GameHandler) GetGameDownload(w http.ResponseWriter, req *http.Request) {
	h.serveGameFile(w, req, h.handler.Setting.GameFileDirectory)
}

func (h *GameHandler) GetGameDownloadIcon(w http.ResponseWriter, req *http.Request) {
	h.serveGameFile(w, req, h.handler.Setting.GameIconDirectory)
}

func (h *GameHandler) getGame(req *http.Request) (game game.Game, err error) {
	vars := mux.Vars(req)
	slug := vars["slug"]

	game, found := h.Games[slug]
	if found == false {
		err = errors.New(fmt.Sprintf("No game '%s' found", slug))
		return
	}

	return
}

func (h *GameHandler) serveGameFile(w http.ResponseWriter, req *http.Request, directory string) {
	game, err := h.getGame(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file := filesystem.SearchFileByName(game.Slug, directory)[0]
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msgf("Failed to retrieve binary data for '%s' in '%s'", game.Slug, directory)
		return
	}

	err = network.ServeFileData(file, w, req)
	if err != nil {
		log.Warn().Err(err).Msgf("Can not serve file '%s'", file)
	}
}

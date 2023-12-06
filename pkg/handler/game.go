package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/network"
)

type Gamehandler struct {
	parent *Handler
	Games  game.Games
}

func (handler *Gamehandler) GetGames(w http.ResponseWriter, req *http.Request) {
	slugsjson, err := json.Marshal(handler.Games.Slugs())
	if err != nil {
		log.Error().Err(err).Msg("failed to encode game slug list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(slugsjson)
	log.Trace().RawJSON("slugs", slugsjson).Msg("GET - /games")
}

func (handler *Gamehandler) GetGame(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	slug := vars["slug"]
	game, err := handler.Games.Get(slug)
	if err != nil {
		log.Warn().Err(err).Str("slug", slug).Msg("game not available")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gamejson, err := json.Marshal(game)
	if err != nil {
		log.Error().Err(err).Str("slug", slug).Msg("failed to encode game")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(gamejson)
	log.Trace().RawJSON("game", gamejson).Msg("GET - /games/:slug")
}

func (handler *Gamehandler) GetGameDownload(w http.ResponseWriter, req *http.Request) {
	handler.serveGameFile(w, req, handler.parent.Setting.GameFileDirectory)
	if req.Method == http.MethodHead {
		log.Trace().Msg("HEAD - /games/:slug/download")
	} else if req.Method == http.MethodGet {
		log.Trace().Msg("GET - /games/:slug/download")
	}
}

func (handler *Gamehandler) GetGameDownloadIcon(w http.ResponseWriter, req *http.Request) {
	handler.serveGameFile(w, req, handler.parent.Setting.GameIconDirectory)
	if req.Method == http.MethodHead {
		log.Trace().Msg("HEAD - /games/:slug/download/icon")
	} else if req.Method == http.MethodGet {
		log.Trace().Msg("GET - /games/:slug/download/icon")
	}
}

func (handler *Gamehandler) serveGameFile(w http.ResponseWriter, req *http.Request, directory string) {
	vars := mux.Vars(req)
	slug := vars["slug"]
	if !handler.Games.HasGame(slug) {
		log.Warn().Str("slug", slug).Msg("no game available")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	files, err := filesystem.SearchFileByNameLazy(slug, directory)
	if err != nil {
		log.Error().Err(err).Str("slug", slug).Msg("failed to retrieve binary data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = network.ServeFileData(files[0], w, req)
	if err != nil {
		log.Warn().Err(err).Str("file", files[0]).Msg("failed to serve file / provide meta-info")
	}
}

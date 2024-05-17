package handler

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/network"
	"github.com/seternate/go-lanty/pkg/setting"
)

type Downloadhandler struct{}

func NewDownloadHandler(parent *Handler) (handler *Downloadhandler) {
	handler = &Downloadhandler{}
	return
}

func (handler *Downloadhandler) GetDownload(w http.ResponseWriter, req *http.Request) {
	files, err := filesystem.SearchFilesBreadthFirst(setting.CLIENT_DOWNLOAD_DIRECTORY, setting.CLIENT_DOWNLOAD_FILE, 0, 1)
	if err != nil || len(files) < 1 {
		log.Error().Err(err).Msg("failed to retrieve lanty client binary data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = network.ServeFileData(files[0], w, req)
	if err != nil {
		log.Warn().Err(err).Str("file", files[0]).Msg("failed to serve file / provide meta-info")
	}

	if req.Method == http.MethodHead {
		log.Trace().Msg("HEAD - /download")
	} else if req.Method == http.MethodGet {
		log.Trace().Msg("GET - /download")
	}
}

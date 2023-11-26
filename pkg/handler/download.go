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
	file, err := filesystem.SearchFileByName(setting.CLIENT_DOWNLOAD_DIRECTORY, setting.CLIENT_DOWNLOAD_FILE, 1)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve lanty client binary data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = network.ServeFileData(file, w, req)
	if err != nil {
		log.Warn().Err(err).Str("file", file).Msg("failed to serve file / provide meta-info")
	}

	if req.Method == http.MethodHead {
		log.Trace().Msg("HEAD - /download")
	} else if req.Method == http.MethodGet {
		log.Trace().Msg("GET - /download")
	}
}

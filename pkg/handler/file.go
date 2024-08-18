package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/file"
	"github.com/seternate/go-lanty/pkg/filesystem"
	"github.com/seternate/go-lanty/pkg/network"
)

type Filehandler struct {
	parent *Handler
}

func NewFileHandler(parent *Handler) (handler *Filehandler) {
	handler = &Filehandler{
		parent: parent,
	}
	return
}

func (handler *Filehandler) PostFile(w http.ResponseWriter, req *http.Request) {
	download, err := network.NewDownloadFromRaw(req.Header, req.Body)
	if err != nil {
		log.Error().Err(err).Msg("error creating download for file upload")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filepath := filepath.Join(handler.parent.Setting.FileUploadDirectory, download.Filename())
	if _, err := os.Stat(handler.parent.Setting.FileUploadDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(handler.parent.Setting.FileUploadDirectory, 0755)
		log.Debug().Str("directory", handler.parent.Setting.FileUploadDirectory).Msg("created missing file-upload directory")
		if err != nil {
			log.Error().Err(err).Msg("error creating file-upload directory")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	f, err := os.Create(filepath)
	if err != nil {
		log.Error().Err(err).Str("file", filepath).Msg("can not create file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	download.Download(req.Context(), f)
	if download.Err != nil {
		log.Error().Err(download.Err).Str("file", filepath).Msg("failed to download file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileurl := req.URL.JoinPath(download.Filename())
	fileurl.Host = req.Host
	if req.TLS != nil {
		fileurl.Scheme = "https"
	} else {
		fileurl.Scheme = "http"
	}
	response := file.FileUploadResponse{
		URL: fileurl.String(),
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Interface("response", response).Msg("failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)
}

func (handler *Filehandler) GetFile(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	filename := vars["filename"]

	files, err := filesystem.SearchFileByNameLazy(filename, handler.parent.Setting.FileUploadDirectory)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("failed to retrieve binary data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = network.ServeFileData(files[0], w, req)
	if err != nil {
		log.Warn().Err(err).Str("file", files[0]).Msg("failed to serve file / provide meta-info")
	}
}

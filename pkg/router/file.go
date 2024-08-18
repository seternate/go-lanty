package router

import (
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
)

func FileRoutes(handler *handler.Handler) (routes Routes) {
	routes = Routes{
		"PostFile":    Route{"PostFile", "POST", "/files", nil},
		"GetFileHEAD": Route{"GetFile", "HEAD", "/files/{filename}", nil},
		"GetFileGET":  Route{"GetFile", "GET", "/files/{filename}", nil},
	}

	if handler == nil || handler.Gamehandler == nil {
		log.Trace().Msg("no HandlerFunc added to routes for files")
		return
	}

	routes.UpdateHandlerFunc("PostFile", handler.Filehandler.PostFile)
	routes.UpdateHandlerFunc("GetFileHEAD", handler.Filehandler.GetFile)
	routes.UpdateHandlerFunc("GetFileGET", handler.Filehandler.GetFile)

	log.Trace().Msg("HandlerFunc added to routes for files")

	return
}

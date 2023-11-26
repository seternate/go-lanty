package router

import (
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
)

func DownloadRoutes(handler *handler.Handler) (routes Routes) {
	routes = Routes{
		"GetDownload": Route{"GetDownload", "GET", "/download", nil},
	}

	if handler == nil || handler.Downloadhandler == nil {
		log.Trace().Msg("no HandlerFunc added to routes for download")
		return
	}

	routes.UpdateHandlerFunc("GetDownload", handler.Downloadhandler.GetDownload)

	log.Trace().Msg("HandlerFunc added to routes for download")

	return
}

package router

import (
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
)

func GameRoutes(handler *handler.Handler) (routes Routes) {
	routes = Routes{
		"GetGames":                Route{"GetGames", "GET", "/games", nil},
		"GetGame":                 Route{"GetGame", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}", nil},
		"GetGameDownloadHEAD":     Route{"GetGameDownload", "HEAD", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download", nil},
		"GetGameDownloadGET":      Route{"GetGameDownload", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download", nil},
		"GetGameDownloadIconHEAD": Route{"GetGameDownloadIcon", "HEAD", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download/icon", nil},
		"GetGameDownloadIconGET":  Route{"GetGameDownloadIcon", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download/icon", nil},
	}

	if handler == nil || handler.Gamehandler == nil {
		log.Trace().Msg("no HandlerFunc added to routes for games")
		return
	}

	routes.UpdateHandlerFunc("GetGames", handler.Gamehandler.GetGames)
	routes.UpdateHandlerFunc("GetGame", handler.Gamehandler.GetGame)
	routes.UpdateHandlerFunc("GetGameDownloadHEAD", handler.Gamehandler.GetGameDownload)
	routes.UpdateHandlerFunc("GetGameDownloadGET", handler.Gamehandler.GetGameDownload)
	routes.UpdateHandlerFunc("GetGameDownloadIconHEAD", handler.Gamehandler.GetGameDownloadIcon)
	routes.UpdateHandlerFunc("GetGameDownloadIconGET", handler.Gamehandler.GetGameDownloadIcon)

	log.Trace().Msg("HandlerFunc added to routes for games")

	return
}

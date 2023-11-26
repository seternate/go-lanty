package router

import (
	"github.com/seternate/go-lanty/pkg/handler"
)

func GameRoutes(h *handler.Handler) (routes Routes) {
	routes = Routes{
		"GetGames":                Route{"GetGames", "GET", "/games", nil},
		"GetGame":                 Route{"GetGame", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}", nil},
		"GetGameDownloadHEAD":     Route{"GetGameDownload", "HEAD", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download", nil},
		"GetGameDownloadGET":      Route{"GetGameDownload", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download", nil},
		"GetGameDownloadIconHEAD": Route{"GetGameDownloadIcon", "HEAD", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download/icon", nil},
		"GetGameDownloadIconGET":  Route{"GetGameDownloadIcon", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download/icon", nil},
	}

	if h == nil || h.GameHandler == nil {
		return
	}

	routes.UpdateHandlerFunc("GetGames", h.GameHandler.GetGames)
	routes.UpdateHandlerFunc("GetGame", h.GameHandler.GetGame)
	routes.UpdateHandlerFunc("GetGameDownloadHEAD", h.GameHandler.GetGameDownload)
	routes.UpdateHandlerFunc("GetGameDownloadGET", h.GameHandler.GetGameDownload)
	routes.UpdateHandlerFunc("GetGameDownloadIconHEAD", h.GameHandler.GetGameDownloadIcon)
	routes.UpdateHandlerFunc("GetGameDownloadIconGET", h.GameHandler.GetGameDownloadIcon)

	return
}

package router

import "github.com/seternate/go-lanty/pkg/handler"

func GameRoutes(h *handler.Handler) (routes Routes) {
	handler := h.GameHandler
	if handler == nil {
		return
	}

	routes = Routes{
		Route{"GetGames", "GET", "/games", handler.GetGames},
		Route{"GetGame", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}", handler.GetGame},
		Route{"GetGameDownload", "HEAD", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download", handler.GetGameDownload},
		Route{"GetGameDownload", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download", handler.GetGameDownload},
		Route{"GetGameDownloadIcon", "HEAD", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download/icon", handler.GetGameDownloadIcon},
		Route{"GetGameDownloadIcon", "GET", "/games/{slug:(?:[a-z]|[0-9]|-)+}/download/icon", handler.GetGameDownloadIcon},
	}

	return
}

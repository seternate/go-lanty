package router

import (
	"github.com/seternate/go-lanty/pkg/handler"
)

func UserRoutes(h *handler.Handler) (routes Routes) {
	handler := h.UserHandler
	if handler == nil {
		return
	}

	routes = Routes{
		Route{"GetUsers", "GET", "/users", handler.GetUsers},
		Route{"GetUser", "GET", "/users/{name:(?:[a-z]|[0-9]|-)+}", handler.GetUser},
		Route{"PostUser", "POST", "/users", handler.PostUser},
	}

	return
}

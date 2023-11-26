package router

import (
	"github.com/seternate/go-lanty/pkg/handler"
)

func UserRoutes(h *handler.Handler) (routes Routes) {
	routes = Routes{
		"GetUsers": Route{"GetUsers", "GET", "/users", nil},
		"GetUser":  Route{"GetUser", "GET", "/users/{name:(?:[a-z]|[0-9]|-)+}", nil},
		"PostUser": Route{"PostUser", "POST", "/users", nil},
	}

	if h == nil || h.UserHandler == nil {
		return
	}

	routes.UpdateHandlerFunc("GetUsers", h.UserHandler.GetUsers)
	routes.UpdateHandlerFunc("GetUser", h.UserHandler.GetUser)
	routes.UpdateHandlerFunc("PostUser", h.UserHandler.PostUser)

	return
}

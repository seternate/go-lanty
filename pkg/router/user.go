package router

import (
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
)

func UserRoutes(handler *handler.Handler) (routes Routes) {
	routes = Routes{
		"GetUsers": Route{"GetUsers", "GET", "/users", nil},
		"GetUser":  Route{"GetUser", "GET", "/users/{ip:(?:[0-9]|.)+}", nil},
		"PostUser": Route{"PostUser", "POST", "/users", nil},
	}

	if handler == nil || handler.Userhandler == nil {
		log.Trace().Msg("no HandlerFunc added to routes for users")
		return
	}

	routes.UpdateHandlerFunc("GetUsers", handler.Userhandler.GetUsers)
	routes.UpdateHandlerFunc("GetUser", handler.Userhandler.GetUser)
	routes.UpdateHandlerFunc("PostUser", handler.Userhandler.PostUser)

	log.Trace().Msg("HandlerFunc added to routes for users")

	return
}

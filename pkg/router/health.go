package router

import (
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
)

func HealthRoutes(handler *handler.Handler) (routes Routes) {
	routes = Routes{
		"GetHealth": Route{"GetHealth", "GET", "/health", nil},
	}

	if handler == nil || handler.Healthhandler == nil {
		log.Trace().Msg("no HandlerFunc added to routes for health")
		return
	}

	routes.UpdateHandlerFunc("GetHealth", handler.Healthhandler.GetHealth)

	log.Trace().Msg("HandlerFunc added to routes for health")

	return
}

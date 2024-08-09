package router

import (
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
)

func ChatRoutes(handler *handler.Handler) (routes Routes) {
	routes = Routes{
		"Chat": Route{"Chat", "GET", "/chat", nil},
	}

	if handler == nil || handler.Chathandler == nil {
		log.Trace().Msg("no HandlerFunc added to routes for chat")
		return
	}

	routes.UpdateHandlerFunc("Chat", handler.Chathandler.Chat)

	log.Trace().Msg("HandlerFunc added to routes for chat")

	return
}

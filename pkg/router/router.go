package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Router struct {
	*mux.Router
	MiddlewareFunc func(http.Request)
}

func middlewareFunc(r http.Request) {
	log.Info().Msgf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)
}

func NewRouter() *Router {
	router := &Router{
		Router:         mux.NewRouter(),
		MiddlewareFunc: middlewareFunc,
	}

	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			router.MiddlewareFunc(*r)
			handler.ServeHTTP(w, r)
		})
	})

	router.StrictSlash(true)

	return router
}

func (router *Router) WithRoutes(routes Routes) *Router {
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
		log.Trace().Interface("route", route).Msg("added route to router")
	}
	return router
}

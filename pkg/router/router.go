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

func NewRouter() *Router {
	router := &Router{mux.NewRouter(), middlewareFunc}
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			router.MiddlewareFunc(*r)
			h.ServeHTTP(w, r)
		})
	})
	router.StrictSlash(true)
	return router
}

func (r *Router) WithRoutes(routes Routes) *Router {
	for _, route := range routes {
		r.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return r
}

func middlewareFunc(r http.Request) {
	log.Info().Msgf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

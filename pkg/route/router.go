package route

import (
	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func NewRouter() *Router {
	router := &Router{mux.NewRouter()}
	router.StrictSlash(true)
	return router
}

func (router *Router) AddRoutes(routes Routes) {
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

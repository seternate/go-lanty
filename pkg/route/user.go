package route

import (
	"github.com/seternate/go-lanty-server/pkg/handlers"
)

var UserRoutes = Routes{
	Route{"GetUsers", "GET", "/users", handlers.HandleGetUsers},
	Route{"GetUser", "GET", "/users/{userId:[0-9]+}", handlers.HandleGetUser},
	Route{"PostUser", "POST", "/users", handlers.HandlePostUser},
}

package route

import (
	"github.com/seternate/go-lanty-server/pkg/handlers"
)

var GameRoutes = Routes{
	Route{"GetGames", "GET", "/games", handlers.HandleGetGames},
	Route{"GetGame", "GET", "/games/{gameId:[0-9]+}", handlers.HandleGetGame},
	Route{"GetGameFile", "HEAD", "/games/{gameId:[0-9]+}/file", handlers.HandleGetGameFile},
	Route{"GetGameFile", "GET", "/games/{gameId:[0-9]+}/file", handlers.HandleGetGameFile},
	Route{"GetGameIcon", "HEAD", "/games/{gameId:[0-9]+}/icon", handlers.HandleGetGameIcon},
	Route{"GetGameIcon", "GET", "/games/{gameId:[0-9]+}/icon", handlers.HandleGetGameIcon},
	Route{"GetGameCover", "GET", "/games/{gameId:[0-9]+}/cover", handlers.HandleGetGameCover},
}

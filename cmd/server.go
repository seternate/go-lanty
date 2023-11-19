package main

import (
	"log"
	"net/http"

	"github.com/seternate/go-lanty-server/pkg/handlers"
	"github.com/seternate/go-lanty-server/pkg/route"
	"github.com/seternate/lanty-api-golang/pkg/game"
)

var games *game.Games

func main() {
	games = new(game.Games)
	games.LoadGames("game")

	handlers.InitGamesHandler(games)
	handlers.InitUsersHandler()

	router := route.NewRouter()
	router.AddRoutes(route.GameRoutes)
	router.AddRoutes(route.UserRoutes)

	log.Fatal(http.ListenAndServe(":8090", router))
}

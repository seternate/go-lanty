package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/handler"
	"github.com/seternate/go-lanty/pkg/router"
	"github.com/seternate/go-lanty/pkg/setting"
)

func main() {
	parseFlags()

	settings, err := setting.LoadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load settings")
	}
	log.Debug().Interface("settings", settings).Msg("successfully loaded settings")

	games, err := game.LoadFromDirectory(settings.GameConfigDirectory)
	if err != nil {
		log.Fatal().Err(err).Str("directory", settings.GameConfigDirectory).Msg("failed to load game configuration files")
	}
	log.Debug().Int("size", games.Size()).Msg("successfully loaded games from configuration files")

	handler := handler.NewHandler(&settings).
		WithGamehandler(games).
		WithUserhandler()
	log.Trace().Msg("handler created")
	gameRoutes := router.GameRoutes(handler)
	userRoutes := router.UserRoutes(handler)
	router := router.NewRouter().
		WithRoutes(gameRoutes).
		WithRoutes(userRoutes)
	log.Trace().Msg("router created")

	address := fmt.Sprintf(":%d", settings.ServerPort)
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create listener")
	}
	log.Fatal().Err(http.Serve(listener, router)).Msg("unexpected error of http server")
}

func parseFlags() {
	loglevel := flag.String("loglevel", "info", "Sets the log level of the application")
	flag.Parse()

	log.Trace().Str("loglevel", *loglevel).Msg("parsed flags")

	switch *loglevel {
	case "disable":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}
}

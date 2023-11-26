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
	logLevel := flag.String("loglevel", "info", "Sets the log level of the application")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	switch *logLevel {
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

	s, err := setting.LoadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load settings")
	}
	log.Debug().Interface("settings", s).Msg("Loaded settings")

	games, err := game.LoadFromDirectory(s.GameConfigDirectory)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to load game configuration files from directory '%s'", s.GameConfigDirectory)
	}
	log.Debug().Int("size", len(games)).Msg("Loaded games")

	handler := handler.NewHandler(&s).
		WithGameHandler(games).
		WithUserHandler()
	gameRoutes := router.GameRoutes(handler)
	userRoutes := router.UserRoutes(handler)
	r := router.NewRouter().
		WithRoutes(gameRoutes).
		WithRoutes(userRoutes)

	address := fmt.Sprintf(":%d", s.ServerPort)
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Fatal().Err(http.Serve(listener, r)).Send()
}

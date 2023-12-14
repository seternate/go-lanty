package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
	"github.com/seternate/go-lanty/pkg/logging"
	"github.com/seternate/go-lanty/pkg/router"
	"github.com/seternate/go-lanty/pkg/setting"
)

func main() {
	logconfig := logging.Config{ConsoleLoggingEnabled: true}
	parseFlags(&logconfig)
	log.Logger = logging.Configure(logconfig)

	settings, err := setting.LoadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load settings")
	}
	log.Debug().Interface("settings", settings).Msg("successfully loaded settings")

	handler := handler.NewHandler(&settings).
		WithGamehandler().
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

func parseFlags(config *logging.Config) {
	flag.StringVar(&config.LogLevel, "loglevel", "info", "Sets the log level")
	flag.BoolVar(&config.FileLoggingEnabled, "logenablefile", false, "Enables logging to file")
	flag.StringVar(&config.Filename, "logfile", "lanty.log", "Sets the log filename")
	flag.StringVar(&config.Directory, "logdir", "log", "Sets the log directory")
	flag.IntVar(&config.MaxBackups, "logbackups", 0, "Sets the number of old logs to remain")
	flag.IntVar(&config.MaxSize, "logfilesize", 10, "Sets the size of the logs before rotating to new file")
	flag.IntVar(&config.MaxAge, "logage", 0, "Sets the maximum number of days to retain old logs")
	flag.Parse()
}

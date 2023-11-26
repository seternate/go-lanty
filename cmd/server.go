package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/handler"
	"github.com/seternate/go-lanty/pkg/logging"
	"github.com/seternate/go-lanty/pkg/router"
	"github.com/seternate/go-lanty/pkg/setting"
	"golang.org/x/sync/errgroup"
)

//TODO debug why server hangs from time to time

func main() {
	signalCtx, cancelSignalCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelSignalCtx()
	errgrp, errCtx := errgroup.WithContext(signalCtx)

	settings, err := setting.LoadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load settings")
	}
	logconfig := logging.Config{ConsoleLoggingEnabled: true}

	parseConfigFlags(&logconfig)
	parseServerFlags(&settings)
	printVersion := flag.Bool("version", false, "Prints the version information")
	flag.Parse()

	log.Logger = logging.Configure(logconfig)

	if *printVersion {
		fmt.Printf("%s - %s", setting.VERSION, runtime.Version())
		os.Exit(0)
	}

	log.Debug().Interface("logconfig", logconfig).Msg("logging configuration")
	log.Debug().Interface("settings", settings).Msg("settings")

	handler := handler.NewHandler(&settings).
		WithGamehandler().
		WithUserhandler(errCtx, errgrp).
		WithDownloadHandler()
	log.Trace().Msg("handler created")
	gameRoutes := router.GameRoutes(handler)
	userRoutes := router.UserRoutes(handler)
	downloadRoutes := router.DownloadRoutes(handler)
	log.Trace().Msg("routes created")
	router := router.NewRouter().
		WithRoutes(gameRoutes).
		WithRoutes(userRoutes).
		WithRoutes(downloadRoutes)
	log.Trace().Msg("router created")

	listener, err := net.Listen("tcp4", ":"+strconv.Itoa(settings.ServerPort))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create listener")
	}
	httpServer := &http.Server{
		Handler: router,
		BaseContext: func(_ net.Listener) context.Context {
			return signalCtx
		},
	}

	var shutdownErr error
	errgrp.Go(func() error {
		log.Info().Msg("starting http server: " + listener.Addr().String())
		return httpServer.Serve(listener)
	})
	errgrp.Go(func() error {
		<-errCtx.Done()
		log.Info().Msgf("graceful shutdown of http server (%ds)", settings.ServerGracefulShutdown)
		ctx, ctxCancel := context.WithTimeout(context.Background(), time.Duration(settings.ServerGracefulShutdown)*time.Second)
		defer ctxCancel()
		shutdownErr = httpServer.Shutdown(ctx)
		return shutdownErr
	})

	err = errgrp.Wait()
	if err != nil || shutdownErr != nil {
		if err == http.ErrServerClosed && shutdownErr != nil {
			log.Fatal().Err(shutdownErr).Msg("server closed forcefully")
		} else if err != http.ErrServerClosed && err != context.Canceled && shutdownErr == nil {
			log.Fatal().Err(err).Msg("unexpected http server error")
		}
		log.Info().Msg("server closed")
	}
}

func parseConfigFlags(config *logging.Config) {
	flag.StringVar(&config.LogLevel, "loglevel", "info", "Sets the log level")
	flag.BoolVar(&config.FileLoggingEnabled, "logenablefile", false, "Enables logging to file")
	flag.StringVar(&config.Filename, "logfile", "lanty.log", "Sets the log filename")
	flag.StringVar(&config.Directory, "logdir", "log", "Sets the log directory")
	flag.IntVar(&config.MaxBackups, "logbackups", 0, "Sets the number of old logs to remain")
	flag.IntVar(&config.MaxSize, "logfilesize", 10, "Sets the size of the logs before rotating to new file")
	flag.IntVar(&config.MaxAge, "logage", 0, "Sets the maximum number of days to retain old logs")
}

func parseServerFlags(settings *setting.Settings) {
	flag.IntVar(&settings.ServerPort, "port", settings.ServerPort, "Port of the http server to listen on")
	flag.IntVar(&settings.ServerGracefulShutdown, "graceful-shutdown", 10, "Timeout in seconds to wait for a graceful shutdown of the server")
	flag.StringVar(&settings.GameConfigDirectory, "game-config-dir", settings.GameConfigDirectory, "Directory of the game configuration files")
	flag.StringVar(&settings.GameFileDirectory, "game-file-dir", settings.GameFileDirectory, "Directory of the game files")
}
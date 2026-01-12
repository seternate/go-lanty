package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	appassetsrv "github.com/seternate/go-lanty/pkg/application/asset"
	appgamesrv "github.com/seternate/go-lanty/pkg/application/game"
	"github.com/seternate/go-lanty/pkg/infrastructure/checksum"
	"github.com/seternate/go-lanty/pkg/infrastructure/database"
	"github.com/seternate/go-lanty/pkg/infrastructure/mimetype"
	"github.com/seternate/go-lanty/pkg/infrastructure/persistence"
	"github.com/seternate/go-lanty/pkg/infrastructure/storageadapter"
	"github.com/seternate/go-lanty/pkg/interface/http/adapter/router"
	"github.com/seternate/go-lanty/pkg/interface/http/adapter/server"
	"github.com/seternate/go-lanty/pkg/interface/http/adapter/swagger"
	"github.com/seternate/go-lanty/pkg/interface/http/controller"
	"github.com/seternate/go-lanty/pkg/logging"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
)

var AppVersion = "dev-build"
var APIVersion = "v1.0.0"
var BasePath = "/api/v1"

// @title Lanty
// @version dev
// @description Lanty is a platform for managing and serving games for LAN parties
// @contact.name Levin Jeck
// @contact.url https://github.com/seternate/go-lanty
// @contact.email seternate@gmail.com
// @license.name License - MIT
// @license.url https://github.com/seternate/go-lanty/blob/main/LICENSE.md
// @externalDocs.description Documentation
// @externalDocs.url https://github.com/seternate/go-lanty/blob/main/README.md
func main() {
	signalCtx, cancelSignalCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancelSignalCtx()
	errgrp, errCtx := errgroup.WithContext(signalCtx)

	config, flagset, err := ParseConfig(os.Args...)
	if err != nil {
		if flagset != nil {
			fmt.Printf("%s\n\n", err)
			flagset.Usage()
			os.Exit(1)
		}
		log.Fatal().Err(err).Msg("error parsing flagset for configuration")
	}

	if config.PrintHelp {
		flagset.Usage()
		os.Exit(0)
	}
	if config.PrintVersion {
		fmt.Printf("%s - %s\n", AppVersion, runtime.Version())
		os.Exit(0)
	}

	log.Logger = logging.Configure(logging.Config{
		LogLevel:           config.LogLevel,
		FileLoggingEnabled: config.EnableFileLogging,
		Directory:          filepath.Dir(filepath.Dir(os.Args[0])),
		Filename:           "lanty.log",
		MaxSize:            config.LogFileSize,
		MaxBackups:         config.LogBackups,
		MaxAge:             config.LogAge,
	})

	db, err := database.New(config.DBString)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	//TODO: FINISH SETUP

	persistance := persistence.NewRepositories(db)

	// Ensure storage directories exist
	osFS := afero.NewOsFs()
	err = osFS.MkdirAll("./blob", 0755)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create blob directory")
	}
	err = osFS.MkdirAll("./icon", 0755)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create icon directory")
	}

	gamefileFS := afero.NewBasePathFs(osFS, "./blob")
	gamefileStorageAdapter := storageadapter.NewFilesystemStorageAdapter(gamefileFS)
	gamefileCalculator := checksum.NewCalculator()
	mimeTypeDetector := mimetype.NewDetector()
	appfileassetservice := appassetsrv.NewService(
		appassetsrv.NewQueryService(db, gamefileStorageAdapter),
		appassetsrv.NewCommandService(persistance.Asset, gamefileStorageAdapter, gamefileCalculator, mimeTypeDetector),
	)

	gameiconFS := afero.NewBasePathFs(osFS, "./icon")
	gameIconStorageAdapter := storageadapter.NewFilesystemStorageAdapter(gameiconFS)
	gameiconCalculator := checksum.NewCalculator()
	appiconassetservice := appassetsrv.NewService(
		appassetsrv.NewQueryService(db, gameIconStorageAdapter),
		appassetsrv.NewCommandService(persistance.Asset, gameIconStorageAdapter, gameiconCalculator, mimeTypeDetector),
	)

	appgamecommandservice := appgamesrv.NewCommandService(
		persistance.Game,
		appiconassetservice,
		appfileassetservice,
	)
	appgameservice := appgamesrv.NewService(
		appgamesrv.NewQueryService(db, appiconassetservice, appfileassetservice),
		appgamecommandservice,
	)

	controller := controller.New(appgameservice)
	engine := router.New(controller)
	swagger.InitInfo(APIVersion, config.Host, config.Port, router.APIBasePath, []string{config.Scheme})
	httpserver := server.Init(errCtx, engine)

	if config.Scheme == "http" {
		errgrp.Go(func() error {
			log.Info().Int("port", config.Port).Msg("starting http server")
			return httpserver.Run(config.Port)
		})
	} else {
		log.Fatal().Msgf("unsupported scheme: %s", config.Scheme)
	}

	//Gracefully shutdown http server
	errgrp.Go(func() error {
		<-errCtx.Done()
		timeout := 10 * time.Second
		log.Info().Str("timeout", timeout.String()).Msg("graceful shutdown of the http server")
		ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
		defer ctxCancel()
		serr := httpserver.Shutdown(ctx)
		if serr != nil {
			log.Error().Err(serr).Msg("http server closed forcefully")
		} else {
			log.Info().Msg("http server closed gracefully")
		}
		return serr
	})

	//Waits for any os.Signal or an error in the buisness loop
	err = errgrp.Wait()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("unexpected application error")
	}
}

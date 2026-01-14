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
	appassetsrv "github.com/seternate/go-lanty/internal/application/asset"
	appgamesrv "github.com/seternate/go-lanty/internal/application/game"
	appusersrv "github.com/seternate/go-lanty/internal/application/user"
	"github.com/seternate/go-lanty/internal/infrastructure/checksum"
	"github.com/seternate/go-lanty/internal/infrastructure/database"
	"github.com/seternate/go-lanty/internal/infrastructure/mimetype"
	"github.com/seternate/go-lanty/internal/infrastructure/persistence"
	"github.com/seternate/go-lanty/internal/infrastructure/storageadapter"
	"github.com/seternate/go-lanty/internal/interface/http/adapter/router"
	"github.com/seternate/go-lanty/internal/interface/http/adapter/server"
	"github.com/seternate/go-lanty/internal/interface/http/adapter/swagger"
	"github.com/seternate/go-lanty/internal/interface/http/controller"
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
	persistance := persistence.NewRepositories(db)

	checksumCalculator := checksum.NewCalculator()
	mimeTypeDetector := mimetype.NewDetector()

	osFS := afero.NewOsFs()
	directories := map[string]string{"blob": "./blob", "icon": "./icon"}
	for dir, path := range directories {
		err = osFS.MkdirAll(path, 0755)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to create %s directory", dir)
		}
	}

	gamefileStorageAdapter := storageadapter.NewFilesystemStorageAdapter(
		afero.NewBasePathFs(osFS, directories["blob"]),
	)
	gameiconStorageAdapter := storageadapter.NewFilesystemStorageAdapter(
		afero.NewBasePathFs(osFS, directories["icon"]),
	)

	appfileassetservice := appassetsrv.NewService(
		appassetsrv.NewQueryService(db, gamefileStorageAdapter),
		appassetsrv.NewCommandService(persistance.Asset, gamefileStorageAdapter, checksumCalculator, mimeTypeDetector),
	)
	appiconassetservice := appassetsrv.NewService(
		appassetsrv.NewQueryService(db, gameiconStorageAdapter),
		appassetsrv.NewCommandService(persistance.Asset, gameiconStorageAdapter, checksumCalculator, mimeTypeDetector),
	)

	appgameservice := appgamesrv.NewService(
		appgamesrv.NewQueryService(
			db,
			appiconassetservice,
			appfileassetservice,
		),
		appgamesrv.NewCommandService(
			persistance.Game,
			appiconassetservice,
			appfileassetservice,
		),
	)

	appuserservice := appusersrv.NewService(
		appusersrv.NewQueryService(db),
		appusersrv.NewCommandService(persistance.User),
	)

	controller := controller.New(appgameservice, appuserservice)
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

	err = errgrp.Wait()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("unexpected application error")
	}
}

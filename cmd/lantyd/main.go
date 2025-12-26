package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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
)

var AppVersion = "dev-build"
var APIVersion = "v1.0.0"
var BasePath = "/api/v1"

//TODO: override Version in build with "-ldflags "-X main.Version=1.0.0""

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
	//SETUP SIGNAL
	// signalCtx, cancelSignalCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// defer cancelSignalCtx()
	// errgrp, errCtx := errgroup.WithContext(signalCtx)

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
	if config.Scheme == "http" {
		err = server.Init(engine).Run(config.Port)
		if err != nil {
			log.Fatal().Err(err).Msg("unexpected error running the http server")
		}
		os.Exit(0)
	}

	log.Fatal().Msgf("can not run server with scheme: %s", config.Scheme)

	//Setup database connection
	// db, err := sqlx.Connect("pgx", config.DBString)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to connect to database")
	// }
	// defer db.Close()
	// log.Debug().Msgf("successfully connected to db: %s", config.DBString)

	//Setup FS
	//TODO: Configurable path

	//Setup buisness logic
	// gamerepository := game.NewPostgresGameRepository(db)
	// gameservice := game.NewService(gamerepository, iconFS)
	//

	//Setup Swagger API
	// docs.SwaggerInfo.Title = "Lanty"
	// docs.SwaggerInfo.Description = "Dummy description"
	// docs.SwaggerInfo.Version = APIVersion
	// docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", config.Port)
	// docs.SwaggerInfo.InfoInstanceName = ""
	// docs.SwaggerInfo.Schemes = []string{"http"}
	// docs.SwaggerInfo.BasePath = BasePath

	// //Setup Handler
	// //TODO: Can be pulled to a router package
	// //TODO: Authentication / Authorization
	// if zerolog.GlobalLevel() > zerolog.DebugLevel {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	// r := gin.New()
	// r.Use(adapter.Wrap(hlog.NewHandler(log.Logger))).
	// 	Use(adapter.Wrap(hlog.RequestIDHandler("requestID", "X-Request-ID"))).
	// 	Use(adapter.Wrap(hlog.RemoteAddrHandler("remote"))).
	// 	Use(adapter.Wrap(hlog.RequestHandler("request"))).
	// 	Use(adapter.Wrap(hlog.ProtoHandler("proto"))).
	// 	Use(server.GinMiddlewareLogger())
	// r.GET("/health", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
	// r.GET("/docs/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	// r.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusMovedPermanently, "/docs/index.html") })
	// api := r.Group(BasePath)
	// //TODO: More generic call: r.Handle(...)
	// api.GET("/games", game.GetGames(gameservice))
	// api.GET("/games/:slug", game.GetGameBySlug(gameservice))
	// api.GET("/games/:slug/icon", game.GetIcon(gameservice))

	//Setup http server
	// server := http.Server{
	// 	Handler: r.Handler(),
	// 	BaseContext: func(net.Listener) context.Context {
	// 		return errCtx
	// 	},
	// }

	//Setup listener for explicit IPv4
	// listener, err := net.Listen("tcp4", ":"+strconv.Itoa(config.Port))
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to create listener")
	// }

	//Run http server
	// errgrp.Go(func() error {
	// 	log.Info().Str("address", listener.Addr().String()).Msg("starting http server")
	// 	return server.Serve(listener)
	// })

	//Gracefully shutdown http server
	// errgrp.Go(func() error {
	// 	<-errCtx.Done()
	// 	timeout := 10 * time.Second
	// 	log.Info().Str("timeout", timeout.String()).Msg("try to gracefully shutdown http server")
	// 	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	// 	defer ctxCancel()
	// 	serr := server.Shutdown(ctx)
	// 	if serr != nil {
	// 		log.Error().Err(serr).Msg("http server closed forcefully")
	// 	} else {
	// 		log.Info().Msg("http server closed gracefully")
	// 	}
	// 	return serr
	// })

	//Waits for any os.Signal or an error in the buisness loop
	// err = errgrp.Wait()
	// if err != nil && err != http.ErrServerClosed {
	// 	log.Fatal().Err(err).Msg("unexpected application error")
	// }
}

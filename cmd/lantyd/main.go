package main

import (
	"os"
	"path/filepath"

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
	flag "github.com/spf13/pflag"
)

//go:generate swag init -d ./ --output ./../../docs --outputTypes go,yaml -pd

var AppVersion = "dev-build"
var APIVersion = "v1.0.0"
var BasePath = "/api/v1"

//TODO: change .github folder for build process
//TODO: change readme
//TODO: override Version in build with "-ldflags "-X main.Version=1.0.0""
//TODO: test for integration tests
//TODO: unit-tests
//TODO: godocs

func main() {
	//SETUP SIGNAL
	// signalCtx, cancelSignalCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// defer cancelSignalCtx()
	// errgrp, errCtx := errgroup.WithContext(signalCtx)

	//LOAD EXTERNAL CONFIG (APP & DB & SERVER)
	// config, err := config.Parse(os.Args[1:]...)
	// dbstring := flag.String("db", "", "[REQUIRED] Connection string to a postgres database (eg. 'postgres://lanty:lanty@localhost:5432/lanty?sslmode=disable') [ENV: LANTY_DB]")
	// port := flag.IntP("port", "p", 8080, "Port of the server")
	// loglevel := flag.String("loglevel", "info", "Log level of the application [disable, trace, debug, info, warning, error, panic, fatal]")
	// printHelp := flag.BoolP("help", "h", false, "Prints help")
	// printVersion := flag.Bool("version", false, "Prints version information")
	// flag.Parse()

	//SETUP CLI
	// flag.Usage = func() {
	// 	fmt.Printf("Usage: %s [options]\nOptions:\n", os.Args[0])
	// 	flag.PrintDefaults()
	// }

	// if len(dbstring) == 0 {
	// 	log.Fatal().Msg("missing db connection string")
	// 	flag.Usage()
	// }

	// if err != nil {
	// 	log.Fatal().Err(err).Msg("error parsing/validating configuration")
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	// if config.PrintHelp {
	// 	flag.Usage()
	// 	os.Exit(0)
	// }
	// if config.PrintVersion {
	// 	fmt.Printf("%s - %s", AppVersion, runtime.Version())
	// 	os.Exit(0)
	// }

	//SETUP LOGGING
	executable, err := os.Executable()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting application path")
		flag.Usage()
		os.Exit(1)
	}

	log.Logger = logging.Configure(logging.Config{
		LogLevel:           "trace",
		FileLoggingEnabled: true,
		Directory:          filepath.Dir(executable),
		Filename:           "lanty.log",
		MaxSize:            10,
		MaxBackups:         3,
		MaxAge:             0,
	})

	//CREATE APPLICATION <- app-config & db-config
	db, err := database.New("postgres://lanty:lanty@localhost:5432/lanty?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

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

	//START SERVER <- server-config
	controller := controller.New(appgameservice)
	engine := router.New(controller)
	swagger.InitInfo("v1.0.0", router.APIBasePath, 8080)
	server.Init(engine).Run(8080)

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

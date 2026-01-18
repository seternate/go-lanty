package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/internal/interface/http/adapter/middleware"
	"github.com/seternate/go-lanty/internal/interface/http/controller"
	"github.com/seternate/go-lanty/pkg/api/paths"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func New(controller *controller.HTTPController) *gin.Engine {
	if zerolog.GlobalLevel() > zerolog.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	//TODO: Do we want to enable this?
	// binding.EnableDecoderDisallowUnknownFields = true
	router.RedirectTrailingSlash = true

	router.Use(
		gin.Recovery(),
		adapter.Wrap(hlog.NewHandler(log.Logger)),
		adapter.Wrap(hlog.RequestIDHandler("requestID", "X-Request-ID")),
		adapter.Wrap(hlog.RemoteAddrHandler("remote")),
		adapter.Wrap(hlog.RequestHandler("request")),
		adapter.Wrap(hlog.ProtoHandler("proto")),
		middleware.Logger,
		middleware.CurlNewlineAppender,
		middleware.ErrorHandler,
	)

	addSwaggerRoutes(router)
	addHealthRoutes(router)
	addAPIRoutes(router, controller)

	return router
}

func addSwaggerRoutes(router *gin.Engine) {
	router.GET("/docs/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusMovedPermanently, "/docs/index.html") })
}

func addHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "ok"}) })
}

func addAPIRoutes(r *gin.Engine, controller *controller.HTTPController) {
	api := r.Group(paths.APIBasePath)

	apiPaths := map[paths.Path][]gin.HandlerFunc{
		paths.GetGames:    {controller.Game.GetGames},
		paths.GetGame:     {controller.Game.GetGame},
		paths.PutGame:     {controller.Game.PutGame},
		paths.DeleteGame:  {controller.Game.DeleteGame},
		paths.GetGameIcon: {controller.Game.GetIcon},
		paths.PutGameIcon: {controller.Game.PutIcon},
		paths.GetGameBlob: {controller.Game.GetBlob},
		paths.PutGameBlob: {controller.Game.PutBlob},
		paths.GetUsers:    {controller.User.GetUsers},
		paths.PutUser:     {controller.User.PutUser},
		paths.DeleteUser:  {controller.User.DeleteUser},
	}

	for path, handlers := range apiPaths {
		api.Handle(path.Method, path.Template, handlers...)
	}

}

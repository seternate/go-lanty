package controller

import (
	appGameSrv "github.com/seternate/go-lanty/internal/application/game"
	appUserSrv "github.com/seternate/go-lanty/internal/application/user"
	"github.com/seternate/go-lanty/internal/interface/http/controller/game"
	"github.com/seternate/go-lanty/internal/interface/http/controller/user"
)

type HTTPController struct {
	Game game.EndpointController
	User user.EndpointController
}

func New(gameSrv *appGameSrv.Service, userSrv *appUserSrv.Service) *HTTPController {
	return &HTTPController{
		Game: game.EndpointController{Service: gameSrv},
		User: user.EndpointController{Service: userSrv},
	}
}

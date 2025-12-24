package controller

import (
	appGameSrv "github.com/seternate/go-lanty/pkg/application/game"
	"github.com/seternate/go-lanty/pkg/interface/http/controller/game"
)

type HTTPController struct {
	Game game.EndpointController
}

func New(gameSrv *appGameSrv.Service) *HTTPController {
	return &HTTPController{
		Game: game.EndpointController{Service: gameSrv},
	}
}

package game

import (
	appGameSrv "github.com/seternate/go-lanty/internal/application/game"
)

type EndpointController struct {
	Service *appGameSrv.Service
}

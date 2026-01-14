package user

import (
	appUserSrv "github.com/seternate/go-lanty/internal/application/user"
)

type EndpointController struct {
	Service *appUserSrv.Service
}

package user

import (
	apimodel "github.com/seternate/go-lanty/pkg/api/models/user"
	appUserSrv "github.com/seternate/go-lanty/internal/application/user"
	domainuser "github.com/seternate/go-lanty/internal/domain/user"
)

func NewUser(userView *appUserSrv.UserView) *apimodel.User {
	return &apimodel.User{
		IPv4Address: userView.IPv4Address,
		Username:    userView.Username,
	}
}

func NewUserFromDomain(user *domainuser.User) *apimodel.User {
	return &apimodel.User{
		IPv4Address: user.IPv4Address.String(),
		Username:    user.Username,
	}
}

func ToCommand(req *apimodel.UpsertUserRequest, ipv4Address string) appUserSrv.UpsertUserCommand {
	return appUserSrv.UpsertUserCommand{
		IPv4Address: ipv4Address,
		Username:    req.Username,
	}
}

package model

import (
	appUserSrv "github.com/seternate/go-lanty/internal/application/user"
	domainuser "github.com/seternate/go-lanty/internal/domain/user"
)

type UserResponse struct {
	IPv4Address string `json:"ipv4Address"`
	Username    string `json:"username"`
}

func NewUserResponse(userView *appUserSrv.UserView) *UserResponse {
	return &UserResponse{
		IPv4Address: userView.IPv4Address,
		Username:    userView.Username,
	}
}

func NewUserResponseFromDomain(user *domainuser.User) *UserResponse {
	return &UserResponse{
		IPv4Address: user.IPv4Address.String(),
		Username:    user.Username,
	}
}

type UpsertUserRequest struct {
	Username string `json:"username"`
}

func (req *UpsertUserRequest) ToCommand(ipv4Address string) appUserSrv.UpsertUserCommand {
	return appUserSrv.UpsertUserCommand{
		IPv4Address: ipv4Address,
		Username:    req.Username,
	}
}

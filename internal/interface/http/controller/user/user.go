package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	errorx "github.com/seternate/go-lanty/internal/interface/http/error"
	httpmodel "github.com/seternate/go-lanty/internal/interface/http/model/user"
	errormodel "github.com/seternate/go-lanty/pkg/api/models/error"
	"github.com/seternate/go-lanty/pkg/api/models/user"
)

var _ = errormodel.Error{}
var _ = usermodel.User{}

// @Summary Get all Users
// @Description Get all available Users
// @Tags users
// @Produce json
// @Success 200 {array} usermodel.User
// @Failure 500 {object} errormodel.Error
// @Router /users [get]
func (ctl *EndpointController) GetUsers(ctx *gin.Context) {
	users, err := ctl.Service.Query.GetUsers()
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to get users: %w", err))
		return
	}

	response := make([]*usermodel.User, 0, len(users))
	for _, user := range users {
		response = append(response, httpmodel.NewUser(&user))
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Create or update a User
// @Description Create or update the given User.
// @Tags users
// @Param ipv4Address path string true "IPv4 Address"
// @Param request body usermodel.UpsertUserRequest true "User to create or update"
// @Accept json
// @Produce json
// @Success 201 {object} usermodel.User "User created"
// @Success 202 {object} usermodel.User "User updated"
// @Failure 400 {object} errormodel.Error
// @Failure 500 {object} errormodel.Error
// @Router /users/{ipv4Address} [put]
func (ctl *EndpointController) PutUser(ctx *gin.Context) {
	ipv4Address := ctx.Param("ipv4Address")

	upsertRequest := &usermodel.UpsertUserRequest{}
	err := ctx.BindJSON(upsertRequest)
	if err != nil {
		errorx.AbortWithError(ctx, errorx.ErrBadRequest("failed to read request body").WithCause(err))
		return
	}

	user, created, err := ctl.Service.Command.UpsertUser(httpmodel.ToCommand(upsertRequest, ipv4Address))
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to upsert user: %w", err))
		return
	}

	if created {
		ctx.JSON(http.StatusCreated, httpmodel.NewUserFromDomain(user))
		return
	}

	ctx.JSON(http.StatusAccepted, httpmodel.NewUserFromDomain(user))
}

// @Summary Delete a User
// @Description Delete a User
// @Tags users
// @Param ipv4Address path string true "IPv4 Address"
// @Success 204 "No Content"
// @Failure 404 {object} errormodel.Error
// @Failure 500 {object} errormodel.Error
// @Router /users/{ipv4Address} [delete]
func (ctl *EndpointController) DeleteUser(ctx *gin.Context) {
	ipv4Address := ctx.Param("ipv4Address")

	err := ctl.Service.Command.DeleteUser(ipv4Address)
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to delete user: %w", err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

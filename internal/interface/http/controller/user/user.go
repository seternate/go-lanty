package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	errorx "github.com/seternate/go-lanty/internal/interface/http/error"
	errormodel "github.com/seternate/go-lanty/internal/interface/http/model/error"
	model "github.com/seternate/go-lanty/internal/interface/http/model/user"
)

var _ = errormodel.ErrorResponse{}

// @Summary Get all Users
// @Description Get all available Users
// @Tags users
// @Produce json
// @Success 200 {array} model.UserResponse
// @Failure 500 {object} errormodel.ErrorResponse
// @Router /users [get]
func (ctl *EndpointController) GetUsers(ctx *gin.Context) {
	users, err := ctl.Service.Query.GetUsers()
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to get users: %w", err))
		return
	}

	response := make([]*model.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, model.NewUserResponse(&user))
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Create or update a User
// @Description Create or update the given User.
// @Tags users
// @Param ipv4Address path string true "IPv4 Address"
// @Param request body model.UpsertUserRequest true "User to create or update"
// @Accept json
// @Produce json
// @Success 201 {object} model.UserResponse "User created"
// @Success 202 {object} model.UserResponse "User updated"
// @Failure 400 {object} errormodel.ErrorResponse
// @Failure 500 {object} errormodel.ErrorResponse
// @Router /users/{ipv4Address} [put]
func (ctl *EndpointController) PutUser(ctx *gin.Context) {
	ipv4Address := ctx.Param("ipv4Address")

	upsertRequest := &model.UpsertUserRequest{}
	err := ctx.BindJSON(upsertRequest)
	if err != nil {
		errorx.AbortWithError(ctx, errorx.ErrBadRequest("failed to read request body").WithCause(err))
		return
	}

	user, created, err := ctl.Service.Command.UpsertUser(upsertRequest.ToCommand(ipv4Address))
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to upsert user: %w", err))
		return
	}

	if created {
		ctx.JSON(http.StatusCreated, model.NewUserResponseFromDomain(user))
		return
	}

	ctx.JSON(http.StatusAccepted, model.NewUserResponseFromDomain(user))
}

// @Summary Delete a User
// @Description Delete a User
// @Tags users
// @Param ipv4Address path string true "IPv4 Address"
// @Success 204 "No Content"
// @Failure 404 {object} errormodel.ErrorResponse
// @Failure 500 {object} errormodel.ErrorResponse
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

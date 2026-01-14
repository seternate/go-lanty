package game

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	errorx "github.com/seternate/go-lanty/internal/interface/http/error"
	errmodel "github.com/seternate/go-lanty/internal/interface/http/model"
	model "github.com/seternate/go-lanty/internal/interface/http/model/game"
)

var _ = errmodel.ErrorResponse{}

// @Summary Get all Games
// @Description Get all available Games
// @Tags games
// @Produce json
// @Success 200 {array} errmodel.GameResponse
// @Failure 500 {object} errmodel.ErrorResponse
// @Router /games [get]
func (ctl *EndpointController) GetGames(ctx *gin.Context) {
	games, err := ctl.Service.Query.GetGames()
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to get games: %w", err))
		return
	}

	response := make([]*model.GameResponse, 0, len(games))
	for _, game := range games {
		response = append(response, model.NewGameResponse(&game))
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get a Game
// @Description Get a Game by its slug
// @Tags games
// @Param slug path string true "Slug"
// @Produce json
// @Success 200 {object} model.GameResponse
// @Failure 404 {object} errmodel.ErrorResponse
// @Failure 500 {object} errmodel.ErrorResponse
// @Router /games/{slug} [get]
func (ctl *EndpointController) GetGame(ctx *gin.Context) {
	slug := ctx.Param("slug")

	game, err := ctl.Service.Query.GetGame(slug)
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to get game: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, model.NewGameResponse(game))
}

// @Summary Create or update a Game
// @Description Create or update the given Game.
// @Tags games
// @Param slug path string true "Slug"
// @Param request body model.UpsertGameRequest true "Game to create or update"
// @Accept json
// @Produce json
// @Success 201 {object} model.GameResponse "Game created"
// @Success 202 {object} model.GameResponse "Game updated"
// @Failure 400 {object} errmodel.ErrorResponse
// @Failure 500 {object} errmodel.ErrorResponse
// @Router /games/{slug} [put]
func (ctl *EndpointController) PutGame(ctx *gin.Context) {
	slug := ctx.Param("slug")

	upsertRequest := &model.UpsertGameRequest{}
	err := ctx.BindJSON(upsertRequest)
	if err != nil {
		errorx.AbortWithError(ctx, errorx.ErrBadRequest("failed to read request body").WithCause(err))
		return
	}

	_, created, err := ctl.Service.Command.UpsertGame(upsertRequest.ToCommand(slug))
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to upsert game: %w", err))
		return
	}

	gameView, err := ctl.Service.Query.GetGame(slug)
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to get game: %w", err))
		return
	}

	if created {
		ctx.JSON(http.StatusCreated, model.NewGameResponse(gameView))
		return
	}

	ctx.JSON(http.StatusAccepted, model.NewGameResponse(gameView))
}

// @Summary Delete a Game
// @Description Delete a Game
// @Tags games
// @Param slug path string true "Slug"
// @Success 204 "No Content"
// @Failure 404 {object} errmodel.ErrorResponse
// @Failure 500 {object} errmodel.ErrorResponse
// @Router /games/{slug} [delete]
func (ctl *EndpointController) DeleteGame(ctx *gin.Context) {
	slug := ctx.Param("slug")

	err := ctl.Service.Command.DeleteGame(slug)
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to delete game: %w", err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

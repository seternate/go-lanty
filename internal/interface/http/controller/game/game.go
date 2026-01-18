package game

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	errorx "github.com/seternate/go-lanty/internal/interface/http/error"
	httpmodel "github.com/seternate/go-lanty/internal/interface/http/model/game"
	errormodel "github.com/seternate/go-lanty/pkg/api/models/error"
	gamemodel "github.com/seternate/go-lanty/pkg/api/models/game"
)

var _ = errormodel.Error{}
var _ = gamemodel.Game{}

// @Summary Get all Games
// @Description Get all available Games
// @Tags games
// @Produce json
// @Success 200 {array} gamemodel.Game
// @Failure 500 {object} errormodel.Error
// @Router /games [get]
func (ctl *EndpointController) GetGames(ctx *gin.Context) {
	games, err := ctl.Service.Query.GetGames()
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to get games: %w", err))
		return
	}

	response := make([]*gamemodel.Game, 0, len(games))
	for _, game := range games {
		response = append(response, httpmodel.NewGame(&game))
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get a Game
// @Description Get a Game by its slug
// @Tags games
// @Param slug path string true "Slug"
// @Produce json
// @Success 200 {object} gamemodel.Game
// @Failure 404 {object} errormodel.Error
// @Failure 500 {object} errormodel.Error
// @Router /games/{slug} [get]
func (ctl *EndpointController) GetGame(ctx *gin.Context) {
	slug := ctx.Param("slug")

	game, err := ctl.Service.Query.GetGame(slug)
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to get game: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, httpmodel.NewGame(game))
}

// @Summary Create or update a Game
// @Description Create or update the given Game.
// @Tags games
// @Param slug path string true "Slug"
// @Param request body gamemodel.UpsertGameRequest true "Game to create or update"
// @Accept json
// @Produce json
// @Success 201 {object} gamemodel.Game "Game created"
// @Success 202 {object} gamemodel.Game "Game updated"
// @Failure 400 {object} errormodel.Error
// @Failure 500 {object} errormodel.Error
// @Router /games/{slug} [put]
func (ctl *EndpointController) PutGame(ctx *gin.Context) {
	slug := ctx.Param("slug")

	upsertRequest := &gamemodel.UpsertGameRequest{}
	err := ctx.BindJSON(upsertRequest)
	if err != nil {
		errorx.AbortWithError(ctx, errorx.ErrBadRequest("failed to read request body").WithCause(err))
		return
	}

	_, created, err := ctl.Service.Command.UpsertGame(httpmodel.ToCommand(upsertRequest, slug))
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
		ctx.JSON(http.StatusCreated, httpmodel.NewGame(gameView))
		return
	}

	ctx.JSON(http.StatusAccepted, httpmodel.NewGame(gameView))
}

// @Summary Delete a Game
// @Description Delete a Game
// @Tags games
// @Param slug path string true "Slug"
// @Success 204 "No Content"
// @Failure 404 {object} errormodel.Error
// @Failure 500 {object} errormodel.Error
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

package game

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seternate/go-lanty/pkg/interface/http/adapter/handler"
	model "github.com/seternate/go-lanty/pkg/interface/http/model/game"
)

// @Summary Get list of Games
// @Description Get a list of all available Games
// @Tags games
// @Produce json
// @Success 200 {array} model.GameResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/games [get]
func (ctl *EndpointController) GetGames(ctx *gin.Context) {
	games, err := ctl.Service.Query.GetGames()
	if err != nil {
		handler.AbortWithError(ctx, err)
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
// @Param slug path string true "Slug of the Game"
// @Produce json
// @Success 200 {object} model.GameResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/games/{slug} [get]
func (ctl *EndpointController) GetGame(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if len(slug) == 0 {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing mandatory parameter: %s", "slug"))
		return
	}

	game, err := ctl.Service.Query.GetGame(slug)
	if err != nil {
		handler.AbortWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, model.NewGameResponse(game))
}

// @Summary Upsert a Game
// @Description Update or insert the given Game
// @Tags games
// @Param slug path string true "Slug of the Game"
// @Param request body model.UpsertGameRequest true "Game to update or insert"
// @Accept json
// @Produce json
// @Success 201 {object} model.GameResponse "Game created"
// @Success 202 {object} model.GameResponse "Game updated"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/games/{slug} [put]
func (ctl *EndpointController) PutGame(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if len(slug) == 0 {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing mandatory parameter %q", "slug"))
		return
	}

	upsertRequest := &model.UpsertGameRequest{}
	err := ctx.BindJSON(upsertRequest)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cmd := upsertRequest.ToCommand(slug)

	_, created, err := ctl.Service.Command.UpsertGame(cmd)
	if err != nil {
		handler.AbortWithError(ctx, err)
		return
	}

	gameView, err := ctl.Service.Query.GetGame(cmd.Slug)
	if err != nil {
		handler.AbortWithError(ctx, err)
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
// @Param slug path string true "Slug of the Game"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/games/{slug} [delete]
func (ctl *EndpointController) DeleteGame(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if len(slug) == 0 {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing mandatory parameter %q", "slug"))
		return
	}

	err := ctl.Service.Command.DeleteGame(slug)
	if err != nil {
		handler.AbortWithError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

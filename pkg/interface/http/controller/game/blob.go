package game

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appGameSrv "github.com/seternate/go-lanty/pkg/application/game"
	domainGame "github.com/seternate/go-lanty/pkg/domain/game"
	"github.com/seternate/go-lanty/pkg/interface/http/adapter/handler"
	"github.com/seternate/go-lanty/pkg/interface/http/adapter/header"
)

// @Summary Get blob for a Game
// @Description Get the blob for a Game with metadata headers
// @Tags games
// @Param slug path string true "Slug of the Game"
// @Produce application/octet-stream
// @Success 200 {file} binary "Blob binary data"
// @Header 200 {string} Content-Length "Size of the blob"
// @Header 200 {string} Content-Type "MIME type of the blob"
// @Header 200 {string} Content-Digest "Checksum in RFC 9530 format (algorithm=base64_checksum)"
// @Failure 400 {object} map[string]string "Bad request: missing slug parameter"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/games/{slug}/blob [get]
func (ctl *EndpointController) GetBlob(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if len(slug) == 0 {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing mandatory parameter %q", "slug"))
		return
	}

	assetContent, err := ctl.Service.Query.FetchBlob(slug)
	if err != nil {
		handler.AbortWithError(ctx, err)
		return
	}
	defer assetContent.Data.Close()

	// Format Content-Digest header
	digestHeader, err := header.EncodeContentDigestHeader(assetContent.Algorithm, assetContent.Checksum)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to format content-digest header: %w", err))
		return
	}

	ctx.DataFromReader(http.StatusOK, int64(assetContent.Size), assetContent.MimeType, assetContent.Data, map[string]string{
		"Content-Digest": digestHeader,
	})
}

// @Summary Upsert blob for a Game
// @Description Creates or updates the blob for a Game. Requires Content-Digest header with checksum. Accepts raw binary data. The Content-Digest header must match the checksum of the uploaded blob data.
// @Tags games
// @Param slug path string true "Slug of the Game"
// @Param Content-Digest header string true "Checksum in RFC 9530 format (algorithm=base64_checksum)"
// @Param request body binary true "Blob binary data"
// @Accept application/octet-stream
// @Produce json
// @Success 201 {object} map[string]string "Blob created"
// @Success 202 {object} map[string]string "Blob updated"
// @Failure 400 {object} map[string]string "Bad request: missing slug parameter, missing/invalid Content-Digest header"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/games/{slug}/blob [put]
func (ctl *EndpointController) PutBlob(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if len(slug) == 0 {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing mandatory parameter: %s", "slug"))
		return
	}

	contentType := ctx.GetHeader("Content-Type")
	if contentType != "" && strings.HasPrefix(contentType, "multipart/form-data") {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("multipart/form-data is not supported for blob uploads; use raw binary data"))
		return
	}

	digestHeader := ctx.GetHeader("Content-Digest")
	if digestHeader == "" {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing Content-Digest header"))
		return
	}

	algorithm, checksum, err := header.DecodeContentDigestHeader(digestHeader)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid Content-Digest header: %w", err))
		return
	}

	created, err := ctl.Service.Command.StoreNewAsset(
		appGameSrv.StoreNewAssetCommand{
			Slug:      slug,
			Role:      domainGame.GAME_ASSET_ROLE_BLOB,
			Checksum:  checksum,
			Algorithm: algorithm,
			Data:      ctx.Request.Body,
		},
	)
	if err != nil {
		handler.AbortWithError(ctx, err)
		return
	}

	if created {
		ctx.Status(http.StatusCreated)
		return
	}
	ctx.Status(http.StatusAccepted)
}

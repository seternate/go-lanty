package game

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appGameSrv "github.com/seternate/go-lanty/internal/application/game"
	domainGame "github.com/seternate/go-lanty/internal/domain/game"
	"github.com/seternate/go-lanty/internal/interface/http/adapter/header"
	errorx "github.com/seternate/go-lanty/internal/interface/http/error"
	errormodel "github.com/seternate/go-lanty/internal/interface/http/model/error"
)

var _ = errormodel.ErrorResponse{}

// @Summary Get a Games blob
// @Description Get the blob for a Game
// @Tags games
// @Param slug path string true "Slug"
// @Produce application/octet-stream, application/json
// @Success 200 {file} file "Blob binary data"
// @Header 200 {string} Content-Digest "Checksum (RFC 9530: algorithm=base64_checksum)"
// @Failure 404 {object} errormodel.ErrorResponse
// @Failure 500 {object} errormodel.ErrorResponse
// @Router /games/{slug}/blob [get]
func (ctl *EndpointController) GetBlob(ctx *gin.Context) {
	slug := ctx.Param("slug")

	assetContent, err := ctl.Service.Query.FetchBlob(slug)
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to fetch blob: %w", err))
		return
	}
	defer assetContent.Data.Close()

	digestHeader, err := header.EncodeContentDigestHeader(assetContent.Algorithm, assetContent.Checksum)
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to encode content-digest header: %w", err))
		return
	}

	ctx.DataFromReader(http.StatusOK, int64(assetContent.Size), assetContent.MimeType, assetContent.Data, map[string]string{
		"Content-Digest": digestHeader,
	})
}

// @Summary Update a Games blob
// @Description Updates a Games blob.
// @Tags games
// @Param slug path string true "Slug"
// @Param Content-Digest header string true "Checksum (RFC 9530: algorithm=base64_checksum)"
// @Param request body string true "Blob binary data"
// @Accept application/octet-stream
// @Produce json
// @Success 201 "Blob created"
// @Success 202 "Blob updated"
// @Failure 400 {object} errormodel.ErrorResponse
// @Failure 404 {object} errormodel.ErrorResponse
// @Failure 500 {object} errormodel.ErrorResponse
// @Router /games/{slug}/blob [put]
func (ctl *EndpointController) PutBlob(ctx *gin.Context) {
	slug := ctx.Param("slug")

	contentType := ctx.GetHeader("Content-Type")
	if contentType != "" && strings.HasPrefix(contentType, "multipart/form-data") {
		errorx.AbortWithError(ctx, errorx.ErrBadRequest("multipart/form-data is not supported for blob uploads; use raw binary data"))
		return
	}

	digestHeader := ctx.GetHeader("Content-Digest")
	if digestHeader == "" {
		errorx.AbortWithError(ctx, errorx.ErrBadRequest("missing Content-Digest header"))
		return
	}

	algorithm, checksum, err := header.DecodeContentDigestHeader(digestHeader)
	if err != nil {
		errorx.AbortWithError(ctx, errorx.ErrBadRequest("invalid Content-Digest header").WithCause(err))
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
		errorx.AbortWithError(ctx, fmt.Errorf("failed to store new asset: %w", err))
		return
	}

	if created {
		ctx.Status(http.StatusCreated)
		return
	}
	ctx.Status(http.StatusAccepted)
}

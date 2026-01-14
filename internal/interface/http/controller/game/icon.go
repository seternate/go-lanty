package game

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appGameSrv "github.com/seternate/go-lanty/internal/application/game"
	domainGame "github.com/seternate/go-lanty/internal/domain/game"
	"github.com/seternate/go-lanty/internal/interface/http/adapter/header"
	errorx "github.com/seternate/go-lanty/internal/interface/http/error"
	model "github.com/seternate/go-lanty/internal/interface/http/model"
)

var _ = model.ErrorResponse{}

// @Summary Get the icon of a Game
// @Description Get the icon of a Game
// @Tags games
// @Param slug path string true "Slug"
// @Produce image/*, application/json
// @Success 200 {file} file "Icon binary data"
// @Header 200 {string} Content-Digest "Checksum (RFC 9530: algorithm=base64_checksum)"
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /games/{slug}/icon [get]
func (ctl *EndpointController) GetIcon(ctx *gin.Context) {
	slug := ctx.Param("slug")

	assetContent, err := ctl.Service.Query.FetchIcon(slug)
	if err != nil {
		errorx.AbortWithError(ctx, fmt.Errorf("failed to fetch icon: %w", err))
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

// @Summary Update a Games icon
// @Description Updates a Games icon.
// @Description Accepts icon data in two formats:
// @Description 1) Raw binary data
// @Description 2) multipart/form-data with exactly one file. The Content-Digest header must match the checksum of the uploaded file data.
// @Tags games
// @Param slug path string true "Slug"
// @Param Content-Digest header string true "Checksum (RFC 9530: algorithm=base64_checksum)"
// @Param request body string false "Icon binary data"
// @Param file formData file false "Icon file (for multipart/form-data uploads). Exactly one file must be provided when using multipart format."
// @Accept image/*, multipart/form-data
// @Produce json
// @Success 201 "Icon created"
// @Success 202 "Icon updated"
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /games/{slug}/icon [put]
func (ctl *EndpointController) PutIcon(ctx *gin.Context) {
	slug := ctx.Param("slug")

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

	var fileData io.ReadCloser

	contentTypeHeader := ctx.ContentType()
	if strings.HasPrefix(contentTypeHeader, "multipart/form-data") {
		form, err := ctx.MultipartForm()
		if err != nil {
			errorx.AbortWithError(ctx, errorx.ErrBadRequest("failed to parse multipart form").WithCause(err))
			return
		}
		defer form.RemoveAll()

		totalFileCount := 0
		var fileHeader *multipart.FileHeader
		for _, fileList := range form.File {
			totalFileCount += len(fileList)
			if fileHeader == nil && len(fileList) > 0 {
				fileHeader = fileList[0]
			}
		}

		if totalFileCount == 0 {
			errorx.AbortWithError(ctx, errorx.ErrBadRequest("no file found in multipart form"))
			return
		}
		if totalFileCount > 1 {
			errorx.AbortWithError(ctx, errorx.ErrBadRequest("multiple files provided in multipart form, expected exactly one file, got %d", totalFileCount))
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			errorx.AbortWithError(ctx, fmt.Errorf("failed to open uploaded file: %w", err))
			return
		}
		fileData = file
	} else {
		fileData = ctx.Request.Body
	}

	created, err := ctl.Service.Command.StoreNewAsset(
		appGameSrv.StoreNewAssetCommand{
			Slug:      slug,
			Role:      domainGame.GAME_ASSET_ROLE_ICON,
			Checksum:  checksum,
			Algorithm: algorithm,
			Data:      fileData,
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

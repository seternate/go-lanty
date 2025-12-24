package game

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appGameSrv "github.com/seternate/go-lanty/pkg/application/game"
	domainGame "github.com/seternate/go-lanty/pkg/domain/game"
	"github.com/seternate/go-lanty/pkg/interface/http/adapter/handler"
	"github.com/seternate/go-lanty/pkg/interface/http/adapter/header"
)

// @Summary Get the icon of a Game
// @Description Get the icon of a Game with metadata headers
// @Tags games
// @Param slug path string true "Slug of the Game"
// @Produce image/*
// @Success 200 {file} binary "Icon binary data"
// @Header 200 {string} Content-Length "Size of the icon"
// @Header 200 {string} Content-Type "MIME type of the icon"
// @Header 200 {string} Content-Digest "Checksum in RFC 9530 format (algorithm=base64_checksum)"
// @Failure 400 {object} map[string]string "Bad request: missing slug parameter"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/games/{slug}/icon [get]
func (ctl *EndpointController) GetIcon(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if len(slug) == 0 {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing mandatory parameter %q", "slug"))
		return
	}

	assetContent, err := ctl.Service.Query.FetchIcon(slug)
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

// @Summary Upsert icon for a Game
// @Description Creates or updates the icon for a Game. Requires Content-Digest header with checksum. Accepts icon data in two formats: 1) Raw binary data, or 2) multipart/form-data with exactly one file. The Content-Digest header must match the checksum of the uploaded file data.
// @Tags games
// @Param slug path string true "Slug of the Game"
// @Param Content-Digest header string true "Checksum in RFC 9530 format (algorithm=base64_checksum)"
// @Param file formData file false "Icon file (for multipart/form-data uploads). Exactly one file must be provided when using multipart format."
// @Param request body binary false "Icon binary data (for raw image/* uploads)."
// @Accept image/*
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} map[string]string "Icon created"
// @Success 202 {object} map[string]string "Icon updated"
// @Failure 400 {object} map[string]string "Bad request: missing/invalid Content-Digest header, no file provided, multiple files provided, or invalid multipart form"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/games/{slug}/icon [put]
func (ctl *EndpointController) PutIcon(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if len(slug) == 0 {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing mandatory parameter: %s", "slug"))
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

	var fileData io.ReadCloser

	contentTypeHeader := ctx.ContentType()
	if strings.HasPrefix(contentTypeHeader, "multipart/form-data") {
		// Handle multipart form data
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse multipart form: %w", err))
			return
		}
		defer form.RemoveAll()

		// Count all files across all field names
		totalFileCount := 0
		var fileHeader *multipart.FileHeader
		for _, fileList := range form.File {
			totalFileCount += len(fileList)
			if fileHeader == nil && len(fileList) > 0 {
				fileHeader = fileList[0]
			}
		}

		if totalFileCount == 0 {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("no file provided in multipart form"))
			return
		}
		if totalFileCount > 1 {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("multiple files provided in multipart form, expected exactly one file, got %d", totalFileCount))
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to open uploaded file: %w", err))
			return
		}
		fileData = file
	} else {
		// Handle raw binary data
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
		handler.AbortWithError(ctx, err)
		return
	}

	if created {
		ctx.Status(http.StatusCreated)
		return
	}
	ctx.Status(http.StatusAccepted)
}

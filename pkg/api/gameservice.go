package api

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/seternate/go-lanty/pkg/game"
	"github.com/seternate/go-lanty/pkg/network"
)

type GameService struct {
	client *Client
}

func (service *GameService) GetGames() (games []string, err error) {
	path, err := service.client.router.Get("GetGames").URLPath()
	if err != nil {
		return
	}
	request, err := service.client.newRESTRequest(http.MethodGet, service.client.buildURL(*path), nil, nil)
	if err != nil {
		return
	}

	response, err := service.client.doREST(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &games)
	if err != nil {
		return
	}
	slices.Sort(games)
	return
}

func (service *GameService) GetGame(slug string) (game game.Game, err error) {
	path, err := service.client.router.Get("GetGame").URLPath("slug", slug)
	if err != nil {
		return
	}
	request, err := service.client.newRESTRequest(http.MethodGet, service.client.buildURL(*path), nil, nil)
	if err != nil {
		return
	}

	response, err := service.client.doREST(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	gamejson, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(gamejson, &game)
	return
}

func (service *GameService) GetIcon(game game.Game) (image image.Image, err error) {
	path, err := service.client.router.Get("GetGameDownloadIcon").URLPath("slug", game.Slug)
	if err != nil {
		return
	}

	download, err := network.NewDownload(service.client.buildURL(*path))
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	download.Download(context.TODO(), buf)
	if download.Err != nil {
		return
	}

	//TODO LOGIC TO USE RIGHT DECODER RATHER THAN JUST PNG

	image, err = png.Decode(buf)
	if err != nil {
		return
	}

	return image, nil
}

func (service *GameService) Download(ctx context.Context, game game.Game, directory string) (download *network.Download, err error) {
	path, err := service.client.router.Get("GetGameDownload").URLPath("slug", game.Slug)
	if err != nil {
		return
	}

	download, err = network.NewDownload(service.client.buildURL(*path))
	if err != nil {
		return
	}

	filepath := filepath.Join(directory, download.Filename())
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Create(filepath)
	if err != nil {
		return
	}

	download.StartDownload(ctx, file)
	go func() {
		<-download.Done
		file.Close()
	}()

	return
}

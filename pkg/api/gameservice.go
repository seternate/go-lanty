package api

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/seternate/go-lanty/pkg/download"
	"github.com/seternate/go-lanty/pkg/game"
)

type GameService struct {
	client *Client
}

func (service *GameService) GetList() (game.Games, error) {
	request, err := service.client.newRESTRequest(http.MethodGet, "/games", nil, nil)
	if err != nil {
		return nil, err
	}

	response, err := service.client.doREST(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	games := game.Games{}
	err = json.Unmarshal(bodyData, &games)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (service *GameService) GetIcon(game game.Game) (image.Image, error) {
	request, err := service.client.newRESTRequest(http.MethodHead, "games/"+game.Slug+"/download/icon", nil, nil)
	if err != nil {
		return nil, err
	}
	response, err := service.client.doREST(request)
	response.Body.Close()
	if err != nil {
		return nil, err
	}

	download, err := download.NewDownload(response)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	err = download.Start(buf)
	if err != nil {
		return nil, err
	}

	//TODO: LOGIC TO USE RIGHT DECODER RATHER THAN JUST PNG

	image, err := png.Decode(buf)
	if err != nil {
		return nil, err
	}

	return image, nil

}

func (service *GameService) GetFile(game game.Game, directory string) (*download.Download, error) {
	request, err := service.client.newRESTRequest(http.MethodHead, "games/"+game.Slug+"/download", nil, nil)
	if err != nil {
		return nil, err
	}
	response, err := service.client.doREST(request)
	response.Body.Close()
	if err != nil {
		return nil, err
	}

	download, err := download.NewDownload(response)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(directory, download.Filename)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	go func() {
		download.Start(out)
		out.Close()
	}()

	return download, nil
}

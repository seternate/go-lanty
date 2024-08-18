package api

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/seternate/go-lanty/pkg/file"
	"github.com/seternate/go-lanty/pkg/network"
)

type FileService struct {
	client *Client
}

func (service *FileService) UploadFile(file string) (response file.FileUploadResponse, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fstat, _ := f.Stat()
	filesize := strconv.FormatInt(fstat.Size(), 10)

	path, err := service.client.router.Get("PostFile").URLPath()
	if err != nil {
		return
	}
	request, err := service.client.newRESTRequest(http.MethodPost, service.client.buildURL(*path), nil, f)
	if err != nil {
		return
	}
	request.Header.Add("Content-Type", mime.TypeByExtension(filepath.Ext(file)))
	request.ContentLength = fstat.Size()
	request.Header.Add("Content-Length", filesize)
	request.Header.Add("Content-Disposition", "attachment; filename="+filepath.Base(file))

	r, err := service.client.doREST(request)
	if err != nil {
		return
	}
	defer r.Body.Close()

	responsejson, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(responsejson, &response)

	return
}

func (service *FileService) GetFile(ctx context.Context, url url.URL, directory string) (download *network.Download, err error) {
	download, err = network.NewDownload(url)
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

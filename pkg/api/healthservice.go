package api

import (
	"errors"
	"net/http"
)

type HealthService struct {
	client *Client
}

func (service *HealthService) Health() (err error) {
	path, err := service.client.router.Get("GetHealth").URLPath()
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

	if response.StatusCode != http.StatusOK {
		return errors.New("server is unhealthy")
	}

	return
}

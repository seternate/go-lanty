package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/seternate/go-lanty/pkg/user"
)

type UserService struct {
	client *Client
}

func (service *UserService) GetUsers() (users []string, err error) {
	path, err := service.client.router.Get("GetUsers").URLPath()
	if err != nil {
		return
	}
	request, err := service.client.newRESTRequest(http.MethodGet, service.client.BuildURL(*path), nil, nil)
	if err != nil {
		return
	}

	response, err := service.client.doREST(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	usernames, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(usernames, &users)
	return
}

func (service *UserService) GetUser(ip string) (user user.User, err error) {
	path, err := service.client.router.Get("GetUser").URLPath("ip", ip)
	if err != nil {
		return
	}
	request, err := service.client.newRESTRequest(http.MethodGet, service.client.BuildURL(*path), nil, nil)
	if err != nil {
		return
	}

	response, err := service.client.doREST(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	userjson, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(userjson, &user)
	return
}

func (service *UserService) CreateNewUser(user user.User) (u user.User, err error) {
	userjson, err := json.Marshal(user)
	if err != nil {
		return
	}

	path, err := service.client.router.Get("PostUser").URLPath()
	if err != nil {
		return
	}
	request, err := service.client.newRESTRequest(http.MethodPost, service.client.BuildURL(*path), nil, bytes.NewReader(userjson))
	if err != nil {
		return
	}

	response, err := service.client.doREST(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	userjson, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(userjson, &u)
	return
}

func (service *UserService) UpdateUser(user user.User) (u user.User, err error) {
	return service.CreateNewUser(user)
}

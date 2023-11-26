package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/seternate/go-lanty/pkg/user"
)

type UserService struct {
	client *Client
}

func (service *UserService) GetList() (user.Users, error) {
	request, err := service.client.newRESTRequest(http.MethodGet, "/users", nil, nil)
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

	users := user.Users{}
	err = json.Unmarshal(bodyData, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *UserService) CreateNewUser(u user.User) (user.User, error) {
	user := user.User{}

	userJson, err := json.Marshal(u)
	if err != nil {
		return user, err
	}

	request, err := service.client.newRESTRequest(http.MethodPost, "/users", nil, userJson)
	if err != nil {
		return user, err
	}

	response, err := service.client.doREST(request)
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	bodyData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(bodyData, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

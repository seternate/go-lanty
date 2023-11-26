package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/seternate/go-lanty/pkg/network"
	"github.com/seternate/go-lanty/pkg/user"
)

type UserHandler struct {
	handler *Handler
	Users   map[string]user.User
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, req *http.Request) {
	keys := make([]string, 0, len(h.Users))
	for k := range h.Users {
		keys = append(keys, k)
	}

	response, err := json.Marshal(keys)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to encode keys of Users map")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	user, err := h.getUser(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn().Err(err).Send()
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msgf("Failed to encode user '%s'", user.Name)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *UserHandler) PostUser(w http.ResponseWriter, req *http.Request) {
	user := &user.User{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn().Err(err).Msg("Can not decode user POST request")
		return
	}

	user.IP = strings.Split(req.RemoteAddr, ":")[0]
	if user.IP == "127.0.0.1" {
		user.IP = network.GetOutboundIP().String()
	}

	if _, found := h.Users[user.Name]; found == true {
		h.Users[user.Name] = *user
		log.Info().Msgf("User '%s' already exists. Updating the user.", user.Name)
		w.WriteHeader(http.StatusOK)
	} else {
		h.Users[user.Name] = *user
		log.Info().Msgf("Added '%s' to UserHandler", user.Name)
		w.WriteHeader(http.StatusCreated)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(user)
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func (h *UserHandler) getUser(req *http.Request) (user user.User, err error) {
	vars := mux.Vars(req)
	name := vars["name"]

	user, found := h.Users[name]
	if found == false {
		err = errors.New(fmt.Sprintf("No user '%s' found", name))
		return
	}

	return
}

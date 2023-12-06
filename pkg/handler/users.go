package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/seternate/go-lanty/pkg/network"
	"github.com/seternate/go-lanty/pkg/user"
)

type Userhandler struct {
	parent *Handler
	Users  map[string]user.User
}

func (handler *Userhandler) GetUsers(w http.ResponseWriter, req *http.Request) {
	names := make([]string, 0, len(handler.Users))
	for name := range handler.Users {
		names = append(names, name)
	}

	namesjson, err := json.Marshal(names)
	if err != nil {
		log.Error().Err(err).Msg("failed to encode user name list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(namesjson)
	log.Trace().RawJSON("names", namesjson).Msg("GET - /users")
}

func (handler *Userhandler) GetUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]
	user, found := handler.Users[name]
	if !found {
		log.Warn().Str("name", name).Msg("user not available")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userjson, err := json.Marshal(user)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("failed to encode user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userjson)
	log.Trace().RawJSON("user", userjson).Msg("GET - /users/:name")
}

func (handler *Userhandler) PostUser(w http.ResponseWriter, req *http.Request) {
	user := &user.User{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(user)
	if err != nil {
		log.Warn().Err(err).Msg("failed to decode user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Trace().Interface("user", user).Msg("POST - /users - Payload")

	user.IP = strings.Split(req.RemoteAddr, ":")[0]
	if user.IP == "127.0.0.1" {
		log.Trace().Msg("POST - /users - localhost request")
		ip, err := network.GetOutboundIP()
		if err != nil {
			log.Debug().Err(err).Msg("can not retrieve local IP")
		} else {
			user.IP = ip.String()
		}
	}

	userjson, err := json.Marshal(user)
	if err != nil {
		log.Error().Err(err).Str("name", user.Name).Msg("failed to encode user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, found := handler.Users[user.Name]
	if found {
		log.Debug().Str("name", user.Name).Msg("updating user")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Debug().Str("name", user.Name).Msg("added user")
		w.WriteHeader(http.StatusCreated)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userjson)
	log.Trace().RawJSON("user", userjson).Msg("POST - /users")
}

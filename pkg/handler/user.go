package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/seternate/go-lanty/pkg/network"
	"github.com/seternate/go-lanty/pkg/user"
)

var UserStaleDuration time.Duration = 15 * time.Second

type Userhandler struct {
	parent *Handler
	users  user.Users
	mutex  sync.RWMutex
}

func NewUserHandler(ctx context.Context, errgrp *errgroup.Group, parent *Handler) (handler *Userhandler) {
	handler = &Userhandler{
		parent: parent,
	}
	errgrp.Go(func() error {
		return handler.run(ctx)
	})
	return
}

func (handler *Userhandler) GetUsers(w http.ResponseWriter, req *http.Request) {
	handler.mutex.RLock()
	ipsjson, err := json.Marshal(handler.users.IPs())
	handler.mutex.RUnlock()
	if err != nil {
		log.Error().Err(err).Msg("failed to encode user ips")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ipsjson)
	log.Trace().RawJSON("ips", ipsjson).Msg("GET - /users")
}

func (handler *Userhandler) GetUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	ip := vars["ip"]
	handler.mutex.RLock()
	user, err := handler.users.Get(ip)
	handler.mutex.RUnlock()
	if err != nil {
		log.Warn().Str("ip", ip).Msg("user not available")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userjson, err := json.Marshal(user)
	if err != nil {
		log.Error().Err(err).Str("ip", ip).Msg("failed to encode user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userjson)
	log.Trace().RawJSON("user", userjson).Msg("GET - /users/:ip")
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

	user.Lastupdate = time.Now()
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
		log.Error().Err(err).Str("ip", user.IP).Msg("failed to encode user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.mutex.RLock()
	_, err = handler.users.Get(user.IP)
	handler.mutex.RUnlock()
	if err != nil {
		handler.mutex.Lock()
		err = handler.users.Add(*user)
		handler.mutex.Unlock()
		if err != nil {
			log.Error().Err(err).Str("ip", user.IP).Msg("error adding user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Debug().Str("ip", user.IP).Msg("added user")
			w.WriteHeader(http.StatusCreated)
		}
	} else {
		handler.mutex.Lock()
		err = handler.users.Update(*user)
		handler.mutex.Unlock()
		if err != nil {
			log.Error().Err(err).Str("ip", user.IP).Msg("error updating user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Debug().Str("ip", user.IP).Msg("updated user")
			w.WriteHeader(http.StatusOK)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userjson)
	log.Trace().RawJSON("user", userjson).Msg("POST - /users")
}

func (handler *Userhandler) run(ctx context.Context) error {
	ticker := time.NewTicker(UserStaleDuration)
	for {
		select {
		case <-ticker.C:
			handler.mutex.RLock()
			users := handler.users.Users()
			handler.mutex.RUnlock()
			for _, user := range users {
				if time.Since(user.Lastupdate) > UserStaleDuration {
					handler.mutex.Lock()
					handler.users.Remove(user)
					handler.mutex.Unlock()
					log.Info().Str("ip", user.IP).Msg("removed stale user")
				}
			}
		case <-ctx.Done():
			log.Debug().Err(ctx.Err()).Msg("stopped UserHandler loop")
			return ctx.Err()
		}
	}
}

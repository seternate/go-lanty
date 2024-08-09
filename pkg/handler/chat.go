package handler

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty/pkg/chat"
)

type ChatHandler struct {
	parent      *Handler
	clients     map[*websocket.Conn]bool
	upgrader    websocket.Upgrader
	broadcaster chan chat.Message
	mutex       sync.RWMutex
}

func NewChatHandler(parent *Handler) (handler *ChatHandler) {
	handler = &ChatHandler{
		parent:  parent,
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		broadcaster: make(chan chat.Message, 100),
	}
	go handler.run()
	return
}

func (handler *ChatHandler) Chat(w http.ResponseWriter, req *http.Request) {
	client, err := handler.upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade incoming websocket connection")
		return
	}
	defer client.Close()
	log.Debug().Str("remoteAddress", client.RemoteAddr().String()).Msg("successfully upgraded incoming websocket connection")

	handler.mutex.Lock()
	handler.clients[client] = true
	handler.mutex.Unlock()

	for {
		messageType, reader, err := client.NextReader()
		if err != nil {
			handler.mutex.Lock()
			delete(handler.clients, client)
			handler.mutex.Unlock()
			log.Error().Err(err).Str("remoteAddress", client.RemoteAddr().String()).Msg("permanent error reading from client websocket")
			client.Close()
			return
		}
		if messageType == websocket.BinaryMessage {
			log.Warn().Str("remoteAddress", client.RemoteAddr().String()).Msg("ignoring message from /chat websocket connection due to BinaryMessage")
			continue
		}

		message, err := chat.ReadMessage(reader)
		if err != nil {
			log.Warn().Err(err).Str("remoteAddress", client.RemoteAddr().String()).Msg("ignoring message due to error reading chat message")
			continue
		}
		if len(strings.TrimSpace(message.GetMessage())) == 0 {
			log.Warn().Str("remoteAddress", client.RemoteAddr().String()).Msg("ignoring message due to empty body")
			continue
		}
		if len(strings.TrimSpace(message.GetUser().Name)) == 0 || len(strings.TrimSpace(message.GetUser().IP)) == 0 {
			log.Warn().Str("remoteAddress", client.RemoteAddr().String()).Interface("user", message.GetUser()).Msg("ignoring message due to malformed user")
			continue
		}
		if message.GetTime().IsZero() {
			log.Warn().Str("remoteAddress", client.RemoteAddr().String()).Msg("added missing timestamp to received message")
			message.SetTime(time.Now())
		}

		log.Debug().Str("remoteAddress", client.RemoteAddr().String()).Msg("broadcasting newly received message")
		handler.broadcaster <- message
	}
}

func (handler *ChatHandler) run() {
	for message := range handler.broadcaster {
		handler.mutex.RLock()
		clients := handler.clients
		handler.mutex.RUnlock()
		for client := range clients {
			client.WriteJSON(message)
		}
	}
}

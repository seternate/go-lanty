package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/seternate/go-lanty/pkg/chat"
)

type ChatService struct {
	client *Client
	ws     *websocket.Conn

	Messages chan chat.Message
	Error    error
}

// If client.baseURL is changed while being connected a Reconnect for this service has to be made
func (service *ChatService) Connect() (resp *http.Response, err error) {
	defer func() { service.Error = err }()
	path, err := service.client.router.Get("Chat").URLPath()
	if err != nil {
		return
	}
	path.Scheme = "ws"
	url := service.client.buildURL(*path)
	service.ws, resp, err = websocket.DefaultDialer.Dial((&url).String(), nil)
	if err != nil {
		err = fmt.Errorf("error connecting to server chat websocket: %w", err)
		return
	}
	service.Error = nil
	go service.run()
	return
}

func (service *ChatService) Disconnect() (err error) {
	if service.ws == nil {
		return
	}
	err = service.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		service.Error = err
	}
	return
}

func (service *ChatService) Reconnect() (err error) {
	defer func() { service.Error = err }()
	err = service.Disconnect()
	_, err = service.Connect()
	return
}

func (service *ChatService) SendMessage(message chat.Message) (err error) {
	if service.ws == nil {
		return fmt.Errorf("can not send chat message: %w", service.Error)
	}
	message.SetTime(time.Now())
	err = service.ws.WriteJSON(message)
	if err != nil {
		service.Error = err
	}
	return
}

func (service *ChatService) run() {
	for service.Error == nil && service.ws != nil {
		messageType, reader, err := service.ws.NextReader()
		if err != nil {
			service.Error = fmt.Errorf("%w: permanent error reading from client websocket", err)
			service.Disconnect()
			return
		}
		if messageType == websocket.BinaryMessage {
			continue
		}

		message, err := chat.ReadMessage(reader)
		if err != nil {
			continue
		}
		if len(strings.TrimSpace(message.GetMessage())) == 0 {
			continue
		}
		if len(strings.TrimSpace(message.GetUser().Name)) == 0 || len(strings.TrimSpace(message.GetUser().IP)) == 0 {
			continue
		}
		if message.GetTime().IsZero() {
			continue
		}

		service.Messages <- message
	}
}

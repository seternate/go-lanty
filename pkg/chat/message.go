package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/seternate/go-lanty/pkg/user"
)

type Message interface {
	GetType() Type
	GetUser() user.User
	GetMessage() string
	GetTime() time.Time
	SetTime(time.Time)
}

type tmpMessage struct {
	Type Type      `json:"type"`
	User user.User `json:"user"`
	Time time.Time `json:"time"`
}

func ReadMessage(r io.Reader) (Message, error) {
	tmpMessageJson, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tmpMessage := &tmpMessage{}
	err = json.Unmarshal(tmpMessageJson, tmpMessage)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
		return nil, err
	}

	switch tmpMessage.Type {
	case TYPE_TEXT:
		textmessage := &TextMessage{}
		err = json.Unmarshal(tmpMessageJson, textmessage)
		if err != nil {
			return nil, err
		}
		return textmessage, err
	case TYPE_FILE:
		filemessage := &FileMessage{}
		err = json.Unmarshal(tmpMessageJson, filemessage)
		if err != nil {
			return nil, err
		}
		return filemessage, err
	}

	return nil, fmt.Errorf("undefined message type: %s", tmpMessage.Type.String())
}

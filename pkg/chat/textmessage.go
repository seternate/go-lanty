package chat

import (
	"time"

	"github.com/seternate/go-lanty/pkg/user"
)

type TextMessage struct {
	Type    Type            `json:"type"`
	User    user.User       `json:"user"`
	Message textMessageText `json:"message"`
	Time    time.Time       `json:"time"`
}

type textMessageText struct {
	Text string `json:"text"`
}

func NewTextMessage(user user.User, message string) *TextMessage {
	return &TextMessage{
		Type:    TYPE_TEXT,
		User:    user,
		Message: textMessageText{message},
		Time:    time.Now(),
	}
}

func (message *TextMessage) GetType() Type {
	return message.Type
}

func (message *TextMessage) GetUser() user.User {
	return message.User
}

func (message *TextMessage) GetMessage() string {
	return message.Message.Text
}

func (message *TextMessage) GetTime() time.Time {
	return message.Time
}

func (message *TextMessage) SetTime(t time.Time) {
	message.Time = t
}

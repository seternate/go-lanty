package chat

import (
	"path/filepath"
	"time"

	"github.com/seternate/go-lanty/pkg/file"
	"github.com/seternate/go-lanty/pkg/user"
)

type FileMessage struct {
	Type    Type                    `json:"type"`
	User    user.User               `json:"user"`
	Message file.FileUploadResponse `json:"message"`
	Time    time.Time               `json:"time"`
}

func NewFileMessage(user user.User, fileuploadresponse file.FileUploadResponse) *FileMessage {
	return &FileMessage{
		Type:    TYPE_FILE,
		User:    user,
		Message: fileuploadresponse,
		Time:    time.Now(),
	}
}

func (message *FileMessage) GetType() Type {
	return message.Type
}

func (message *FileMessage) GetUser() user.User {
	return message.User
}

func (message *FileMessage) GetMessage() string {
	return filepath.Base(message.Message.URL)
}

func (message *FileMessage) GetTime() time.Time {
	return message.Time
}

func (message *FileMessage) SetTime(t time.Time) {
	message.Time = t
}

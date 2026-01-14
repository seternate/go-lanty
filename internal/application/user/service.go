package user

import (
	"time"

	"github.com/seternate/go-lanty/internal/domain/user"
)

type CommandService interface {
	UpsertUser(cmd UpsertUserCommand) (*user.User, bool, error)
	DeleteUser(ipv4Address string) error
}

type UpsertUserCommand struct {
	IPv4Address string
	Username    string
}

type QueryService interface {
	GetUsers() (UsersView, error)
}

type UserView struct {
	IPv4Address string
	Username    string
	CreatedAt   time.Time
}

type UsersView []UserView

type Service struct {
	Query   QueryService
	Command CommandService
}

func NewService(query QueryService, command CommandService) *Service {
	return &Service{
		Query:   query,
		Command: command,
	}
}

package user

import (
	"errors"
	"fmt"

	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	"github.com/seternate/go-lanty/internal/domain/user"
)

var _ CommandService = (*commandServiceImpl)(nil)

type commandServiceImpl struct {
	repository user.UserRepository
}

func NewCommandService(repository user.UserRepository) CommandService {
	return &commandServiceImpl{
		repository: repository,
	}
}

func (service *commandServiceImpl) UpsertUser(cmd UpsertUserCommand) (*user.User, bool, error) {
	existingUser, err := service.repository.GetUser(cmd.IPv4Address)
	if err != nil {
		var notFound domainerr.NotFound
		if errors.As(err, &notFound) {
			newUser, err := user.NewUser(cmd.Username, cmd.IPv4Address)
			if err != nil {
				return nil, false, fmt.Errorf("failed to create new user ipv4Address=%s username=%s: %w", cmd.IPv4Address, cmd.Username, err)
			}

			err = service.repository.SaveUser(newUser)
			if err != nil {
				return nil, false, fmt.Errorf("failed to save new user to repository: %w", err)
			}

			return newUser, true, nil
		}
		return nil, false, fmt.Errorf("failed to get user ipv4Address=%s from repository: %w", cmd.IPv4Address, err)
	}

	err = existingUser.SetUsername(cmd.Username)
	if err != nil {
		return nil, false, fmt.Errorf("failed to update user ipv4Address=%s username: %w", cmd.IPv4Address, err)
	}

	err = service.repository.SaveUser(existingUser)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save updated user ipv4Address=%s to repository: %w", cmd.IPv4Address, err)
	}

	return existingUser, false, nil
}

func (service *commandServiceImpl) DeleteUser(ipv4Address string) error {
	deleteUser, err := service.repository.GetUser(ipv4Address)
	if err != nil {
		return fmt.Errorf("failed to get user from repository: %w", err)
	}

	err = service.repository.DeleteUser(deleteUser.IPv4Address.String())
	if err != nil {
		return fmt.Errorf("failed to delete user from repository: %w", err)
	}

	return nil
}

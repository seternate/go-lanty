package user

import (
	"fmt"

	"github.com/seternate/go-lanty/internal/domain/user"
)

type UserRow struct {
	IPv4Address string `db:"ipv4_address"`
	Username    string `db:"username"`
}

func (row UserRow) Assemble() (*user.User, error) {
	user, err := user.RehydrateUser(row.Username, row.IPv4Address)
	if err != nil {
		return nil, fmt.Errorf("failed to rehydrate user ipv4Address=%s username=%s (database corruption): %w", row.IPv4Address, row.Username, err)
	}
	return user, nil
}

func Disassemble(u *user.User) *UserRow {
	return &UserRow{
		IPv4Address: u.IPv4Address.String(),
		Username:    u.Username,
	}
}

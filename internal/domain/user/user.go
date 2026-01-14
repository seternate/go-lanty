package user

import (
	"net"
	"regexp"

	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

type User struct {
	IPv4Address net.IP
	Username    string
}

func NewUser(username string, ipv4Address string) (*User, error) {
	user, err := hydrateUser(username, ipv4Address)
	if err != nil {
		return nil, domainerr.InvariantViolationErr("user", ipv4Address).WithCause(err)
	}

	return user, nil
}

func RehydrateUser(username string, ipv4Address string) (*User, error) {
	user, err := hydrateUser(username, ipv4Address)
	if err != nil {
		return nil, domainerr.TrustedInvariantViolationErr("user", ipv4Address).WithCause(err)
	}

	return user, nil
}

func hydrateUser(username string, ipv4Address string) (*User, error) {
	validationErrors := domainerr.ValidationErrs()

	err := validationErrors.Wrap(validateUsername(username))
	if err != nil {
		return nil, err
	}

	ip, err := validateIPv4(ipv4Address)
	if err != nil {
		validationErrors.Wrap(err)
		return nil, validationErrors
	}

	if validationErrors.HasErrors() {
		return nil, validationErrors
	}

	return &User{
		IPv4Address: ip,
		Username:    username,
	}, nil
}

func (user *User) SetUsername(username string) error {
	err := validateUsername(username)
	if err != nil {
		return domainerr.InvariantViolationErr("user", user.IPv4Address.String()).WithCause(err)
	}

	user.Username = username
	return nil
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func validateUsername(username string) error {
	validationErrors := domainerr.ValidationErrs()

	if len(username) == 0 {
		validationErrors.Wrap(domainerr.ValidationErr("username", "can not be empty"))
	}

	if len(username) > 64 {
		validationErrors.Wrap(domainerr.ValidationErr("username", "must be at most 64 characters").WithGot(username))
	}

	if !usernameRegex.MatchString(username) {
		validationErrors.Wrap(domainerr.ValidationErr("username", "does not match expected format").WithExpected("pattern=%s", usernameRegex.String()).WithGot(username))
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func validateIPv4(ipv4Address string) (net.IP, error) {
	ip := net.ParseIP(ipv4Address)
	if ip == nil {
		return nil, domainerr.ValidationErr("IPv4 address", "must be a valid IP address").WithGot(ipv4Address)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil, domainerr.ValidationErr("IPv4 address", "must be a valid IPv4 address").WithGot(ip.String())
	}

	return ipv4, nil
}

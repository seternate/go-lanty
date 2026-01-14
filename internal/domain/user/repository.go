package user

type UserRepository interface {
	GetUser(ipv4Address string) (*User, error)
	SaveUser(*User) error
	DeleteUser(ipv4Address string) error
}

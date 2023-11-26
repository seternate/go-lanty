package user

import (
	"errors"
	"slices"
)

type Users struct {
	users map[string]User
	ips   []string
}

func (users *Users) HasUser(ip string) bool {
	return slices.Contains(users.ips, ip)
}

func (users *Users) Add(user User) (err error) {
	if len(users.ips) == 0 {
		users.users = make(map[string]User)
	} else if users.HasUser(user.IP) {
		return errors.New("user already in list")
	}
	users.ips = append(users.ips, user.IP)
	users.users[user.IP] = user
	return
}

func (users *Users) Equal(u Users) bool {
	equalIps := slices.Compare(users.IPs(), u.IPs())
	if equalIps != 0 {
		return false
	}
	for _, ip := range users.IPs() {
		left, _ := users.Get(ip)
		right, _ := u.Get(ip)
		if left != right {
			return false
		}
	}
	return true
}

func (users *Users) Get(ip string) (user User, err error) {
	user, found := users.users[ip]
	if !found {
		err = errors.New("no user with specified ip found")
	}
	return
}

func (users *Users) Remove(user User) (err error) {
	if users.HasUser(user.IP) {
		delete(users.users, user.IP)
		index := slices.Index(users.ips, user.IP)
		users.ips = slices.Delete(users.ips, index, index+1)
	} else {
		return errors.New("no user with specified ip found to remove")
	}

	return
}

func (users *Users) IPs() []string {
	return users.ips
}

func (users *Users) Users() (userlist []User) {
	for _, ip := range users.ips {
		user, _ := users.Get(ip)
		userlist = append(userlist, user)
	}
	return
}

func (users *Users) Update(user User) (err error) {
	if !users.HasUser(user.IP) {
		return errors.New("no such user to update")
	}
	users.users[user.IP] = user
	return
}

package user

import "time"

type User struct {
	Name       string    `json:"name"`
	IP         string    `json:"ip"`
	Lastupdate time.Time `json:"-"`
}

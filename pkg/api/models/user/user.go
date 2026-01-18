package usermodel

type User struct {
	IPv4Address string `json:"ipv4Address"`
	Username    string `json:"username"`
}

type UpsertUserRequest struct {
	Username string `json:"username"`
}

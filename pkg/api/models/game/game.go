package gamemodel

import (
	"time"
)

type Game struct {
	Slug        string       `json:"slug"`
	Name        string       `json:"name"`
	Executables []Executable `json:"executables"`
	CreatedAt   time.Time    `json:"createdat"`
}

type UpsertGameRequest struct {
	Name        string                        `json:"name"`
	Executables []UpsertGameExecutableRequest `json:"executables"`
}

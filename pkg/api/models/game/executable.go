package gamemodel

import (
	"time"
)

type Executable struct {
	Role              string   `json:"role"`
	Path              string   `json:"path"`
	RequiresAdmin     *bool    `json:"requiresadmin,omitempty"`
	Format            *string  `json:"format,omitempty"`
	ArgumentSeperator *string  `json:"argumentseperator,omitempty"`
	Args              []Arg    `json:"args"`
	CreatedAt         time.Time `json:"createdat"`
}

type UpsertGameExecutableRequest struct {
	Role              string                           `json:"role"`
	Path              string                           `json:"path"`
	RequiresAdmin     *bool                            `json:"requiresadmin"`
	Format            *string                          `json:"format"`
	ArgumentSeperator *string                          `json:"argumentseperator"`
	Args              []UpsertGameExecutableArgRequest `json:"args"`
}

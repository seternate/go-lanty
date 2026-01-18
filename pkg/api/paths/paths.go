package paths

import (
	"net/http"
	"strings"
)

const (
	// APIBasePath is the base path for all API endpoints
	APIBasePath = "/api/v1"

	// Games endpoints
	GamesPath    = "/games"
	GamePath     = "/games/:slug"
	GameIconPath = "/games/:slug/icon"
	GameBlobPath = "/games/:slug/blob"

	// Users endpoints
	UsersPath = "/users"
	UserPath  = "/users/:ipv4Address"
)

var (
	GetGames = Path{
		Method:   http.MethodGet,
		Template: "/games",
	}

	GetGame = Path{
		Method:   http.MethodGet,
		Template: "/games/:slug",
	}

	PutGame = Path{
		Method:   http.MethodPut,
		Template: "/games/:slug",
	}

	DeleteGame = Path{
		Method:   http.MethodDelete,
		Template: "/games/:slug",
	}

	GetGameIcon = Path{
		Method:   http.MethodGet,
		Template: "/games/:slug/icon",
	}

	PutGameIcon = Path{
		Method:   http.MethodPut,
		Template: "/games/:slug/icon",
	}

	GetGameBlob = Path{
		Method:   http.MethodGet,
		Template: "/games/:slug/blob",
	}

	PutGameBlob = Path{
		Method:   http.MethodPut,
		Template: "/games/:slug/blob",
	}

	GetUsers = Path{
		Method:   http.MethodGet,
		Template: "/users",
	}

	PutUser = Path{
		Method:   http.MethodPut,
		Template: "/users/:ipv4Address",
	}

	DeleteUser = Path{
		Method:   http.MethodDelete,
		Template: "/users/:ipv4Address",
	}
)

type Path struct {
	Template string
	Method   string
}

func (p Path) Build(params Params) (string, error) {
	if params == nil {
		return p.Template, nil
	}

	return params.Apply(p.Template)
}

type Params interface {
	Apply(path string) (string, error)
}

func replaceParams(path string, key string, value string) string {
	return strings.ReplaceAll(path, ":"+key, value)
}

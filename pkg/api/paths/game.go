package paths

import "errors"

type GameParams struct {
	Slug string
}

func (p GameParams) Apply(path string) (string, error) {
	if p.Slug == "" {
		return "", errors.New("slug is required")
	}
	return replaceParams(path, "slug", p.Slug), nil
}

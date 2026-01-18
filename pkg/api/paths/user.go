package paths

import "errors"

type UserParams struct {
	IPv4Address string
}

func (p UserParams) Apply(path string) (string, error) {
	if p.IPv4Address == "" {
		return "", errors.New("ipv4Address is required")
	}
	return replaceParams(path, "ipv4Address", p.IPv4Address), nil
}

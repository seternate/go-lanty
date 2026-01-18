package api

import (
	"net/http"

	"github.com/seternate/go-lanty/pkg/api/internal"
)

type APIError = internal.APIError

func CheckAPIResponse(resp *http.Response) error {
	return internal.CheckResponse(resp)
}

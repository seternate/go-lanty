package internal

import (
	"context"
	"io"
	"net/http"

	"github.com/seternate/go-lanty/pkg/api/paths"
)

type ClientInterface interface {
	BuildURL(paths.Path, paths.Params) (string, error)
	NewRESTRequestWithContext(ctx context.Context, method string, url string, headers map[string]string, queryParams map[string]string, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
	RESTRequest(ctx context.Context, method string, path paths.Path, pathParams paths.Params, headers map[string]string, queryParams map[string]string, body io.Reader) (*http.Response, error)
	HTTPClient() *http.Client
	Headers() map[string]string
}

func CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return ParseErrorResponse(resp)
}

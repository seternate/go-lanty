package api

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/seternate/go-lanty/pkg/router"
)

type Client struct {
	baseURL    *url.URL
	httpclient *http.Client
	router     *router.Router

	Game *GameService
	User *UserService
}

func NewClient(baseURL *url.URL, timeout time.Duration) (client *Client) {
	httpclient := &http.Client{
		Timeout: timeout,
	}

	router := router.NewRouter().
		WithRoutes(router.GameRoutes(nil)).
		WithRoutes(router.UserRoutes(nil))

	client = &Client{
		baseURL:    baseURL,
		httpclient: httpclient,
		router:     router,
	}

	client.Game = &GameService{client: client}
	client.User = &UserService{client: client}

	return
}

func (c *Client) BuildURL(url url.URL) url.URL {
	url.Scheme = "http"
	url.Host = c.baseURL.Host
	return url
}

func (c *Client) newRESTRequest(method string, url url.URL, params map[string]string, body io.Reader) (request *http.Request, err error) {
	request, err = http.NewRequest(method, url.String(), body)
	if err != nil {
		return
	}

	q := request.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	request.URL.RawQuery = q.Encode()

	return
}

func (c *Client) doREST(request *http.Request) (*http.Response, error) {
	return c.httpclient.Do(request)
}

package api

import (
	"net/http"
	"net/url"
	"time"

	"github.com/seternate/go-lanty/pkg/router"
)

type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
	router     *router.Router

	Game *GameService
}

func NewClient(baseURL *url.URL, apiKey, apiSecret string, timeout time.Duration) (*Client, error) {
	httpclient := &http.Client{
		Timeout: timeout,
	}
	router := router.NewRouter().
		WithRoutes(router.GameRoutes(nil)).
		WithRoutes(router.UserRoutes(nil))
	client := &Client{
		BaseURL:    baseURL,
		httpClient: httpclient,
		router:     router,
	}

	client.Game = &GameService{client: client}

	return client, nil
}

func (c *Client) newRESTRequest(method, path string, params map[string]string, body interface{}) (*http.Request, error) {
	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(urlPath)
	request, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// if body != nil {
	// 	request.Header.Set("Content-Type", "application/json")
	// }
	//request.Header.Set("Accept", "application/json")

	q := request.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	request.URL.RawQuery = q.Encode()

	return request, nil
}

func (c *Client) doREST(request *http.Request) (*http.Response, error) {
	return c.httpClient.Do(request)
}

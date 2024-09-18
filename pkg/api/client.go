package api

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/seternate/go-lanty/pkg/chat"
	"github.com/seternate/go-lanty/pkg/router"
)

type Client struct {
	baseURL    *url.URL
	httpclient *http.Client
	router     *router.Router

	Health *HealthService
	Game   *GameService
	User   *UserService
	Chat   *ChatService
	File   *FileService
}

func NewClient(baseURL string, timeout time.Duration) (client *Client, err error) {
	httpclient := &http.Client{
		Timeout: timeout,
	}

	router := router.NewRouter().
		WithRoutes(router.HealthRoutes(nil)).
		WithRoutes(router.GameRoutes(nil)).
		WithRoutes(router.UserRoutes(nil)).
		WithRoutes(router.ChatRoutes(nil)).
		WithRoutes(router.FileRoutes(nil))

	client = &Client{
		httpclient: httpclient,
		router:     router,
	}
	err = client.SetBaseURL(baseURL)
	client.Health = &HealthService{client: client}
	client.Game = &GameService{client: client}
	client.User = &UserService{client: client}
	client.Chat = &ChatService{
		client:   client,
		Messages: make(chan chat.Message, 100),
	}
	client.File = &FileService{client: client}

	return
}

func (c *Client) SetBaseURL(baseURL string) (err error) {
	u, err := url.Parse(baseURL)
	if len(u.Host) == 0 || err != nil {
		u, err = url.Parse("http://" + baseURL)
	}
	if err == nil {
		c.baseURL = u
	}
	return
}

func (c *Client) buildURL(url url.URL) url.URL {
	if len(url.Scheme) == 0 {
		url.Scheme = c.baseURL.Scheme
	}
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

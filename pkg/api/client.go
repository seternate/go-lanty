package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/seternate/go-lanty/pkg/api/games"
	"github.com/seternate/go-lanty/pkg/api/internal"
	"github.com/seternate/go-lanty/pkg/api/paths"
	"github.com/seternate/go-lanty/pkg/api/users"
)

var _ internal.ClientInterface = &Client{}

type Client struct {
	config      *Config
	httpClient  *http.Client
	gamesClient *games.Client
	usersClient *users.Client
}

func New(baseURL string, opts ...Option) (*Client, error) {
	config, err := NewConfig(baseURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	client := &Client{
		config:     config,
		httpClient: httpClient,
	}

	client.gamesClient = games.New(client, nil, nil)
	client.usersClient = users.New(client)

	return client, nil
}

func (c *Client) BuildURL(path paths.Path, params paths.Params) (string, error) {
	pathStr, err := path.Build(params)
	if err != nil {
		return "", err
	}

	pathStr = strings.TrimPrefix(pathStr, "/")

	joinedURL, err := url.JoinPath(c.config.BaseURL.String(), pathStr)
	if err != nil {
		return "", fmt.Errorf("failed to join URL paths: %w", err)
	}

	return joinedURL, nil
}

func (c *Client) NewRESTRequestWithContext(ctx context.Context, method string, url string, headers map[string]string, queryParams map[string]string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request with context: %w", err)
	}

	q := request.URL.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	request.URL.RawQuery = q.Encode()

	for key, value := range c.Headers() {
		request.Header.Set(key, value)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	return request, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) RESTRequest(ctx context.Context, method string, path paths.Path, pathParams paths.Params, headers map[string]string, queryParams map[string]string, body io.Reader) (*http.Response, error) {
	url, err := c.BuildURL(paths.GetGames, pathParams)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := c.NewRESTRequestWithContext(ctx, paths.GetGames.Method, url, headers, queryParams, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if err := CheckAPIResponse(resp); err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("API response error: %w", err)
	}

	return resp, nil
}

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) Headers() map[string]string {
	return c.config.Headers
}

func (c *Client) Games() *games.Client {
	return c.gamesClient
}

func (c *Client) Users() *users.Client {
	return c.usersClient
}

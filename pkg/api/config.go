package api

import (
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Headers    map[string]string
	Timeout    time.Duration
}

func NewConfig(baseURL string, opts ...Option) (*Config, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse baseURL: %w", err)
	}

	if parsedURL.Scheme == "" {
		return nil, errors.New("baseURL must contain a scheme (e.g., http:// or https://)")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, errors.New("baseURL scheme must be http or https")
	}
	if parsedURL.Host == "" {
		return nil, errors.New("baseURL must contain a host")
	}

	config := &Config{
		BaseURL: parsedURL,
		Timeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(config)
	}

	return config, nil
}

type Option func(*Config)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

func WithHeader(key, value string) Option {
	return func(c *Config) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		c.Headers[key] = value
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(c *Config) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		maps.Copy(c.Headers, headers)
	}
}

package client

import (
	"net/http"
	"os"
)

type Option func(*clientImpl)

func WithAuth(username, password string) Option {
	return func(c *clientImpl) {
		c.username = username
		c.password = password
	}
}

func WithEnvAuth() Option {
	return WithAuth(os.Getenv("BAR_USER"), os.Getenv("BAR_PASS"))
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *clientImpl) {
		c.client = httpClient
	}
}

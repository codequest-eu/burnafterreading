package client

import (
	"fmt"
	"io"
	"net/http"
)

type clientImpl struct {
	url, username, password string
}

// BasicHTTP is a client for the server using Basic HTTP Auth as it's
// authorization mechanism.
func BasicHTTP(url, username, password string) *clientImpl {
	return &clientImpl{url, username, password}
}

func (c *clientImpl) Get(key string) (io.ReadCloser, error) {
	res, err := c.dial("GET", key, nil)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (c *clientImpl) Put(key string, source io.Reader) error {
	_, err := c.dial("PUT", key, source)
	return err
}

func (c *clientImpl) dial(method, key string, data io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.urlFor(key), data)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.password)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid response code: %d", res.StatusCode)
	}
	return res, nil
}

func (c *clientImpl) urlFor(key string) string {
	return fmt.Sprintf("%s/?key=%s", c.url, key)
}

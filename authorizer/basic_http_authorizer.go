package authorizer

import "net/http"

type basicHTTPAuthorizer struct {
	username, password string
}

func BasicHTTPAuthorizer(username, password string) *basicHTTPAuthorizer {
	return &basicHTTPAuthorizer{username, password}
}

func (a *BasicHTTPAuthorizer) Authorize(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	return ok && username == a.username && password == a.password
}

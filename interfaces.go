package burnafterreading

import (
	"io"
	"net/http"
)

// Authorizer allows the user to exernally supply the HTTP(s) authorization
// logic.
type Authorizer interface {
	// Authorize answers whether a particular request is legit.
	Authorize(r *http.Request) bool
}

// Storage allows the user to externally supply the storage mechanism.s
type Storage interface {
	// Put provides a sink for the data. The implementation may choose to
	// not allow duplicate data to be saved for a given key but this is not
	// a strict requirement.
	Put(key string) (io.WriteCloser, error)

	// Get provides a source the data for a given key.
	Get(key string) (io.ReadCloser, error)
}

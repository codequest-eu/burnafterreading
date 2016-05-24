package lib

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

// ReadDeleter is an interface for getting and deleting data from a given
// source. It's implemented by both the Storage underlying the server and the
// Client which communicates with it.
type ReadDeleter interface {
	// Get provides a source the data for a given key.
	Get(key string) (io.ReadCloser, error)

	// Delete removes an entry for a given key.
	Delete(key string) error
}

// Storage allows the user to externally supply the storage mechanisms.
type Storage interface {
	ReadDeleter

	// Put provides a sink for the data. The implementation may choose to
	// not allow duplicate data to be saved for a given key but this is not
	// a strict requirement.
	Put(key string) (io.WriteCloser, error)
}

// Client represents the other side of an exchange.
type Client interface {
	ReadDeleter

	// Put takes the data from the reader and pipes it to the server,
	// returning the status of the operation in form of an error.
	Put(key string, source io.Reader) error
}

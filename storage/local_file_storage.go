package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	// ErrAlreadyExists is an error is returned to signify that an entry for
	// a given key already exists in the service.
	ErrAlreadyExists = errors.New("entry already exists")
)

func LocalFileStorage(basePath string) (*localFileStorage, error) {
	if fi, err := os.Stat(basePath); err != nil || !fi.IsDir() {
		return nil, errors.New("invalid path")
	}
	return &localFileStorage{basePath}, nil
}

type localFileStorage struct {
	basePath string
}

// Put provides a writer for saving the entry as a local file.
func (lfs *localFileStorage) Put(key string) (io.WriteCloser, error) {
	return os.Create(lfs.pathFor(key))
}

// Get provides a reader for reading the entry from a local file.
func (lfs *localFileStorage) Get(key string) (io.ReadCloser, error) {
	path := lfs.pathFor(key)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &readCloserWithCallback{
		file,
		func() error {
			return os.Remove(path)
		},
	}, nil
}

func (lfs *localFileStorage) pathFor(key string) string {
	return filepath.Join(lfs.basePath, key)
}

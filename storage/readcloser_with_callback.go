package storage

import "io"

type readCloserWithCallback struct {
	io.ReadCloser
	callback func() error
}

func (rc *readCloserWithCallback) Close() error {
	if err := rc.ReadCloser.Close(); err != nil {
		return err
	}
	return rc.callback()
}

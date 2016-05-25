package mocks

import "io"

// MockStorage is a mock implementation of the Storage interface.
type MockStorage struct {
	MockReadDeleter
}

// Put is a mock implementation of the Storage's Put method.
func (ms *MockStorage) Put(key string) (io.WriteCloser, error) {
	args := ms.Called(key)
	return args.Get(0).(io.WriteCloser), args.Error(0)
}

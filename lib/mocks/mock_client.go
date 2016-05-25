package mocks

import "io"

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	MockReadDeleter
}

// Put is a mock implementation of the Client's Put method.
func (mc *MockClient) Put(key string, source io.Reader) error {
	return mc.Called(key, source).Error(0)
}

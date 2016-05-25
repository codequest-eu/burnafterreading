package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"
)

// MockReadDeleter is a mock implementation of the ReadDeleter interface.
type MockReadDeleter struct {
	mock.Mock
}

// Get is a mock implementation of the ReadDeleter's Get method.
func (mrd *MockReadDeleter) Get(key string) (io.ReadCloser, error) {
	args := mrd.Called(key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// Delete is a mock implementation of the ReadDeleter's Delete method.
func (mrd *MockReadDeleter) Delete(key string) error {
	return mrd.Called(key).Error(0)
}

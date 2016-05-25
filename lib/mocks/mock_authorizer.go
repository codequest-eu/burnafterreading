package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockAuthorizer is a mock implementation of the Authorizer interface.
type MockAuthorizer struct {
	mock.Mock
}

// Authorize is a mock implementation of the Authorizer's Authorize method.
func (ma *MockAuthorizer) Authorize(r *http.Request) bool {
	return ma.Called(r).Bool(0)
}

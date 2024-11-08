// Code generated by mockery v2.46.3. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// MockResponse is an autogenerated mock type for the Response type
type MockResponse struct {
	mock.Mock
}

type MockResponse_Expecter struct {
	mock *mock.Mock
}

func (_m *MockResponse) EXPECT() *MockResponse_Expecter {
	return &MockResponse_Expecter{mock: &_m.Mock}
}

// NewMockResponse creates a new instance of MockResponse. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockResponse(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockResponse {
	mock := &MockResponse{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

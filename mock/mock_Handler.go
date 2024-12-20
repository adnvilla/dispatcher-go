// Code generated by mockery v2.46.3. DO NOT EDIT.

package mock

import (
	context "context"

	dispatcher "github.com/adnvilla/dispatcher-go"
	mock "github.com/stretchr/testify/mock"
)

// MockHandler is an autogenerated mock type for the Handler type
type MockHandler[TRequest dispatcher.Request, TResponse dispatcher.Response] struct {
	mock.Mock
}

type MockHandler_Expecter[TRequest dispatcher.Request, TResponse dispatcher.Response] struct {
	mock *mock.Mock
}

func (_m *MockHandler[TRequest, TResponse]) EXPECT() *MockHandler_Expecter[TRequest, TResponse] {
	return &MockHandler_Expecter[TRequest, TResponse]{mock: &_m.Mock}
}

// Handle provides a mock function with given fields: ctx, request
func (_m *MockHandler[TRequest, TResponse]) Handle(ctx context.Context, request TRequest) (TResponse, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Handle")
	}

	var r0 TResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, TRequest) (TResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, TRequest) TResponse); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(TResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, TRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHandler_Handle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Handle'
type MockHandler_Handle_Call[TRequest dispatcher.Request, TResponse dispatcher.Response] struct {
	*mock.Call
}

// Handle is a helper method to define mock.On call
//   - ctx context.Context
//   - request TRequest
func (_e *MockHandler_Expecter[TRequest, TResponse]) Handle(ctx interface{}, request interface{}) *MockHandler_Handle_Call[TRequest, TResponse] {
	return &MockHandler_Handle_Call[TRequest, TResponse]{Call: _e.mock.On("Handle", ctx, request)}
}

func (_c *MockHandler_Handle_Call[TRequest, TResponse]) Run(run func(ctx context.Context, request TRequest)) *MockHandler_Handle_Call[TRequest, TResponse] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(TRequest))
	})
	return _c
}

func (_c *MockHandler_Handle_Call[TRequest, TResponse]) Return(_a0 TResponse, _a1 error) *MockHandler_Handle_Call[TRequest, TResponse] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHandler_Handle_Call[TRequest, TResponse]) RunAndReturn(run func(context.Context, TRequest) (TResponse, error)) *MockHandler_Handle_Call[TRequest, TResponse] {
	_c.Call.Return(run)
	return _c
}

// NewMockHandler creates a new instance of MockHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHandler[TRequest dispatcher.Request, TResponse dispatcher.Response](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHandler[TRequest, TResponse] {
	mock := &MockHandler[TRequest, TResponse]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

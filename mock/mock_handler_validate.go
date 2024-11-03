package mock

import context "context"

// TODO: hack to include handler.(Validator[TRequest]) on tests
func (_m *MockHandler[TRequest, TResponse]) Validate(ctx context.Context, request TRequest) error {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, TRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

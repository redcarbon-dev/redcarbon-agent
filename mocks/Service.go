// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

type Service_Expecter struct {
	mock *mock.Mock
}

func (_m *Service) EXPECT() *Service_Expecter {
	return &Service_Expecter{mock: &_m.Mock}
}

// RunService provides a mock function with given fields: ctx
func (_m *Service) RunService(ctx context.Context) {
	_m.Called(ctx)
}

// Service_RunService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RunService'
type Service_RunService_Call struct {
	*mock.Call
}

// RunService is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Service_Expecter) RunService(ctx interface{}) *Service_RunService_Call {
	return &Service_RunService_Call{Call: _e.mock.On("RunService", ctx)}
}

func (_c *Service_RunService_Call) Run(run func(ctx context.Context)) *Service_RunService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Service_RunService_Call) Return() *Service_RunService_Call {
	_c.Call.Return()
	return _c
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Deleter is an autogenerated mock type for the Deleter type
type Deleter struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0
func (_m *Deleter) Delete(_a0 int64) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDeleter interface {
	mock.TestingT
	Cleanup(func())
}

// NewDeleter creates a new instance of Deleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDeleter(t mockConstructorTestingTNewDeleter) *Deleter {
	mock := &Deleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

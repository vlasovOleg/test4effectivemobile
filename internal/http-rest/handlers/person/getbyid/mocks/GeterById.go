// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	model "test4effectivemobile/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// GeterById is an autogenerated mock type for the GeterById type
type GeterById struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: _a0
func (_m *GeterById) GetByID(_a0 int64) (model.Person, error) {
	ret := _m.Called(_a0)

	var r0 model.Person
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (model.Person, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(int64) model.Person); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(model.Person)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGeterById interface {
	mock.TestingT
	Cleanup(func())
}

// NewGeterById creates a new instance of GeterById. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGeterById(t mockConstructorTestingTNewGeterById) *GeterById {
	mock := &GeterById{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

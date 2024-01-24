// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	model "test4effectivemobile/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// Saver is an autogenerated mock type for the Saver type
type Saver struct {
	mock.Mock
}

// Save provides a mock function with given fields: _a0
func (_m *Saver) Save(_a0 model.Person) (int64, error) {
	ret := _m.Called(_a0)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(model.Person) (int64, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(model.Person) int64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(model.Person) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSaver interface {
	mock.TestingT
	Cleanup(func())
}

// NewSaver creates a new instance of Saver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSaver(t mockConstructorTestingTNewSaver) *Saver {
	mock := &Saver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

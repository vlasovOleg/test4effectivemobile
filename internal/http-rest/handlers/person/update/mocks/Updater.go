// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	model "test4effectivemobile/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// Updater is an autogenerated mock type for the Updater type
type Updater struct {
	mock.Mock
}

// Update provides a mock function with given fields: _a0
func (_m *Updater) Update(_a0 model.Person) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.Person) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUpdater interface {
	mock.TestingT
	Cleanup(func())
}

// NewUpdater creates a new instance of Updater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUpdater(t mockConstructorTestingTNewUpdater) *Updater {
	mock := &Updater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

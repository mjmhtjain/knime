// Code generated by mockery v2.52.4. DO NOT EDIT.

package mocks

import (
	model "github.com/mjmhtjain/knime/src/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// INatsRepository is an autogenerated mock type for the INatsRepository type
type INatsRepository struct {
	mock.Mock
}

// PublishMessage provides a mock function with given fields: message
func (_m *INatsRepository) PublishMessage(message *model.OutboxMessageEntity) error {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for PublishMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.OutboxMessageEntity) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewINatsRepository creates a new instance of INatsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewINatsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *INatsRepository {
	mock := &INatsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

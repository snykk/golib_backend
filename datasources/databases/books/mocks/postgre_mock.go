// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	books "github.com/snykk/golib_backend/domains/books"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Repository) Delete(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: ctx
func (_m *Repository) GetAll(ctx context.Context) ([]books.Domain, error) {
	ret := _m.Called(ctx)

	var r0 []books.Domain
	if rf, ok := ret.Get(0).(func(context.Context) []books.Domain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]books.Domain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: ctx, id
func (_m *Repository) GetById(ctx context.Context, id int) (books.Domain, error) {
	ret := _m.Called(ctx, id)

	var r0 books.Domain
	if rf, ok := ret.Get(0).(func(context.Context, int) books.Domain); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(books.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, book
func (_m *Repository) Store(ctx context.Context, book *books.Domain) (books.Domain, error) {
	ret := _m.Called(ctx, book)

	var r0 books.Domain
	if rf, ok := ret.Get(0).(func(context.Context, *books.Domain) books.Domain); ok {
		r0 = rf(ctx, book)
	} else {
		r0 = ret.Get(0).(books.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *books.Domain) error); ok {
		r1 = rf(ctx, book)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, book
func (_m *Repository) Update(ctx context.Context, book *books.Domain) error {
	ret := _m.Called(ctx, book)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *books.Domain) error); ok {
		r0 = rf(ctx, book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

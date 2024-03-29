// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	reviews "github.com/snykk/golib_backend/domains/reviews"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, domain
func (_m *Repository) Delete(ctx context.Context, domain *reviews.Domain) (int, error) {
	ret := _m.Called(ctx, domain)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, *reviews.Domain) int); ok {
		r0 = rf(ctx, domain)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *reviews.Domain) error); ok {
		r1 = rf(ctx, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *Repository) GetAll(ctx context.Context) ([]reviews.Domain, error) {
	ret := _m.Called(ctx)

	var r0 []reviews.Domain
	if rf, ok := ret.Get(0).(func(context.Context) []reviews.Domain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]reviews.Domain)
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

// GetByBookId provides a mock function with given fields: ctx, bookId
func (_m *Repository) GetByBookId(ctx context.Context, bookId int) ([]reviews.Domain, error) {
	ret := _m.Called(ctx, bookId)

	var r0 []reviews.Domain
	if rf, ok := ret.Get(0).(func(context.Context, int) []reviews.Domain); ok {
		r0 = rf(ctx, bookId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]reviews.Domain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, bookId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: ctx, id
func (_m *Repository) GetById(ctx context.Context, id int) (reviews.Domain, error) {
	ret := _m.Called(ctx, id)

	var r0 reviews.Domain
	if rf, ok := ret.Get(0).(func(context.Context, int) reviews.Domain); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(reviews.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUserId provides a mock function with given fields: ctx, userId
func (_m *Repository) GetByUserId(ctx context.Context, userId int) ([]reviews.Domain, error) {
	ret := _m.Called(ctx, userId)

	var r0 []reviews.Domain
	if rf, ok := ret.Get(0).(func(context.Context, int) []reviews.Domain); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]reviews.Domain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserReview provides a mock function with given fields: ctx, bookId, userId
func (_m *Repository) GetUserReview(ctx context.Context, bookId int, userId int) (reviews.Domain, error) {
	ret := _m.Called(ctx, bookId, userId)

	var r0 reviews.Domain
	if rf, ok := ret.Get(0).(func(context.Context, int, int) reviews.Domain); ok {
		r0 = rf(ctx, bookId, userId)
	} else {
		r0 = ret.Get(0).(reviews.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, bookId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, domain
func (_m *Repository) Store(ctx context.Context, domain *reviews.Domain) (reviews.Domain, error) {
	ret := _m.Called(ctx, domain)

	var r0 reviews.Domain
	if rf, ok := ret.Get(0).(func(context.Context, *reviews.Domain) reviews.Domain); ok {
		r0 = rf(ctx, domain)
	} else {
		r0 = ret.Get(0).(reviews.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *reviews.Domain) error); ok {
		r1 = rf(ctx, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, domain
func (_m *Repository) Update(ctx context.Context, domain *reviews.Domain) error {
	ret := _m.Called(ctx, domain)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *reviews.Domain) error); ok {
		r0 = rf(ctx, domain)
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

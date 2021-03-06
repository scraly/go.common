// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import time "time"

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Add provides a mock function with given fields: key, value, expire
func (_m *Store) Add(key string, value interface{}, expire time.Duration) error {
	ret := _m.Called(key, value, expire)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, time.Duration) error); ok {
		r0 = rf(key, value, expire)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connect provides a mock function with given fields:
func (_m *Store) Connect() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: key
func (_m *Store) Delete(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Flush provides a mock function with given fields:
func (_m *Store) Flush() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: key, value
func (_m *Store) Get(key string, value interface{}) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *Store) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Replace provides a mock function with given fields: key, data, expire
func (_m *Store) Replace(key string, data interface{}, expire time.Duration) error {
	ret := _m.Called(key, data, expire)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, time.Duration) error); ok {
		r0 = rf(key, data, expire)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Set provides a mock function with given fields: key, value, expire
func (_m *Store) Set(key string, value interface{}, expire time.Duration) error {
	ret := _m.Called(key, value, expire)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, time.Duration) error); ok {
		r0 = rf(key, value, expire)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

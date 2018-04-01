package mocks

import (
	club "github.com/ElSiipo/Klubbkoll/club"
	mock "github.com/stretchr/testify/mock"
)

type ClubRepository struct {
	mock.Mock
}

// Delete mock
func (_m *ClubRepository) Delete(clubID string) (bool, error) {
	ret := _m.Called(clubID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(clubID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clubID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID mock
func (_m *ClubRepository) GetByID(clubID string) (*club.Club, error) {
	ret := _m.Called(clubID)

	var r0 *club.Club
	if rf, ok := ret.Get(0).(func(string) *club.Club); ok {
		r0 = rf(clubID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*club.Club)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clubID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store mock
func (_m *ClubRepository) Store(c *club.Club) (string, error) {
	ret := _m.Called(c)

	var r0 string
	if rf, ok := ret.Get(0).(func(*club.Club) string); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*club.Club) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update mock
func (_m *ClubRepository) Update(_c0 *club.Club) (*club.Club, error) {
	ret := _m.Called(_c0)

	var r0 *club.Club
	if rf, ok := ret.Get(0).(func(*club.Club) *club.Club); ok {
		r0 = rf(_c0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*club.Club)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*club.Club) error); ok {
		r1 = rf(_c0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

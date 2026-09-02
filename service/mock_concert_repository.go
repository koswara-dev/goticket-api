package service

import (
	"gotiket-api/model"

	"github.com/stretchr/testify/mock"
)

type MockConcertRepository struct {
	mock.Mock
}

func (m *MockConcertRepository) Create(concert *model.Concert) error {
	args := m.Called(concert)
	if args.Error(0) == nil {
		concert.ID = 1
	}
	return args.Error(0)
}

func (m *MockConcertRepository) FindByID(id uint) (model.Concert, error) {
	args := m.Called(id)
	return args.Get(0).(model.Concert), args.Error(1)
}

func (m *MockConcertRepository) FindByUserID(userID uint) (model.Concert, error) {
	args := m.Called(userID)
	return args.Get(0).(model.Concert), args.Error(1)
}

func (m *MockConcertRepository) FindAll() ([]model.Concert, error) {
	args := m.Called()
	return args.Get(0).([]model.Concert), args.Error(1)
}

func (m *MockConcertRepository) FindAllPagination(page int, limit int, search string, sort string) ([]model.Concert, int64, error) {
	args := m.Called(page, limit, search, sort)
	return args.Get(0).([]model.Concert), args.Get(1).(int64), args.Error(2)
}

func (m *MockConcertRepository) Update(concert *model.Concert) error {
	args := m.Called(concert)
	return args.Error(0)
}

func (m *MockConcertRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

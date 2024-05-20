package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockLocationGateway struct {
	mock.Mock
}

func (m *mockLocationGateway) GetLocationByZipCode(ctx context.Context, zipCode string) (string, error) {
	args := m.Called(zipCode)
	return args.Get(0).(string), args.Error(1)
}

type mockTemperatureGateway struct {
	mock.Mock
}

func (m *mockTemperatureGateway) GetTemperatureByLocation(ctx context.Context, location string) (float64, error) {
	args := m.Called(location)
	return args.Get(0).(float64), args.Error(1)
}

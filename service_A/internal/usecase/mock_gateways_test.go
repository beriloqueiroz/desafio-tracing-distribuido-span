package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockOrchestrationGateway struct {
	mock.Mock
}

func (m *mockOrchestrationGateway) GetTemperatureByZipCode(ctx context.Context, zipCode string) (*GetTemperByZipCodeUseCaseOutput, error) {
	args := m.Called(zipCode)
	return args.Get(0).(*GetTemperByZipCodeUseCaseOutput), args.Error(1)
}

package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TesteExecuteUseCase(t *testing.T) {
	mockOrchestrationGateway := new(mockOrchestrationGateway)
	mockOrchestrationGateway.On("GetTemperatureByZipCode", "12365478").Return(10.5, nil)

	usecase := GetTemperByZipCodeUseCase{
		OrchestrationGateway: mockOrchestrationGateway,
	}
	output, err := usecase.Execute(context.Background(), "60541646")
	assert.Nil(t, err)
	assert.InDelta(t, 10.5, output.TempC, 0.00001)
	assert.InDelta(t, 50.9, output.TempF, 0.00001)
	assert.InDelta(t, 283.5, output.TempK, 0.00001)
}

func TesteExecuteUseCaseWhenInvalidZipCode(t *testing.T) {
	mockOrchestrationGateway := new(mockOrchestrationGateway)

	usecase := GetTemperByZipCodeUseCase{
		OrchestrationGateway: mockOrchestrationGateway,
	}
	output, err := usecase.Execute(context.Background(), "60541646")
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid zipcode")
}

func TesteExecuteUseCaseWhenNotFoundZipCode(t *testing.T) {
	mockOrchestrationGateway := new(mockOrchestrationGateway)
	mockOrchestrationGateway.On("GetTemperatureByZipCode", "12365478").Return(nil, errors.New("not found"))

	usecase := GetTemperByZipCodeUseCase{
		OrchestrationGateway: mockOrchestrationGateway,
	}

	output, err := usecase.Execute(context.Background(), "60541646")
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "can not find zipcode")
}

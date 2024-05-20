package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TesteExecuteUseCase(t *testing.T) {
	mockLocationGateway := new(mockLocationGateway)
	mockTemperatureGateway := new(mockTemperatureGateway)
	mockLocationGateway.On("GetLocationByZipCode", "60541646").Return("12365478", nil)
	mockTemperatureGateway.On("GetTemperatureByLocation", "12365478").Return(10.5, nil)

	usecase := GetTemperByZipCodeUseCase{
		LocationGateway:    mockLocationGateway,
		TemperatureGateway: mockTemperatureGateway,
	}
	output, err := usecase.Execute(context.Background(), "60541646")
	assert.Nil(t, err)
	assert.InDelta(t, 10.5, output.TempC, 0.00001)
	assert.InDelta(t, 50.9, output.TempF, 0.00001)
	assert.InDelta(t, 283.5, output.TempK, 0.00001)
}

func TesteExecuteUseCaseWhenInvalidZipCode(t *testing.T) {
	mockLocationGateway := new(mockLocationGateway)
	mockTemperatureGateway := new(mockTemperatureGateway)

	usecase := GetTemperByZipCodeUseCase{
		LocationGateway:    mockLocationGateway,
		TemperatureGateway: mockTemperatureGateway,
	}
	output, err := usecase.Execute(context.Background(), "60541646")
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid zipcode")
}

func TesteExecuteUseCaseWhenNotFoundZipCode(t *testing.T) {
	mockLocationGateway := new(mockLocationGateway)
	mockTemperatureGateway := new(mockTemperatureGateway)
	mockLocationGateway.On("GetLocationByZipCode", "60541646").Return(nil, errors.New("not found"))

	usecase := GetTemperByZipCodeUseCase{
		LocationGateway:    mockLocationGateway,
		TemperatureGateway: mockTemperatureGateway,
	}
	output, err := usecase.Execute(context.Background(), "60541646")
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "can not find zipcode")
}

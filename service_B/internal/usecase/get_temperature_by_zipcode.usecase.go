package usecase

import (
	"context"
	"errors"

	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/entity"
)

type GetTemperByZipCodeUseCase struct {
	LocationGateway    LocationGateway
	TemperatureGateway TemperatureGateway
}

type GetTemperByZipCodeUseCaseOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewGetTemperByZipCodeUseCase(locationGateway LocationGateway, temperatureGateway TemperatureGateway) *GetTemperByZipCodeUseCase {
	return &GetTemperByZipCodeUseCase{
		LocationGateway:    locationGateway,
		TemperatureGateway: temperatureGateway,
	}
}

func (uc *GetTemperByZipCodeUseCase) Execute(ctx context.Context, zipCode string) (GetTemperByZipCodeUseCaseOutput, error) {
	output := GetTemperByZipCodeUseCaseOutput{}
	zipCodeObj, err := entity.NewZipCode(zipCode)
	if err != nil {
		return output, err
	}
	location, err := uc.LocationGateway.GetLocationByZipCode(ctx, zipCode)
	if err != nil {
		return output, errors.New("can not find zipcode")
	}
	temperature, err := uc.TemperatureGateway.GetTemperatureByLocation(ctx, location)
	if err != nil {
		return output, err
	}
	tempLocation, err := entity.NewTemperatureLocation(zipCodeObj, temperature, location)
	if err != nil {
		return output, err
	}
	output.City = tempLocation.City
	output.TempC = tempLocation.TempC
	output.TempF = tempLocation.TempF
	output.TempK = tempLocation.TempK
	return output, nil
}

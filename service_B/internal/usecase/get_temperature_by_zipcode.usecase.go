package usecase

import (
	"context"

	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/entity"
)

type GetTemperByZipCodeUseCase struct {
	TemperatureGateway TemperatureGateway
}

type GetTemperByZipCodeUseCaseOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewGetTemperByZipCodeUseCase(temperatureGateway TemperatureGateway) *GetTemperByZipCodeUseCase {
	return &GetTemperByZipCodeUseCase{
		TemperatureGateway: temperatureGateway,
	}
}

func (uc *GetTemperByZipCodeUseCase) Execute(ctx context.Context, zipCode string) (GetTemperByZipCodeUseCaseOutput, error) {
	output := GetTemperByZipCodeUseCaseOutput{}
	zipCodeObj, err := entity.NewZipCode(zipCode)
	if err != nil {
		return output, err
	}
	temperature, city, err := uc.TemperatureGateway.GetTemperatureByZipCode(ctx, zipCodeObj.Value)
	if err != nil {
		return output, err
	}
	tempLocation, err := entity.NewTemperatureLocation(zipCodeObj, *temperature, *city)
	if err != nil {
		return output, err
	}
	output.City = tempLocation.City
	output.TempC = tempLocation.TempC
	output.TempF = tempLocation.TempF
	output.TempK = tempLocation.TempK
	return output, nil
}

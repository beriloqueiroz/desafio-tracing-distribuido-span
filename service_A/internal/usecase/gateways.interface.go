package usecase

import "context"

type GetTemperByZipCodeUseCaseOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type OrchestrationGateway interface {
	GetTemperatureByZipCode(ctx context.Context, zipCode string) (*GetTemperByZipCodeUseCaseOutput, error)
}

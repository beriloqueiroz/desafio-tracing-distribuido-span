package usecase

import (
	"context"

	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/entity"
)

type GetTemperByZipCodeUseCase struct {
	OrchestrationGateway OrchestrationGateway
}

func NewGetTemperByZipCodeUseCase(OrchestrationGateway OrchestrationGateway) *GetTemperByZipCodeUseCase {
	return &GetTemperByZipCodeUseCase{
		OrchestrationGateway: OrchestrationGateway,
	}
}

func (uc *GetTemperByZipCodeUseCase) Execute(ctx context.Context, zipCode string) (*GetTemperByZipCodeUseCaseOutput, error) {
	zipCodeObj, err := entity.NewZipCode(zipCode)
	if err != nil {
		return nil, err
	}
	tempLocation, err := uc.OrchestrationGateway.GetTemperatureByZipCode(ctx, zipCodeObj.Value)
	if err != nil {
		return nil, err
	}
	return tempLocation, nil
}

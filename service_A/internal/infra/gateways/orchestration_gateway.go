package gateways

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/usecase"
)

type OrchestrationGatewayImpl struct {
	Ctx context.Context
	Url string
}

func (gt *OrchestrationGatewayImpl) GetTemperatureByZipCode(ctx context.Context, zipCode string) (*usecase.GetTemperByZipCodeUseCaseOutput, error) {
	resp, error := http.DefaultClient.Get(gt.Url + "?cep=" + zipCode)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c usecase.GetTemperByZipCodeUseCaseOutput
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	gt.Ctx.Done()
	return &c, nil
}

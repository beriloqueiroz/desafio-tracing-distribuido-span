package gateways

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/usecase"
)

type OrchestrationGatewayImpl struct {
	Ctx    context.Context
	Url    string
	Client http.Client
}

func (gt *OrchestrationGatewayImpl) GetTemperatureByZipCode(ctx context.Context, zipCode string) (*usecase.GetTemperByZipCodeUseCaseOutput, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", gt.Url+"?cep="+zipCode, nil)
	if err != nil {
		return nil, err
	}

	resp, err := gt.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var c usecase.GetTemperByZipCodeUseCaseOutput
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	gt.Ctx.Done()
	return &c, nil
}

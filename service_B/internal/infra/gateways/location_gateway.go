package gateways

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type GetLocationGatewayImpl struct {
	Ctx    context.Context
	Client http.Client
}

type viaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

func (gt *GetLocationGatewayImpl) GetLocationByZipCode(ctx context.Context, zipCode string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+zipCode+"/json/", nil)
	if err != nil {
		return "", err
	}
	resp, err := gt.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var c viaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return "", err
	}
	gt.Ctx.Done()
	return c.Localidade, nil
}

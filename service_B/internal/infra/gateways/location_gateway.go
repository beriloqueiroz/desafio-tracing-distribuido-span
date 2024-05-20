package gateways

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type GetLocationGatewayImpl struct {
	Ctx context.Context
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
	location, err := buscaCep(zipCode)
	if err != nil {
		return "", err
	}
	gt.Ctx.Done()
	return location.Localidade, nil
}

func buscaCep(cep string) (*viaCEP, error) {
	resp, error := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c viaCEP
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}

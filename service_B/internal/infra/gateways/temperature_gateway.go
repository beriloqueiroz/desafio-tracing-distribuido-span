package gateways

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type GetTemperatureGatewayImpl struct {
	Ctx     context.Context
	BaseUrl string
	Key     string
	Client  http.Client
}

func (gt *GetTemperatureGatewayImpl) GetTemperatureByZipCode(ctx context.Context, zipCode string) (*float64, *string, error) {
	location, err := gt.buscaCep(ctx, zipCode)
	if err != nil {
		return nil, nil, err
	}
	output, err := gt.buscaTemp(ctx, *location)
	if err != nil {
		return nil, nil, err
	}
	return &output.Current.TempC, location, nil
}

func (gt *GetTemperatureGatewayImpl) buscaCep(ctx context.Context, zipCode string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+zipCode+"/json/", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var c viaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	if c.HasError {
		return nil, errors.New("can not find zipcode")
	}
	return &c.Localidade, nil
}

func (gt *GetTemperatureGatewayImpl) buscaTemp(ctx context.Context, city string) (*temperatureInfo, error) {
	uri := gt.BaseUrl + "?q=" + url.QueryEscape(city) + "&lang=pt-br&key=" + gt.Key
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	resp, error := http.DefaultClient.Do(req)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var t temperatureInfo
	error = json.Unmarshal(body, &t)

	if error != nil {
		return nil, error
	}
	return &t, nil
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
	HasError    bool   `json:"erro"`
}

type temperatureInfo struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch float64 `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch float64 `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            float64 `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree float64 `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   float64 `json:"humidity"`
		Cloud      float64 `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

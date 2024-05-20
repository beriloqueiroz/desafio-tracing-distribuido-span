package main

import (
	"context"
	"fmt"

	config "github.com/beriloqueiroz/desafio-temperatura-por-cep/configs"
	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/gateways"
	routes "github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/web/routes/api"
	webserver "github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/web/server"
	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/usecase"
)

func main() {
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	server := webserver.NewWebServer(configs.WebServerPort)

	initCtx := context.Background()

	locationGateway := &gateways.GetLocationGatewayImpl{
		Ctx: initCtx,
	}
	temperatureGateway := &gateways.GetTemperatureGatewayImpl{
		Ctx:     initCtx,
		BaseUrl: configs.TempBaseUrl,
		Key:     configs.TempApiKey,
	}

	getTemperUseCase := usecase.NewGetTemperByZipCodeUseCase(
		locationGateway,
		temperatureGateway,
	)
	getTemperatureRoute := routes.NewGetTemperatureRouteApi(*getTemperUseCase)
	server.AddRoute("GET /", getTemperatureRoute.Handler)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	server.Start()
}

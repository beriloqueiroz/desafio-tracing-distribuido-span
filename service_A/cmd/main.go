package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	config "github.com/beriloqueiroz/desafio-temperatura-por-cep/configs"
	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/gateways"
	routes "github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/web/routes/api"
	webserver "github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/web/server"
	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/usecase"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	initCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// load environment configs
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	// start web server
	server := webserver.NewWebServer(configs.WebServerPort)
	orchestrationGateway := &gateways.OrchestrationGatewayImpl{
		Ctx: initCtx,
		Url: configs.ServiceBUrl,
	}
	getTemperUseCase := usecase.NewGetTemperByZipCodeUseCase(
		orchestrationGateway,
	)
	getTemperatureRoute := routes.NewGetTemperatureRouteApi(*getTemperUseCase)
	server.AddRoute("POST /", getTemperatureRoute.Handler)
	srvErr := make(chan error, 1)
	go func() {
		fmt.Println("Starting web server on port", configs.WebServerPort)
		srvErr <- server.Start()
	}()

	// Wait for interruption.
	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-initCtx.Done():
		log.Println("Shutting down due to other reason...")
	}
}

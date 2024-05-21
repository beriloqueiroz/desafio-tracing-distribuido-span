package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	config "github.com/beriloqueiroz/desafio-temperatura-por-cep/configs"
	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/gateways"
	routes "github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/web/routes/api"
	webserver "github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/web/server"
	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/usecase"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func main() {

	service_name := "service_A"

	// graceful exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	initCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// load environment configs
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	shutdown, err := initTraceProvider(service_name, configs.OtelExporterEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(initCtx); err != nil {
			log.Fatal("failed shutdown TraceProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("request service A")
	// start web server
	server := webserver.NewWebServer(configs.WebServerPort, tracer)
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
		fmt.Println("Starting web server "+service_name+" on port", configs.WebServerPort)
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

func initTraceProvider(serviceName string, collectorUrl string) (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(serviceName),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	// conn, err := grpc.DialContext(ctx, collectorUrl,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithBlock(),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create gRPC connection to collector: %W", err)
	// }

	// traceExport, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

	traceExport, err := otlptracehttp.New(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create export trace: %w", err)
	}

	// create span para envio em batch
	bsp := sdktrace.NewBatchSpanProcessor(traceExport)

	//create tracer provider with span bsp
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tracerProvider.Shutdown, nil
}

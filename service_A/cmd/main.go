package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	config "github.com/beriloqueiroz/desafio-temperatura-por-cep/configs"
	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/gateways"
	routes "github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/web/routes/api"
	webserver "github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/infra/web/server"
	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/usecase"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	serviceName := "service_A"

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

	shutdown, err := initTraceProvider(serviceName, configs.OtelExporterEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(initCtx); err != nil {
			log.Fatal("failed shutdown TraceProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("request service A")
	server := webserver.NewWebServer(configs.WebServerPort, tracer)

	//insert automatic transport, any request with this client will be track
	//there are a generic interceptor
	client := &http.Client{
		Transport: otelhttp.NewTransport(Interceptor{http.DefaultTransport}),
	}

	// include gateways e use cases
	orchestrationGateway := &gateways.OrchestrationGatewayImpl{
		Url:    configs.ServiceBUrl,
		Client: client,
	}
	getTemperUseCase := usecase.NewGetTemperByZipCodeUseCase(
		orchestrationGateway,
	)

	// add routes and run server
	getTemperatureRoute := routes.NewGetTemperatureRouteApi(*getTemperUseCase)
	getTemperatureRoute.TestDelay = time.Millisecond * time.Duration(configs.TestDelay)
	server.AddRoute("POST /get-cep", getTemperatureRoute.Handler)
	srvErr := make(chan error, 1)
	go func() {
		fmt.Println("Starting web server "+serviceName+" on port", configs.WebServerPort)
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

	conn, err := grpc.DialContext(ctx, collectorUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %W", err)
	}

	traceExport, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

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

type Interceptor struct {
	core http.RoundTripper
}

func (Interceptor) modifyRequest(r *http.Request) *http.Request {
	// otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(r.Header))
	fmt.Println("LOG - Host: " + r.URL.Host)
	return r
}

func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	// modify before the request is sent
	newReq := i.modifyRequest(r)

	// send the request using the DefaultTransport
	return i.core.RoundTrip(newReq)
}

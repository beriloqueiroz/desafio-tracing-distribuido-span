package routes

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type GetTemperatureRouteApi struct {
	getTemperatureByZipCodeUseCase usecase.GetTemperByZipCodeUseCase
	OtelTracer                     trace.Tracer
}

func NewGetTemperatureRouteApi(getTemperatureByZipCodeUseCase usecase.GetTemperByZipCodeUseCase, otelTracer trace.Tracer) *GetTemperatureRouteApi {
	return &GetTemperatureRouteApi{
		getTemperatureByZipCodeUseCase: getTemperatureByZipCodeUseCase,
		OtelTracer:                     otelTracer,
	}
}

func (cr *GetTemperatureRouteApi) Handler(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := cr.OtelTracer.Start(ctx, r.URL.Path)
	defer span.End()
	zipCode := r.URL.Query().Get("cep")

	output, err := cr.getTemperatureByZipCodeUseCase.Execute(ctx, zipCode)

	if err != nil {
		if err.Error() == "invalid zipcode" {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if err.Error() == "can not find zipcode" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

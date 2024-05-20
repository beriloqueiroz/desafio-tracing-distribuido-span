package routes

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/desafio-temperatura-por-cep/internal/usecase"
)

type GetTemperatureRouteApi struct {
	getTemperatureByZipCodeUseCase usecase.GetTemperByZipCodeUseCase
}

func NewGetTemperatureRouteApi(getTemperatureByZipCodeUseCase usecase.GetTemperByZipCodeUseCase) *GetTemperatureRouteApi {
	return &GetTemperatureRouteApi{getTemperatureByZipCodeUseCase}
}

func (cr *GetTemperatureRouteApi) Handler(w http.ResponseWriter, r *http.Request) {
	zipCode := r.URL.Query().Get("cep")

	output, err := cr.getTemperatureByZipCodeUseCase.Execute(r.Context(), zipCode)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err.Error() == "invalid zipcode" {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		if err.Error() == "can not find zipcode" {
			w.WriteHeader(http.StatusNotFound)
		}
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

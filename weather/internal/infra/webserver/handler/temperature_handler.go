package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/api"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/dto"
	"github.com/go-chi/chi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type ApiTemperatureResponse struct {
	cepApi          api.Cep
	weatherApi      api.Weather
	Tracer          trace.Tracer
	OtelRequestName string
}

func NewApiTemperatureResponse(cepApi api.Cep, weatherApi api.Weather, tracer trace.Tracer, otelRequestName string) *ApiTemperatureResponse {
	return &ApiTemperatureResponse{
		cepApi:          cepApi,
		weatherApi:      weatherApi,
		Tracer:          tracer,
		OtelRequestName: otelRequestName,
	}
}

func (at *ApiTemperatureResponse) GetTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	fmt.Println(cep)

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := at.Tracer.Start(ctx, "Chamada externa"+at.OtelRequestName)
	defer span.End()

	if len(cep) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cepResponse, err := at.cepApi.GetViaCepResponse(cep, ctx)
	if err != nil {
		switch err.Error() {
		case "error making request":
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(dto.Error{Message: err.Error()}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		case "invalid zipcode":
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err := json.NewEncoder(w).Encode(dto.Error{Message: err.Error()}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		case "cannot find zipcode":
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(dto.Error{Message: err.Error()}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	weatherResponse, err := at.weatherApi.GetWeatherApiResponse(cepResponse.Localidade, ctx)
	fmt.Println(weatherResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	temperatureResponse := &dto.TemperatureResponse{
		TemperatureC: weatherResponse.Current.TempC,
		TemperatureF: weatherResponse.Current.TempF,
		TemperatureK: CelsiusToKelvin(weatherResponse.Current.TempC),
		Location:     cepResponse.Localidade,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(temperatureResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

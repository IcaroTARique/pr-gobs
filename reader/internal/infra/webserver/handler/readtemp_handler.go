package handler

import (
	"encoding/json"
	"github.com/IcaroTARique/pr-gobs/internal/dto"
	"github.com/IcaroTARique/pr-gobs/internal/infra/weather_consumer"
	"github.com/go-chi/chi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type ApiTemperatureHandler struct {
	weather         weather_consumer.WeatherConsumer
	Tracer          trace.Tracer
	Zipkin          trace.Tracer
	OtelRequestName string
}

func NewApiTemperatureHandler(weather weather_consumer.WeatherConsumer, tracer, zipTracer trace.Tracer, otelRequestName string) *ApiTemperatureHandler {
	return &ApiTemperatureHandler{
		weather:         weather,
		Tracer:          tracer,
		Zipkin:          zipTracer,
		OtelRequestName: otelRequestName,
	}
}

func (th *ApiTemperatureHandler) NewApiTemperatureHandler(w http.ResponseWriter, r *http.Request) {

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	//Zipkin
	ctx, zipSpan := th.Zipkin.Start(ctx, "Chamada Zipkin: "+th.OtelRequestName)
	defer zipSpan.End()

	ctx, span := th.Tracer.Start(ctx, "Chamada externa "+th.OtelRequestName)
	defer span.End()

	cep := chi.URLParam(r, "cep")
	if len(cep) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := th.weather.GetTemperature(cep, ctx)
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

package handler

import (
	"encoding/json"
	"github.com/IcaroTARique/pr-gobs/internal/infra/weather_consumer"
	"github.com/go-chi/chi"
	"net/http"
)

type ApiTemperatureHandler struct {
	weather weather_consumer.WeatherConsumer
}

func NewApiTemperatureHandler(weather weather_consumer.WeatherConsumer) *ApiTemperatureHandler {
	return &ApiTemperatureHandler{
		weather: weather,
	}
}

func (th *ApiTemperatureHandler) NewApiTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if len(cep) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := th.weather.GetTemperature(cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

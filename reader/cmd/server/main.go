package main

import (
	"github.com/IcaroTARique/pr-gobs/configs"
	"github.com/IcaroTARique/pr-gobs/internal/infra/weather_consumer"
	"github.com/IcaroTARique/pr-gobs/internal/infra/webserver/handler"

	"github.com/go-chi/chi"
	"net/http"
)

func main() {

	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	weatherConsumer := weather_consumer.NewWeatherConsumer(conf.ReaderHost, conf.ReaderPort)
	weatherHandler := handler.NewApiTemperatureHandler(*weatherConsumer)

	r := chi.NewRouter()
	r.Get("/temperature/{cep}", weatherHandler.NewApiTemperatureHandler)

	if err := http.ListenAndServe(":"+conf.WebServerPort, r); err != nil {
		panic(err)
	}
}

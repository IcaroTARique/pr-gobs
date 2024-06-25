package weather_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IcaroTARique/pr-gobs/internal/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

type WeatherConsumer struct {
	Protocol        string
	Domain          string
	Port            string
	Endpoint        string
	OtelRequestName string
}

func NewWeatherConsumer(Domain, Port string) *WeatherConsumer {
	return &WeatherConsumer{
		Protocol: "http://",
		Domain:   Domain,
		Port:     Port,
		Endpoint: "/temperature/%s",
	}
}

func (w *WeatherConsumer) GetTemperature(cep string, ctx context.Context) (dto.TemperatureResponse, error) {

	Url := fmt.Sprintf(w.Protocol+w.Domain+":"+w.Port+w.Endpoint, cep)
	fmt.Println("URL :", Url)
	req, err := http.NewRequestWithContext(ctx, "GET", Url, nil)
	if err != nil {
		return dto.TemperatureResponse{}, fmt.Errorf("error making request")
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.TemperatureResponse{}, fmt.Errorf("error making request")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusBadRequest {
			return dto.TemperatureResponse{}, fmt.Errorf("error making request")
		}
		if res.StatusCode == http.StatusUnprocessableEntity {
			return dto.TemperatureResponse{}, fmt.Errorf("invalid zipcode")
		}
		if res.StatusCode == http.StatusNotFound {
			return dto.TemperatureResponse{}, fmt.Errorf("cannot find zipcode")
		}
		return dto.TemperatureResponse{}, fmt.Errorf("internal server error")
	}

	var temperatureResponse dto.TemperatureResponse
	err = json.NewDecoder(res.Body).Decode(&temperatureResponse)
	if err != nil {
		return dto.TemperatureResponse{}, fmt.Errorf("error parsing response")
	}
	fmt.Println(temperatureResponse)

	return temperatureResponse, nil
}

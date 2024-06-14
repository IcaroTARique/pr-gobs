package weather_consumer

import (
	"encoding/json"
	"fmt"
	"github.com/IcaroTARique/pr-gobs/internal/dto"
	"net/http"
)

type WeatherConsumer struct {
	Protocol string
	Domain   string
	Port     string
	Endpoint string
}

func NewWeatherConsumer(Domain, Port string) *WeatherConsumer {
	return &WeatherConsumer{
		Protocol: "http://",
		Domain:   Domain,
		Port:     Port,
		Endpoint: "/temperature/%s",
	}
}

func (w *WeatherConsumer) GetTemperature(cep string) (dto.TemperatureResponse, error) {
	Url := fmt.Sprintf(w.Protocol+w.Domain+":"+w.Port+w.Endpoint, cep)
	fmt.Println("URL :", Url)
	res, err := http.Get(Url)
	if err != nil {
		return dto.TemperatureResponse{}, fmt.Errorf("error making request")
	}

	var temperatureResponse dto.TemperatureResponse
	err = json.NewDecoder(res.Body).Decode(&temperatureResponse)
	if err != nil {
		return dto.TemperatureResponse{}, fmt.Errorf("error parsing response")
	}
	fmt.Println(temperatureResponse)

	return temperatureResponse, nil
}

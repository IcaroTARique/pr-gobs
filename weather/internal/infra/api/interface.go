package api

import "github.com/IcaroTARique/pr-locate-weather/internal/infra/dto"

type Weather interface {
	GetWeatherApiResponse(cityName string) (dto.WeatherApiResponse, error)
}

type Cep interface {
	GetViaCepResponse(cep string) (dto.ViaCepResponse, error)
}

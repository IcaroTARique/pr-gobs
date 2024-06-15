package api

import (
	"context"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/dto"
)

type Weather interface {
	GetWeatherApiResponse(cityName string, ctx context.Context) (dto.WeatherApiResponse, error)
}

type Cep interface {
	GetViaCepResponse(cep string, ctx context.Context) (dto.ViaCepResponse, error)
}

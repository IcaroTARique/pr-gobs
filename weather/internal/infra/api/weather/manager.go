package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io"
	"net/http"
	"strings"
	"unicode"
)

type ApiWeather struct {
	Url             string
	XApi            string
	OtelRequestName string
}

func NewApiWeather(otelRequestName string) *ApiWeather {
	return &ApiWeather{
		Url:             "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		XApi:            "e547194a521a49ddbcf220303241206",
		OtelRequestName: otelRequestName,
	}
}

func (aw *ApiWeather) GetWeatherApiResponse(cityName string, ctx context.Context) (dto.WeatherApiResponse, error) {

	treatedCityName := UnicodeFormatCityNameString(cityName)
	webUrlFormatTreatedCityName := WebUrlFormatCityNameString(treatedCityName)

	url := fmt.Sprintf(aw.Url, aw.XApi, webUrlFormatTreatedCityName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.WeatherApiResponse{}, fmt.Errorf("error making request")
	}
	defer res.Body.Close()

	printableBody, err := io.ReadAll(res.Body)

	var weatherApiResponse dto.WeatherApiResponse
	err = json.Unmarshal(printableBody, &weatherApiResponse)
	if err != nil {
		fmt.Println(err)
	}

	return weatherApiResponse, nil
}

func WebUrlFormatCityNameString(cityName string) string {
	return strings.ReplaceAll(cityName, " ", "%20")
}

func UnicodeFormatCityNameString(cityName string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, cityName)
	return result
}

package weather

import (
	"encoding/json"
	"fmt"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/dto"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io"
	"net/http"
	"strings"
	"unicode"
)

type ApiWeather struct {
	Url  string
	XApi string
}

func NewApiWeather() *ApiWeather {
	return &ApiWeather{
		Url:  "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		XApi: "e547194a521a49ddbcf220303241206",
	}
}

func (aw *ApiWeather) GetWeatherApiResponse(cityName string) (dto.WeatherApiResponse, error) {

	treatedCityName := UnicodeFormatCityNameString(cityName)
	webUrlFormatTreatedCityName := WebUrlFormatCityNameString(treatedCityName)

	url := fmt.Sprintf(aw.Url, aw.XApi, webUrlFormatTreatedCityName)
	fmt.Println(url)

	res, err := http.Get(url)
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

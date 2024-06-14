package weather

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWeatherApiResponse(t *testing.T) {
	api := NewApiWeather()
	weather, err := api.GetWeatherApiResponse("João Pessoa")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(weather)
	assert.NotNil(t, weather)
}

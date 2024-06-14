package weather_consumer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWeatherConsumer_GetTemperature(t *testing.T) {
	temp, err := NewWeatherConsumer().GetTemperature("58046320")
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, temp)
	assert.NotEmpty(t, temp.TemperatureC)
	assert.NotEmpty(t, temp.TemperatureF)
	assert.NotEmpty(t, temp.TemperatureK)
	assert.NotEmpty(t, temp.Location)
}

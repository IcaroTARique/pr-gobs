package dto

type TemperatureResponse struct {
	TemperatureC float64 `json:"temperature_c"`
	TemperatureF float64 `json:"temperature_f"`
	TemperatureK float64 `json:"temperature_k"`
	Location     string  `json:"location"`
}

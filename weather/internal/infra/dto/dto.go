package dto

type ViaCepResponse struct {
	Cep        string `json:"cep"`
	Localidade string `json:"localidade"`
}

type Error struct {
	Message string `json:"message"`
}

type ViaCepError struct {
	Message string `json:"erro"`
}

type WeatherApiResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type TemperatureResponse struct {
	TemperatureC float64 `json:"temperature_c"`
	TemperatureF float64 `json:"temperature_f"`
	TemperatureK float64 `json:"temperature_k"`
	Location     string  `json:"location"`
}

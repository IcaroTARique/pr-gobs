package entity

type Cep struct {
	Cep        string `json:"cep"`
	Localidade string `json:"localidade"`
}

type Error struct {
	Message string `json:"message"`
}

type Weather struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
		TempK float64 `json:"temp_k"`
	} `json:"current"`
}

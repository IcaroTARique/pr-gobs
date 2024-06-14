package cep

import (
	"encoding/json"
	"fmt"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/dto"
	"net/http"
)

type ApiCep struct {
	Url string
}

func NewApiCep() *ApiCep {
	return &ApiCep{
		Url: "https://viacep.com.br/ws/%s/json/",
	}
}

func (ac *ApiCep) GetViaCepResponse(cep string) (dto.ViaCepResponse, error) {

	res, err := http.Get(fmt.Sprintf(ac.Url, cep))
	if err != nil {
		return dto.ViaCepResponse{}, fmt.Errorf("error making request")
	}
	if res.StatusCode == http.StatusBadRequest {
		return dto.ViaCepResponse{}, fmt.Errorf("invalid zipcode")
	}
	defer res.Body.Close()

	var viaCepResponse dto.ViaCepResponse
	err = json.NewDecoder(res.Body).Decode(&viaCepResponse)
	if err != nil {
		return dto.ViaCepResponse{}, fmt.Errorf("error parsing response")
	}
	if viaCepResponse.Localidade == "" && viaCepResponse.Cep == "" {
		return dto.ViaCepResponse{}, fmt.Errorf("cannot find zipcode")
	}

	return viaCepResponse, nil
}

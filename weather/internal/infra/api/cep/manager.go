package cep

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type ApiCep struct {
	Url             string
	Tracer          trace.Tracer
	OtelRequestName string
}

func NewApiCep(tracer trace.Tracer, otelRequestName string) *ApiCep {
	return &ApiCep{
		Url:             "https://viacep.com.br/ws/%s/json/",
		Tracer:          tracer,
		OtelRequestName: otelRequestName,
	}
}

func (ac *ApiCep) GetViaCepResponse(cep string, ctx context.Context) (dto.ViaCepResponse, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(ac.Url, cep), nil)
	if err != nil {
		return dto.ViaCepResponse{}, fmt.Errorf("error making request")
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	res, err := http.DefaultClient.Do(req)
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

package cep

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCep(t *testing.T) {
	cepApi := NewApiCep()
	res, _ := cepApi.GetViaCepResponse("58046320")

	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Cep)
	assert.NotEmpty(t, res.Localidade)
}

func TestGetWrongCep(t *testing.T) {
	cepApi := NewApiCep()
	res, err := cepApi.GetViaCepResponse("5804632")

	assert.Equal(t, "", res.Cep)
	assert.Equal(t, "", res.Localidade)
	assert.Error(t, err, "cannot find zipcode")
}

func TestGetCepWithWrongFormat(t *testing.T) {
	cepApi := NewApiCep()
	res, err := cepApi.GetViaCepResponse("58046321")
	assert.Error(t, err, "invalid zipcode")
	assert.Equal(t, "", res.Cep)
	assert.Equal(t, "", res.Localidade)

}

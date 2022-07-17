package cep

import (
	"errors"

	"testing"

	"github.com/gonzalesMK/brcep/cep/basecep"
	"github.com/stretchr/testify/assert"
)

type MockAPI struct {
	shouldErr    error
	shouldReturn *basecep.BrAddress
}

func (a *MockAPI) Fetch(cep string) (*basecep.BrAddress, error) {
	if a.shouldErr != nil {
		return nil, a.shouldErr
	}
	return a.shouldReturn, nil
}

func TestHandleShouldReturnErrorIfFetchReturnsError(t *testing.T) {

	var cepHandler = &CepHandler{
		CepApis: map[string]basecep.API{
			"mock": &MockAPI{
				shouldErr: errors.New("unknown error"),
			},
		},
	}

	// Assertions
	_, err := cepHandler.GetCep("11111111")
	if assert.Error(t, err) {
		assert.Equal(t, "no API responded successfully", err.Error())
	}
}

func TestHandleShouldReturnErrorIfURLIsInvalid(t *testing.T) {

	var cepHandler = &CepHandler{
		CepApis: map[string]basecep.API{
			"mock": &MockAPI{
				shouldErr: nil,
			},
		},
	}

	_, err := cepHandler.GetCep("12345678")

	// Assertions
	if assert.Error(t, err) {
		assert.Equal(t, "no API responded successfully", err.Error())
	}

}

func TestHandleShouldReturnErrorIfWithoutCepAndFormat(t *testing.T) {

	var cepHandler = &CepHandler{
		CepApis: map[string]basecep.API{
			"mock": &MockAPI{
				shouldErr: nil,
			},
		},
	}

	_, err := cepHandler.GetCep("abc")

	// Assertions
	if assert.Error(t, err) {
		assert.Equal(t, "CEP not identified", err.Error())
	}
}

func TestHandleShouldSucceed(t *testing.T) {
	example := basecep.BrAddress{
		Cep:         "01001-000",
		Endereco:    "Praça da Sé",
		Complemento: "lado ímpar",
		Cidade:      "São Paulo",
		Uf:          "SP",
		Bairro:      "Sé",
		Ibge:        "3550308",
	}
	var cepHandler = &CepHandler{
		CepApis: map[string]basecep.API{
			"mock": &MockAPI{
				shouldErr:    nil,
				shouldReturn: &example,
			},
		},
		Cache: make(map[string]*basecep.BrAddress),
	}

	result, err := cepHandler.GetCep("01001-000")

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, example, *result)
		assert.NotNil(t, cepHandler.Cache)
	}

}

func TestHandleCachedShouldHitCacheAndSucceed(t *testing.T) {
	var cached = basecep.BrAddress{
		Cep:         "01001000",
		Endereco:    "Praça da Sé",
		Complemento: "lado ímpar",
		Cidade:      "São Paulo",
		Uf:          "SP",
		Bairro:      "Sé",
		Ibge:        "3550308",
	}

	var cepHandler = &CepHandler{
		CepApis: map[string]basecep.API{
			"mock": &MockAPI{
				shouldErr:    nil,
				shouldReturn: nil,
			},
		},
		Cache: make(map[string]*basecep.BrAddress),
	}
	cepHandler.Cache["00789000"] = &cached

	result, err := cepHandler.GetCep("00789-000")

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, cached, *result)
		assert.NotNil(t, cepHandler.Cache)

		cached, found := cepHandler.Cache["00789000"]
		assert.Equal(t, true, found)
		assert.NotNil(t, cached)
	}

}

package viacep

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gonzalesMK/brcep/cep/basecep"
)

const (
	// ID holds the identifier of this implementation
	ID                  = "viacep"
	defaultViaCepAPIURL = "http://viacep.com.br/"
)

// API holds the API implementation for viacep.com.br
type API struct {
	url    string
	client *http.Client
}

type viacepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Complemento string `json:"complemento"`
	Cidade      string `json:"localidade"`
	Estado      string `json:"uf"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	DDD         string `json:"ddd"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
}

// NewViaCepAPI creates and return a new API struct
func NewViaCepAPI(url string, client *http.Client) *API {
	if len(url) <= 0 {
		url = defaultViaCepAPIURL
	}
	if client == nil {
		client = http.DefaultClient
	}

	return &API{url, client}
}

// Fetch implements API.Fetch by fetching the viacep.com.br and
// returning a BrCepResult
func (api *API) Fetch(cep string) (*basecep.BrAddress, error) {

	resp, err := api.client.Get(fmt.Sprintf(api.url+"ws/%s/json/", cep))
	if err != nil {
		return nil, fmt.Errorf("ViaCepApi.Fetch %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("ViaCepApi.Fetch %v %d", "invalid status code", resp.StatusCode)
	}

	// Decode result
	var result viacepResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("ViaCepApi.Fetch %v", err)
	}

	return result.toCEP(), nil
}

func (r viacepResponse) toCEP() *basecep.BrAddress {
	var result = new(basecep.BrAddress)

	result.Cep = r.Cep
	result.Endereco = r.Logradouro
	result.Bairro = r.Bairro
	result.Complemento = r.Complemento
	result.Cidade = r.Cidade
	result.Uf = r.Estado
	result.Latitude = r.Latitude
	result.Longitude = r.Longitude
	result.DDD = r.Ibge
	result.Unidade = r.Unidade
	result.Ibge = r.Ibge

	return result
}

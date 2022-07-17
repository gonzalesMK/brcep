package cepaberto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/leogregianin/brcep/api"
)

const (
	// ID holds the identifier of this implementation
	ID = "cepaberto"

	defaultCepAbertoAPIURL = "http://www.cepaberto.com/"
)

// API holds the API implementation for cepaberto.com
type API struct {
	url    string
	token  string
	client *http.Client
}

type responsePayload struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     struct {
		Nome string `json:"nome"`
		DDD  int    `json:"ddd"`
		Ibge string `json:"ibge"`
	}
	Estado struct {
		Sigla string `json:"sigla"`
	}
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// NewCepAbertoAPI creates and returns new API struct
func NewCepAbertoAPI(url, token string, client *http.Client) *API {
	if len(url) <= 0 {
		url = defaultCepAbertoAPIURL
	}
	if client == nil {
		client = http.DefaultClient
	}

	return &API{url, token, client}
}

// Fetch implements API.Fetch by fetching the cepaberto.com and
// returning a BrCepResult
func (api *API) Fetch(cep string) (*api.BrCepResult, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(api.url+"api/v3/cep?cep=%s", url.QueryEscape(cep)), nil)
	if err != nil {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf(`Token token=%s`, api.token))

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v %d", "invalid status code", resp.StatusCode)
	}

	var result responsePayload
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v", err)
	}

	return result.toBrCepResult(), nil
}

func (r responsePayload) toBrCepResult() *api.BrCepResult {
	var result = new(api.BrCepResult)

	result.Cep = r.Cep
	result.Endereco = r.Logradouro
	result.Bairro = r.Bairro
	result.Cidade = r.Cidade.Nome
	result.Uf = r.Estado.Sigla
	result.Latitude = r.Latitude
	result.Longitude = r.Longitude
	result.DDD = strconv.Itoa(r.Cidade.DDD)
	result.Ibge = r.Cidade.Ibge

	return result
}

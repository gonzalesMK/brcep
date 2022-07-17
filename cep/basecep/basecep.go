package basecep

import (
	"regexp"
)

// BrAddress holds the standardized JSON result from the API
type BrAddress struct {
	Cep         string `json:"cep"`
	Endereco    string `json:"endereco"`
	Bairro      string `json:"bairro"`
	Complemento string `json:"complemento"`
	Cidade      string `json:"cidade"`
	Uf          string `json:"uf"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	DDD         string `json:"ddd"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
}

// API holds the interface that represents an API capable of fetching
// a CEP an return a standardized struct
type API interface {
	// Fetch should fetch the result from the
	// API and return as BrCepResult
	Fetch(cep string) (*BrAddress, error)
}

var CepSanitizer = regexp.MustCompile("[^0-9]+")

// Sanitize replaces not number with "" ..
func (r *BrAddress) Sanitize() {
	r.Cep = CepSanitizer.ReplaceAllString(r.Cep, "")
}

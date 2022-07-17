package cep

import (
	"errors"

	"github.com/gonzalesMK/brcep/cep/basecep"
)

type (
	// CepHandler provides a handler for a http.Server
	// used to interface the API implementations and
	// the server request
	CepHandler struct {
		CepApis map[string]basecep.API
		Cache   map[string]*basecep.BrAddress
	}

	channelResponse struct {
		ApiID  string
		Result *basecep.BrAddress
		Error  error
	}
)

func (h *CepHandler) fetchCep(cep_str string) (*basecep.BrAddress, error) {

	var (
		ch         = make(chan *channelResponse)
		apiResults *basecep.BrAddress
	)

	for id, currentAPI := range h.CepApis {
		go h.fetch(cep_str, id, currentAPI, ch)
	}

	for range h.CepApis {
		var current = <-ch
		if current.Error == nil {
			apiResults = current.Result
			break
		}
	}

	if apiResults == nil {
		return nil, errors.New("no API responded successfully")
	}

	apiResults.Sanitize()

	h.Cache[cep_str] = apiResults

	return apiResults, nil
}

func (h *CepHandler) GetCep(cep_str string) (*basecep.BrAddress, error) {

	cep_str = basecep.CepSanitizer.ReplaceAllString(cep_str, "")
	if len(cep_str) != 8 {
		return nil, errors.New("CEP not identified")
	}

	if cached, found := h.Cache[cep_str]; found {
		return cached, nil
	}

	return h.fetchCep(cep_str)
}

func (h *CepHandler) fetch(cep string, id string, api basecep.API, ch chan<- *channelResponse) {

	result, err := api.Fetch(cep)

	ch <- &channelResponse{
		ApiID:  id,
		Result: result,
		Error:  err,
	}
}

package correios

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CorreiosSuite struct{}

func TestNewCorreiosApiSetDefaultUrl(t *testing.T) {
	var correiosAPI = NewCorreiosAPI("", nil)
	assert.Equal(t, correiosAPI.url, "https://apps.correios.com.br/")
	assert.NotNil(t, correiosAPI.client)
}

func (s *CorreiosSuite) TestFetchShouldFailWhenInvalidStatusCode(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var correiosAPI = NewCorreiosAPI("", httpClient)
	_, err := correiosAPI.Fetch("78048000")

	assert.NotNil(t, err)
}

func (s *CorreiosSuite) TestFetchShouldFailWhenInvalidJSON(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var correiosAPI = NewCorreiosAPI("", httpClient)
	_, err := correiosAPI.Fetch("78048000")

	assert.NotNil(t, err)
}

func (s *CorreiosSuite) TestFetchShouldSucceedWhenCorrectRemoteResponse(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
	<ns2:consultaCEPResponse xmlns:ns2="http://cliente.bean.master.sigep.bsb.correios.com.br/">
		<return>
			<bairro>Alvorada</bairro>
			<cep>78048000</cep>
			<cidade>Cuiabá</cidade>
			<complemento2>- de 5686 a 6588 - lado par</complemento2>
			<end>Avenida Miguel Sutil</end>
			<uf>MT</uf>
		</return>
	</ns2:consultaCEPResponse>
</soap:Body></soap:Envelope>`))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var correiosAPI = NewCorreiosAPI("http://localhost:8080/", httpClient)
	result, err := correiosAPI.Fetch("78048000")

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Cep, "78048000")
	assert.Equal(t, result.Endereco, "Avenida Miguel Sutil")
	assert.Equal(t, result.Complemento, "- de 5686 a 6588 - lado par")
	assert.Equal(t, result.Bairro, "Alvorada")
	assert.Equal(t, result.Cidade, "CuiabÃ¡")
	assert.Equal(t, result.Uf, "MT")
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}

package correios

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/felipefrizzo/brazilian-zipcode-api/internal/address"
	"github.com/felipefrizzo/brazilian-zipcode-api/internal/helpers"
	"github.com/richardwilkes/toolbox/errs"
)

const (
	requestBody = `
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cli="http://cliente.bean.master.sigep.bsb.correios.com.br/">
			<soapenv:Header/>
			<soapenv:Body>
				<cli:consultaCEP>
					<cep>%s</cep>
				</cli:consultaCEP>
			</soapenv:Body>
		</soapenv:Envelope>
	`

	zipcodeNotFound string = "CEP NAO ENCONTRADO"
	invalidZipcode  string = "CEP INVÁLIDO"
)

// Correios represents the Correios service interface.
// api reference documentation: https://www.correios.com.br/atendimento/developers/arquivos/manual-para-integracao-correios-api
type Correios interface {
	GetAddressByZipcode(ctx context.Context, zipcode string) (*address.Address, error)
}

type correios struct {
	url        string
	httpClient *http.Client
}

type soapResponseError struct {
	FaultError string `xml:"Body>Fault>faultstring"`
}

type soapResponse struct {
	Address address.Address `xml:"Body>consultaCEPResponse>return"`
}

// New creates a new instance of the Correios service.
func New(url string) Correios {
	return &correios{
		url:        url,
		httpClient: &http.Client{},
	}
}

// GetAddressByZipcode retrieves the address information for a given zip code.
func (c *correios) GetAddressByZipcode(ctx context.Context, zipcode string) (*address.Address, error) {
	var payload []byte
	payload = fmt.Appendf(payload, requestBody, zipcode)
	resp, err := c.sendHTTPRequest(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		var soapError soapResponseError
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errs.NewWithCause("failed to read response body", err)
		}

		if err := xml.Unmarshal(helpers.ISO8859ToUTF8(body), &soapError); err != nil {
			return nil, errs.NewWithCause("failed to unmarshal SOAP error", err)
		}

		if soapError.FaultError == zipcodeNotFound {
			return nil, errs.NewWithCause("zipcode not found", errors.New(soapError.FaultError))
		}

		if soapError.FaultError == invalidZipcode {
			return nil, errs.NewWithCause("invalid zipcode", errors.New(soapError.FaultError))
		}

		return nil, errs.NewWithCause("unexpected status code", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.NewWithCause("failed to read response body", err)
	}

	var soap soapResponse
	if err := xml.Unmarshal(helpers.ISO8859ToUTF8(body), &soap); err != nil {
		return nil, errs.NewWithCause("failed to unmarshal SOAP error", err)
	}

	return &soap.Address, nil
}

func (c *correios) sendHTTPRequest(ctx context.Context, method, url string, payload io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return nil, errs.NewWithCause("failed to create a request", err)
	}

	req.Header.Add("Content-Type", "application/soap+xml; charset=iso-8859-1")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, errs.NewWithCause("failed to send a request", err)
	}

	return resp, nil
}

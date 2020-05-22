package helpers

import (
	"bytes"
	"fmt"
	"net/http"
)

// RequestXML abstraction function to make HTTP request in XML
func RequestXML(url string, body string, fields []interface{}) (*http.Response, error) {
	client := &http.Client{}
	payload := []byte(fmt.Sprintf(body, fields...))

	request, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/soap+xml; charset=iso-8859-1")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

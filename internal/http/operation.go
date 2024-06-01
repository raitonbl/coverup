package internal

import (
	"bytes"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
	"github.com/raitonbl/coverup/pkg"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
)

type Operation struct {
	Context        context.Context
	body           []byte
	response       *http.Response
	requestHeaders map[string]string
}

func (instance *Operation) withHttpRequest() error {
	instance.requestHeaders = nil
	return nil
}

func (instance *Operation) withRequestBody(body *godog.DocString) error {
	instance.body = []byte(body.Content)
	return nil
}

func (instance *Operation) withRequestHeaders(headers *godog.Table) error {
	if instance.requestHeaders == nil {
		instance.requestHeaders = make(map[string]string)
	}
	for _, row := range headers.Rows {
		instance.requestHeaders[row.Cells[0].Value] = row.Cells[1].Value
	}
	return nil
}

func (instance *Operation) withGetMethod(url string) error {
	//TODO LOG ignoring payload if set
	return instance.doSendHttpRequest("GET", url, nil)
}

func (instance *Operation) withPostMethod(url string) error {
	return instance.doSendHttpRequest("POST", url, instance.body)
}

func (instance *Operation) withPutMethod(url string) error {
	return instance.doSendHttpRequest("PUT", url, instance.body)
}

func (instance *Operation) withPatchMethod(url string) error {
	return instance.doSendHttpRequest("PATCH", url, instance.body)
}

func (instance *Operation) withDeleteMethod(url string) error {
	//TODO LOG ignoring payload if set
	return instance.doSendHttpRequest("DELETE", url, nil)
}

func (instance *Operation) doSendHttpRequest(method, url string, payload []byte) error {
	var body io.Reader = nil
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	serverURI := url
	if strings.HasPrefix(serverURI, "/") {
		serverURI = instance.Context.GetServerURL() + serverURI
	}
	req, err := http.NewRequest(method, serverURI, body)
	if err != nil {
		return err
	}
	if instance.requestHeaders != nil {
		for k, v := range instance.requestHeaders {
			req.Header.Set(k, v)
		}
	}
	httpClient := instance.Context.GetHttpClient()
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	instance.response = res
	return nil
}

func (instance *Operation) withStatusCode(statusCode int) error {
	assert.Equal(nil, statusCode, instance.response.StatusCode)
	return nil
}

func (instance *Operation) withHttpResponseHeader(headerName, headerValue string) error {
	assert.Equal(nil, headerValue, instance.response.Header.Get(headerName))
	return nil
}

func (instance *Operation) withHttpResponseHeaders(headers *godog.Table) error {
	for _, row := range headers.Rows {
		assert.Equal(nil, row.Cells[1].Value, instance.response.Header.Get(row.Cells[0].Value))
	}
	return nil
}

func (instance *Operation) withResponseBodyURI(uri string) error {
	var b []byte
	var err error
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		b, err = pkg.ReadFromURL(instance.Context.GetResourcesHttpClient(), uri)
	} else if strings.HasPrefix(uri, "file://") {
		b, err = pkg.ReadFromFile(instance.Context.GetWorkDirectory(), uri)
	} else {
		return fmt.Errorf("unsupported URI: %s", uri)
	}
	if err != nil {
		return err
	}
	return instance.withHttpResponseBodyEqualTo(string(b))
}

func (instance *Operation) withHttpResponseBody(body *godog.DocString) error {
	return instance.withHttpResponseBodyEqualTo(body.Content)
}

func (instance *Operation) withHttpResponseBodyEqualTo(body string) error {
	responseBody, err := io.ReadAll(instance.response.Body)
	if err != nil {
		return err
	}
	assert.JSONEq(nil, body, string(responseBody))
	return nil
}

package internal

import (
	"bytes"
	"errors"
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
	body           []byte
	response       *http.Response
	Context        context.Context
	requestHeaders map[string]string
}

// HTTP Request Methods

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

func (instance *Operation) sendHttpRequest(method, url string) error {
	return instance.doSendHttpRequest(method, url, instance.body)
}

func (instance *Operation) doSendHttpRequest(method, url string, payload []byte) error {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	serverURI := url
	if strings.HasPrefix(serverURI, "/") {
		serverURI = instance.Context.GetServerURL() + serverURI
	}
	req, err := http.NewRequest(method, serverURI, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	for k, v := range instance.requestHeaders {
		req.Header.Set(k, v)
	}
	httpClient := instance.Context.GetHttpClient()
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	instance.response = res
	return nil
}

// Simplified HTTP Methods

func (instance *Operation) withGetMethod(url string) error {
	return instance.sendHttpRequest("GET", url)
}

func (instance *Operation) withPostMethod(url string) error {
	return instance.sendHttpRequest("POST", url)
}

func (instance *Operation) withPutMethod(url string) error {
	return instance.sendHttpRequest("PUT", url)
}

func (instance *Operation) withPatchMethod(url string) error {
	return instance.sendHttpRequest("PATCH", url)
}

func (instance *Operation) withDeleteMethod(url string) error {
	return instance.sendHttpRequest("DELETE", url)
}

// HTTP Response Verification Methods

func (instance *Operation) withStatusCode(statusCode int) error {
	if instance.response.StatusCode != statusCode {
		return fmt.Errorf("expected status code %d, got %d", statusCode, instance.response.StatusCode)
	}
	return nil
}

func (instance *Operation) withHttpResponseHeader(headerName, headerValue string) error {
	if instance.response.Header.Get(headerName) != headerValue {
		return fmt.Errorf("expected header %s to be %s, got %s", headerName, headerValue, instance.response.Header.Get(headerName))
	}
	return nil
}

func (instance *Operation) withHttpResponseHeaders(headers *godog.Table) error {
	for _, row := range headers.Rows {
		if err := instance.withHttpResponseHeader(row.Cells[0].Value, row.Cells[1].Value); err != nil {
			return err
		}
	}
	return nil
}

func (instance *Operation) withResponseBodyURI(uri string) error {
	var b []byte
	var err error
	switch {
	case strings.HasPrefix(uri, "http://"), strings.HasPrefix(uri, "https://"):
		b, err = pkg.ReadFromURL(instance.Context.GetResourcesHttpClient(), uri)
	case strings.HasPrefix(uri, "file://"):
		b, err = pkg.ReadFromFile(instance.Context.GetWorkDirectory(), uri)
	default:
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
	responseBody, err := instance.getResponseBody()
	if err != nil {
		return err
	}
	assert.JSONEq(nil, body, string(responseBody))
	return nil
}

// Helper for Header Verification
func (instance *Operation) extractHeaderFromResponse(header string) func() (any, error) {
	return func() (any, error) {
		value := instance.response.Header.Get(header)
		if value == "" {
			return nil, errors.New("header[" + header + "] is undefined")
		}
		return value, nil
	}
}

// Generic response verification

func (instance *Operation) withResponseLinePartEqualTo(componentType, k, value string, extractValueFrom func() (any, error)) error {
	return isEqualTo(extractValueFrom, value, func(v1, v2 any) error {
		return fmt.Errorf("$%s.%s=%v isn't equal to %v", componentType, k, v1, v2)
	})
}

func (instance *Operation) withResponseLinePartStartsWith(componentType, k, value string, extractValueFrom func() (any, error)) error {
	return startsWith(extractValueFrom, value, func(v1, v2 string) error {
		return fmt.Errorf("$%s.%s=%v doesn't start with %v", componentType, k, v1, v2)
	})
}

func (instance *Operation) withResponseLinePartEndsWith(componentType, k, value string, extractValueFrom func() (any, error)) error {
	return endsWith(extractValueFrom, value, func(v1, v2 string) error {
		return fmt.Errorf("$%s.%s=%v doesn't end with %v", componentType, k, v1, v2)
	})
}

func (instance *Operation) withResponseLinePartGreaterOrEqualTo(componentType, k string, value float64, extractValueFrom func() (any, error)) error {
	return isGreaterOrEqualTo(extractValueFrom, value, func(v1 any, v2 float64) error {
		return fmt.Errorf("$%s.%s=%v isn't greater nor equal to %v", componentType, k, v1, v2)
	})
}

func (instance *Operation) withResponseLinePartLesserOrEqualTo(componentType, k string, value float64, extractValueFrom func() (any, error)) error {
	return isLesserOrEqualTo(extractValueFrom, value, func(v1 any, v2 float64) error {
		return fmt.Errorf("$%s.%s=%v isn't lesser nor equal to %v", componentType, k, v1, v2)
	})
}

func (instance *Operation) withResponseLinePartGreaterThan(componentType, k string, value float64, extractValueFrom func() (any, error)) error {
	return isGreaterThan(extractValueFrom, value, func(v1 any, v2 float64) error {
		return fmt.Errorf("$%s.%s=%v isn't greater than %v", componentType, k, v1, v2)
	})
}

func (instance *Operation) withResponseLinePartLesserThan(componentType, k string, value float64, extractValueFrom func() (any, error)) error {
	return isLesserThan(extractValueFrom, value, func(v1 any, v2 float64) error {
		return fmt.Errorf("$%s.%s=%v isn't lesser than %v", componentType, k, v1, v2)
	})
}

func (instance *Operation) withResponseLinePartMatchesPattern(componentType, k string, value string, extractValueFrom func() (any, error)) error {
	return matchesPattern(extractValueFrom, value, func(v1 string, v2 string) error {
		return fmt.Errorf("$%s.%s=%v doesn't match the pattern %v", componentType, k, v1, v2)
	})
}

// Generic Verification Method

func (instance *Operation) withHeaderEqualTo(header, value string) error {
	return instance.withResponseLinePartEqualTo("headers", header, value, instance.extractHeaderFromResponse(header))
}

func (instance *Operation) withHeaderStartsWith(header, value string) error {
	return instance.withResponseLinePartStartsWith("headers", header, value, instance.extractHeaderFromResponse(header))
}

func (instance *Operation) withHeaderEndsWith(header, value string) error {
	return instance.withResponseLinePartEndsWith("headers", header, value, instance.extractHeaderFromResponse(header))
}

func (instance *Operation) withHeaderPathMatches(header, value string) error {
	return instance.withResponseLinePartMatchesPattern("headers", header, value, instance.extractHeaderFromResponse(header))
}

// Simplified Body Path Verification Methods
func (instance *Operation) getResponseBody() ([]byte, error) {
	if instance.body == nil {
		binary, err := io.ReadAll(instance.response.Body)
		if err != nil {
			return nil, err
		}
		instance.body = binary
	}
	return instance.body, nil
}

func (instance *Operation) extractPathFromResponse(d string) func() (any, error) {
	return func() (any, error) {
		binary, err := instance.getResponseBody()
		if err != nil {
			return nil, err
		}
		if len(binary) == 0 {
			return nil, nil
		}
		contentType := instance.response.Header.Get("content-type")
		if contentType == "" {
			contentType = "application/json"
		}
		return extractValueFromResponsePath(contentType, d, binary)
	}
}

func (instance *Operation) withBodyPathEqualTo(path, value string) error {
	return instance.withResponseLinePartEqualTo("body", path, value, instance.extractPathFromResponse(path))
}

func (instance *Operation) withBodyPathStartsWith(path, value string) error {
	return instance.withResponseLinePartStartsWith("body", path, value, instance.extractPathFromResponse(path))
}

func (instance *Operation) withBodyPathEndsWith(path, value string) error {
	return instance.withResponseLinePartEndsWith("body", path, value, instance.extractPathFromResponse(path))
}

func (instance *Operation) withBodyPathMatches(path, pattern string) error {
	return instance.withResponseLinePartMatchesPattern("body", path, pattern, instance.extractPathFromResponse(path))
}

// Numeric Verification Methods

func (instance *Operation) withHeaderGreaterOrEqualTo(header string, value float64) error {
	return instance.withResponseLinePartGreaterOrEqualTo("headers", header, value, instance.extractHeaderFromResponse(header))
}

func (instance *Operation) withHeaderLesserOrEqualTo(header string, value float64) error {
	return instance.withResponseLinePartLesserOrEqualTo("headers", header, value, instance.extractHeaderFromResponse(header))
}

func (instance *Operation) withHeaderGreaterThan(header string, value float64) error {
	return instance.withResponseLinePartGreaterThan("headers", header, value, instance.extractHeaderFromResponse(header))
}

func (instance *Operation) withHeaderLesserThan(header string, value float64) error {
	return instance.withResponseLinePartLesserThan("headers", header, value, instance.extractHeaderFromResponse(header))
}

func (instance *Operation) withBodyPathGreaterOrEqualTo(path string, value float64) error {
	return instance.withResponseLinePartGreaterOrEqualTo("body", path, value, instance.extractPathFromResponse(path))
}

func (instance *Operation) withBodyPathLesserOrEqualTo(path string, value float64) error {
	return instance.withResponseLinePartLesserOrEqualTo("body", path, value, instance.extractPathFromResponse(path))
}

func (instance *Operation) withBodyPathGreaterThan(path string, value float64) error {
	return instance.withResponseLinePartGreaterThan("body", path, value, instance.extractPathFromResponse(path))
}

func (instance *Operation) withBodyPathLesserThan(path string, value float64) error {
	return instance.withResponseLinePartLesserThan("body", path, value, instance.extractPathFromResponse(path))
}

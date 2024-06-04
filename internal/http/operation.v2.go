package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
	"github.com/raitonbl/coverup/pkg"
	"github.com/thoas/go-funk"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"net/http"
	"strings"
)

const defaultAlias = ""
const componentType = "HttpRequest"

var schemaCache = make(map[string]gojsonschema.JSONLoader)

func CreateHttpRequest(instance *context.Builder) func() error {
	f := CreateHttpRequestWithAlias(instance)
	return func() error {
		return f(defaultAlias)
	}
}

func CreateHttpRequestWithAlias(instance *context.Builder) func(string) error {
	return func(alias string) error {
		return instance.WithComponent(componentType, &Request{
			headers: make(map[string]string),
		}, alias)
	}
}

func getHttpRequest(instance *context.Builder, alias string) (*Request, error) {
	valueOf := instance.GetComponent(componentType, alias)
	if r, isHttpRequest := valueOf.(*Request); isHttpRequest {
		return r, nil
	} else {
		return nil, fmt.Errorf("please define %s.%s before mention", componentType, componentType)
	}
}
func SetRequestHeaders(instance *context.Builder) func(*godog.Table) error {
	return func(table *godog.Table) error {
		var req, err = getHttpRequest(instance, defaultAlias)
		if err != nil {
			return err
		}
		if req.headers == nil {
			req.headers = make(map[string]string)
		}
		for _, row := range table.Rows {
			valueOf, prob := instance.GetValue(row.Cells[1].Value)
			if prob != nil {
				return prob
			}
			req.headers[row.Cells[0].Value] = fmt.Sprintf("%v", valueOf)
		}
		return nil
	}
}

func SetRequestOperation(instance *context.Builder, method string) func(string) error {
	return func(v string) error {
		var req, err = getHttpRequest(instance, defaultAlias)
		if err != nil {
			return err
		}
		valueOf, prob := instance.GetValue(v)
		if prob != nil {
			return prob
		}
		url := fmt.Sprintf("%v", valueOf)
		if strings.HasPrefix(url, "/") {
			if req.serverURL == "" {
				req.serverURL = instance.GetServerURL()
			}
			req.uri = url
		} else {
			req.uri = ""
			req.serverURL = url
		}
		req.method = method
		return nil
	}
}

func SetRequestBody(instance *context.Builder) func(string) error {
	return func(v string) error {
		return setHttpRequestBody(instance, v)
	}
}

func SetRequestBodyFromURI(instance *context.Builder, uriSchema string) func(string) error {
	return func(uri string) error {
		b, err := fetchFromURI(instance, uriSchema, uri)
		if err != nil {
			return err
		}
		return setHttpRequestBody(instance, string(b))
	}
}

func setHttpRequestBody(instance *context.Builder, v string) error {
	var req, err = getHttpRequest(instance, defaultAlias)
	if err != nil {
		return err
	}
	req.body = []byte(v)
	return nil
}

func SubmitsHttpRequest(instance *context.Builder) func() error {
	return func() error {
		return submitsHttpRequest(instance, "")
	}
}

func SubmitsHttpRequestWhenAlias(instance *context.Builder) func(string) error {
	return func(alias string) error {
		return submitsHttpRequest(instance, alias)
	}
}

func submitsHttpRequest(instance *context.Builder, alias string) error {
	req, err := getHttpRequest(instance, alias)
	if err != nil {
		return err
	}
	if req.response != nil {
		return nil
	}
	return doSubmitHttpRequest(instance, req)
}

func doSubmitHttpRequest(instance *context.Builder, src *Request) error {
	var body io.Reader
	if src.body != nil {
		body = bytes.NewReader(src.body)
	}
	serverURI := src.serverURL
	if serverURI == "" {
		serverURI = instance.GetServerURL()
	}
	if src.uri != "" {
		serverURI += src.uri
	}
	req, err := http.NewRequest(src.method, serverURI, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	for k, v := range src.headers {
		req.Header.Set(k, v)
	}
	httpClient := instance.GetHttpClient()
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	binary, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	headers := make(map[string]string)
	for k, v := range res.Header {
		req.Header.Set(k, strings.Join(v, ","))
	}
	src.response = &Response{
		body:       binary,
		headers:    headers,
		statusCode: res.StatusCode,
		pathCache:  make(map[string]any),
	}
	return nil
}

func AssertHttpResponseStatusCode(instance *context.Builder) func(statusCode int) error {
	return func(statusCode int) error {
		return assertHttpResponseStatusCodeWhenAlias(instance, statusCode, "")
	}
}

func AssertHttpResponseStatusCodeWhenAlias(instance *context.Builder) func(string, int) error {
	return func(alias string, statusCode int) error {
		return assertHttpResponseStatusCodeWhenAlias(instance, statusCode, alias)
	}
}

func AssertHttpResponseBodySchemaOnURIWhenAlias(instance *context.Builder, uriSchema string) func(string, string) error {
	return func(alias, s string) error {
		return assertHttpResponseBodySchemaOnURIWhenAlias(instance, uriSchema, s, alias)
	}
}

func AssertHttpResponseBodySchemaOnURI(instance *context.Builder, uriSchema string) func(string) error {
	return func(s string) error {
		return assertHttpResponseBodySchemaOnURIWhenAlias(instance, uriSchema, s, "")
	}
}

func assertHttpResponseBodySchemaOnURIWhenAlias(instance *context.Builder, uriSchema, s, alias string) error {
	resp, err := getHttpResponse(instance, alias)
	if err != nil {
		return err
	}
	valueOf, err := instance.GetValue(s)
	if err != nil {
		return err
	}
	key := uriSchema + "://" + s
	schemaLoader, hasValue := schemaCache[key]
	if !hasValue {
		binary, prob := fetchFromURI(instance, uriSchema, valueOf.(string))
		if prob != nil {
			return prob
		}
		schemaLoader = gojsonschema.NewBytesLoader(binary)
		schemaCache[key] = schemaLoader
	}
	documentLoader := gojsonschema.NewBytesLoader(resp.body)
	r, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}
	if !r.Valid() {
		return errors.New(
			strings.Join(funk.Map(r.Errors(), func(desc gojsonschema.ResultError) string {
				return fmt.Sprintf("- %s", desc)
			}).([]string), "\n"),
		)
	}
	return nil
}

func assertHttpResponseStatusCodeWhenAlias(instance *context.Builder, statusCode int, alias string) error {
	var response, err = getHttpResponse(instance, alias)
	if err != nil {
		return err
	}
	if response.statusCode != statusCode {
		return fmt.Errorf("expected status code %d, got %d", statusCode, response.statusCode)
	}
	return nil
}

func getHttpResponse(instance *context.Builder, alias string) (*Response, error) {
	var req, err = getHttpRequest(instance, alias)
	if err != nil {
		return nil, err
	}
	if req.response == nil {
		if err = doSubmitHttpRequest(instance, req); err != nil {
			return nil, err
		}
	}
	return req.response, nil
}

func fetchFromURI(instance *context.Builder, schemaType, value string) ([]byte, error) {
	var b []byte
	var err error
	switch schemaType {
	case "http", "https":
		b, err = pkg.ReadFromURL(instance.GetResourcesHttpClient(), value)
	case "file":
		b, err = pkg.ReadFromFile(instance.GetWorkDirectory(), value)
	default:
		return nil, fmt.Errorf("unsupported URI schema %s", value)
	}
	return b, err
}

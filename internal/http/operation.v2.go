package internal

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
	"github.com/raitonbl/coverup/pkg"
	"strings"
)

const defaultAlias = ""
const componentType = "HttpRequest"

func CreateHttpRequest(instance *context.Builder) func(string) error {
	f := CreateHttpRequestWithAlias(instance)
	return func(s string) error {
		return f(s)
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

func CreateHttpRequestWithAlias(instance *context.Builder) func(string) error {
	return func(alias string) error {
		return instance.WithComponent(componentType, &Request{
			headers: make(map[string]string),
		}, alias)
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
		var b []byte
		var err error
		switch uriSchema {
		case "http", "https":
			b, err = pkg.ReadFromURL(instance.GetResourcesHttpClient(), uri)
		case "file":
			b, err = pkg.ReadFromFile(instance.GetWorkDirectory(), uri)
		default:
			return fmt.Errorf("unsupported URI schema %s", uri)
		}
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

func AssertHttpResponseStatusCode(instance *context.Builder) func(statusCode int) error {
	return func(statusCode int) error {
		return assertHttpResponseStatusCodeWhenAlias(instance, statusCode, "")
	}
}

func AssertHttpResponseStatusCodeWhenAlias(instance *context.Builder) func(int, string) error {
	return func(statusCode int, alias string) error {
		return assertHttpResponseStatusCodeWhenAlias(instance, statusCode, alias)
	}
}

func assertHttpResponseStatusCodeWhenAlias(instance *context.Builder, statusCode int, alias string) error {
	var req, err = getHttpRequest(instance, alias)
	if err != nil {
		return err
	}
	// TODO: FETCH RESPONSE FROM REQ(OR APPLY REQUEST)
	if req.response.statusCode != statusCode {
		return fmt.Errorf("expected status code %d, got %d", statusCode, req.response.statusCode)
	}
	return nil
}

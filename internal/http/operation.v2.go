package internal

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
	"strings"
)

const componentType = "HttpRequest"

func CreateHttpRequest(instance *context.Builder) func(string) error {
	f := CreateHttpRequestWithAlias(instance)
	return func(s string) error {
		return f(s)
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
		var req, err = getRequest(instance)
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
		var req, err = getRequest(instance)
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
		return nil
	}
}

func SetRequestBodyFromURI(instance *context.Builder, uriSchema string) func(string) error {
	return func(v string) error {
		return nil
	}
}

func getRequest(instance *context.Builder) (*Request, error) {
	valueOf := instance.GetComponent(componentType, "")
	if r, isHttpRequest := valueOf.(*Request); isHttpRequest {
		return r, nil
	} else {
		return nil, fmt.Errorf("before setting %s.headers, please define %s", componentType, componentType)
	}
}

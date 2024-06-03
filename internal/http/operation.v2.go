package internal

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
	"strings"
)

const componentType = "HttpRequest"

func CreateHttpRequest(instance *context.Builder) func(string) error {
	f := CreateHttpRequestAndAlias(instance)
	return func(s string) error {
		return f(s)
	}
}

func CreateHttpRequestAndAlias(instance *context.Builder) func(string) error {
	return func(alias string) error {
		return instance.WithComponent(componentType, &Request{
			headers: make(map[string]string),
		}, alias)
	}
}

func CreateHttpRequestHeaders(instance *context.Builder) func(*godog.Table) error {
	return func(table *godog.Table) error {
		var req, err = getRequest(instance)
		if err != nil {
			return err
		}
		if req.headers == nil {
			req.headers = make(map[string]string)
		}
		for _, row := range table.Rows {
			req.headers[row.Cells[0].Value] = row.Cells[1].Value
		}
		return nil
	}
}

func CreateHttpRequestOperation(instance *context.Builder, method string) func(string) error {
	return func(url string) error {
		var req, err = getRequest(instance)
		if err != nil {
			return err
		}
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

func getRequest(instance *context.Builder) (*Request, error) {
	valueOf := instance.GetComponent(componentType, "")
	if r, isHttpRequest := valueOf.(*Request); isHttpRequest {
		return r, nil
	} else {
		return nil, fmt.Errorf("before setting %s.headers, please define %s", componentType, componentType)
	}
}

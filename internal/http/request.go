package internal

import (
	"fmt"
	"strings"
)

type Request struct {
	method    string
	serverURL string
	uri       string
	body      []byte
	response  *Response
	headers   map[string]string
}

func (instance Request) GetPathValue(x string) (any, error) {
	switch x {
	case "method":
		return instance.method, nil
	case "server_url":
		return instance.serverURL, nil
	case "uri":
		return instance.uri, nil
	case "body":
		return string(instance.body), nil
	case "response":
		{
			if instance.response == nil {
				return nil, fmt.Errorf("%s mustn't be undefined", x)
			}
			return instance.response, nil
		}
	case "headers":
		return instance.headers, nil
	default:
		{
			if strings.HasPrefix(x, "response.") {
				if instance.response == nil {
					return nil, fmt.Errorf("%s mustn't be undefined", x)
				}
				return instance.response.GetPathValue(x[9:])
			}
			return nil, fmt.Errorf("cannot resolve %s", x)
		}
	}
}

type Response struct {
	statusCode int
	body       []byte
	pathCache  map[string]any
	headers    map[string]string
}

func (instance *Response) GetPathValue(x string) (any, error) {
	switch x {
	case "status_code":
		return instance.statusCode, nil
	case "headers":
		return instance.headers, nil
	case "body":
		return instance.body, nil
	default:
		{
			if strings.HasPrefix(x, "body.") {
				return instance.getValueFromPath(x[5:])
			}
			return nil, fmt.Errorf("cannot resolve %s", x)
		}
	}
}

func (instance *Response) getContentType() string {
	valueOf, hasValue := instance.headers["content-type"]
	if hasValue {
		return valueOf
	}
	return "text/plain"
}

func (instance *Response) getValueFromPath(x string) (any, error) {
	valueOf, hasValue := instance.pathCache[x]
	if hasValue {
		return valueOf, nil
	}
	valueOf, err := extractValueFromResponsePath(instance.getContentType(), x, instance.body)
	if err != nil {
		return nil, err
	}
	if instance.pathCache == nil {
		instance.pathCache = make(map[string]any)
	}
	instance.pathCache[x] = valueOf
	return valueOf, nil
}

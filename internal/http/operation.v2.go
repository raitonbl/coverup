package internal

import (
	"fmt"
	"github.com/raitonbl/coverup/internal/context"
)

type Request struct {
	method   string
	url      string
	body     []byte
	response *Response
	headers  map[string]string
}

type Response struct {
	statusCode int
	body       []byte
	headers    map[string]string
}

type Builder struct {
	request  *Request
	context  context.Context
	requests map[string]Request
}

func (instance *Builder) withHttpRequest() error {
	return instance.withHttpRequestAndAlias("")
}

func (instance *Builder) withHttpRequestAndAlias(name string) error {
	req := Request{
		headers: make(map[string]string),
	}
	if name != "" {
		_, hasValue := instance.requests[name]
		if hasValue {
			return fmt.Errorf("HttpRequest with alias %s cannot be defined more than once", name)
		}
		instance.requests[name] = req
	}
	instance.request = &req
	return nil
}

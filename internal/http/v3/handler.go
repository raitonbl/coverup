package v3

import (
	"github.com/cucumber/godog"
	v3 "github.com/raitonbl/coverup/internal/v3"
)

type HttpContext struct {
	ctx v3.ScenarioContext
}

func (instance *HttpContext) WithRequest() error {
	return instance.withRequest("")
}

func (instance *HttpContext) WithRequestWhenAlias(alias string) error {
	return instance.withRequest(alias)
}

func (instance *HttpContext) withRequest(alias string) error {
	return nil
}

func (instance *HttpContext) WithHeaders(table *godog.Table) error {
	return nil
}

func (instance *HttpContext) WithHeader(name, value string) error {
	return nil
}

func (instance *HttpContext) WithMethod(method string) error {
	return nil
}

func (instance *HttpContext) WithPath(path string) error {
	return nil
}

func (instance *HttpContext) WithHttpPath(url string) error {
	return nil
}

func (instance *HttpContext) WithHttpsPath(url string) error {
	return nil
}

func (instance *HttpContext) WithServerURL(url string) error {
	return nil
}

func (instance *HttpContext) WithBody(body *godog.DocString) error {
	return nil
}

func (instance *HttpContext) WithBodyFileURI(value string) error {
	return nil
}

func (instance *HttpContext) WithAcceptHeader(value string) error {
	return nil
}

func (instance *HttpContext) WithContentTypeHeader(value string) error {
	return nil
}

func (instance *HttpContext) WithFormEncType(value string) error {
	return nil
}

func (instance *HttpContext) WithFormAttribute(name, value string) error {
	return nil
}

func (instance *HttpContext) WithFormFile(name, value string) error {
	return nil
}

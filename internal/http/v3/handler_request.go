package v3

import (
	"fmt"
	"github.com/cucumber/godog"
	v3 "github.com/raitonbl/coverup/internal/v3"
	"github.com/raitonbl/coverup/pkg"
	"strings"
)

const fileURISchema = "file://"
const ComponentType = "HttpRequest"

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
	return instance.ctx.Register(ComponentType, &HttpRequest{
		headers: make(map[string]string),
	}, alias)
}

func (instance *HttpContext) onHttpRequest(f func(*HttpRequest) error) error {
	valueOf, err := instance.ctx.GetComponent(ComponentType, "")
	if err != nil {
		return err
	}
	if r, isHttpRequest := valueOf.(*HttpRequest); isHttpRequest {
		return f(r)
	}
	return fmt.Errorf("please define %s before using it", ComponentType)
}

func (instance *HttpContext) WithHeaders(table *godog.Table) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		if req.headers == nil {
			req.headers = make(map[string]string)
		}
		for _, row := range table.Rows {
			key := row.Cells[0].Value
			valueOf, prob := instance.ctx.GetValue(row.Cells[1].Value)
			if prob != nil {
				return fmt.Errorf(`cannot determine header "%s" value. caused by:%v"`, key, prob)
			}
			err := instance.withHeader(req, key, fmt.Sprintf("%v", valueOf))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (instance *HttpContext) WithHeader(name, value string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		return instance.withHeader(req, name, value)
	})
}

func (instance *HttpContext) withHeader(req *HttpRequest, name, value string) error {
	if req.headers == nil {
		req.headers = make(map[string]string)
	}
	valueOf, err := instance.ctx.GetValue(value)
	if err != nil {
		return err
	}
	req.headers[name] = fmt.Sprintf("%v", valueOf)
	return nil
}

func (instance *HttpContext) WithMethod(method string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		switch method {
		case "OPTIONS":
			return instance.withMethod(req, method)
		case "HEAD":
			return instance.withMethod(req, method)
		case "GET":
			return instance.withMethod(req, method)
		case "PUT":
			return instance.withMethod(req, method)
		case "POST":
			return instance.withMethod(req, method)
		case "PATCH":
			return instance.withMethod(req, method)
		case "DELETE":
			return instance.withMethod(req, method)
		default:
			return fmt.Errorf("cannot assign %s.method to %s ", ComponentType, method)
		}
	})
}

func (instance *HttpContext) withMethod(req *HttpRequest, method string) error {
	req.method = method
	return nil
}

func (instance *HttpContext) WithPath(path string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		valueOf, err := instance.ctx.GetValue(path)
		if err != nil {
			return err
		}
		req.path = fmt.Sprintf("%v", valueOf)
		return nil
	})
}

func (instance *HttpContext) WithHttpPath(url string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("http://%v", valueOf)
		return nil
	})
}

func (instance *HttpContext) WithHttpsPath(url string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("https://%v", valueOf)
		return nil
	})
}

func (instance *HttpContext) WithServerURL(url string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("%v", valueOf)
		return nil
	})
}

func (instance *HttpContext) WithBody(body *godog.DocString) error {
	return instance.withBody(func() ([]byte, error) {
		return []byte(body.Content), nil
	})
}

func (instance *HttpContext) withBody(s func() ([]byte, error)) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		binary, err := s()
		if err != nil {
			return err
		}
		req.form = nil
		req.body = binary
		return nil
	})
}

func (instance *HttpContext) WithBodyFileURI(value string) error {
	return instance.withBody(func() ([]byte, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return instance.readURI(fmt.Sprintf("%s%v", fileURISchema, valueOf))
	})
}

func (instance *HttpContext) readURI(uri string) ([]byte, error) {
	var binary []byte
	var err error
	switch {
	case strings.HasPrefix(uri, "http://"), strings.HasPrefix(uri, "https://"):
		binary, err = pkg.ReadFromURL(instance.ctx.GetResourcesHttpClient(), uri)
	case strings.HasPrefix(uri, fileURISchema):
		binary, err = pkg.ReadFromFile(instance.ctx.GetWorkDirectory(), uri)
	default:
		return nil, fmt.Errorf("unsupported URI: %s", uri)
	}
	return binary, err
}

func (instance *HttpContext) WithAcceptHeader(value string) error {
	return instance.WithHeader("Accept", value)
}

func (instance *HttpContext) WithContentTypeHeader(value string) error {
	return instance.WithHeader("Content-Type", value)
}

func (instance *HttpContext) WithFormEncType(value string) error {
	return instance.onForm(func(form *Form) error {
		if value == "multipart/form-data" || value == "application/x-www-form-urlencoded" {
			form.encType = value
			return nil
		}
		return fmt.Errorf("encType %s is not supported", value)
	})
}

func (instance *HttpContext) WithFormAttribute(name, value string) error {
	return instance.onFormAttribute(name, func() (any, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%v", valueOf), nil
	})
}

func (instance *HttpContext) WithFormFile(name, value string) error {
	return instance.onFormAttribute(name, func() (any, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return instance.readURI(fmt.Sprintf("%s%v", fileURISchema, valueOf))
	})
}

func (instance *HttpContext) onFormAttribute(name string, f func() (any, error)) error {
	return instance.onForm(func(form *Form) error {
		valueOf, err := f()
		if err != nil {
			return err
		}
		form.attributes[name] = valueOf
		return nil
	})
}

func (instance *HttpContext) onForm(s func(*Form) error) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		form := req.form
		if form == nil {
			req.form = &Form{
				encType:    "multipart/form-data",
				attributes: make(map[string]any),
			}
			form = req.form
		}
		req.body = nil
		return s(form)
	})
}

func (instance *HttpContext) SubmitHttpRequest() error {
	return nil
}

func (instance *HttpContext) SubmitHttpRequestOnBehalfOfEntity(id string) error {
	return nil
}

func (instance *HttpContext) SubmitNamedHttpRequest(alias string) error {
	return nil
}

func (instance *HttpContext) SubmitNamedHttpRequestOnBehalfOfEntity(alias, id string) error {
	return nil
}

func (instance *HttpContext) AssertResponseStatusCode(statusCode int) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseStatusCode(alias string, statusCode int) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseExactHeaders(alias string, table *godog.Table) error {
	return nil
}

func (instance *HttpContext) AssertResponseExactHeaders(table *godog.Table) error {
	return nil
}

func (instance *HttpContext) AssertResponseContainsHeaders(table *godog.Table) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseContainsHeaders(alias string, table *godog.Table) error {
	return nil
}

func (instance *HttpContext) AssertResponseHeader(name, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponseContentType(value string) error {
	return instance.AssertResponseHeader("content-type", value)
}

func (instance *HttpContext) AssertNamedHttpRequestResponseContentType(alias, value string) error {
	return instance.AssertNamedHttpRequestResponseHeader(alias, "content-type", value)
}

func (instance *HttpContext) AssertNamedHttpRequestResponseHeader(alias string, header string, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponseIsValidAgainstSchema(file string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseIsValidAgainstSchema(alias, file string) error {
	return nil
}

func (instance *HttpContext) AssertResponseBodyPathEqualsTo(alias, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponseBodyPathEqualsToValue(alias, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponseBodyPathIsEqualToFloat64(p string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponseBodyPathIsEqualToBoolean(alias string, value bool) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseBodyPathEqualsTo(alias, value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseBodyPathEqualsToValue(alias, value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseBodyPathIsEqualToFloat64(alias string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseBodyPathIsEqualToBoolean(alias string, value bool) error {
	return nil
}

func (instance *HttpContext) AssertResponseBodyEqualsToFile(value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseBodyEqualsToFile(alias string, value bool) error {
	return nil
}

func (instance *HttpContext) AssertResponseBodyEqualsTo(value *godog.DocString) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseBodyEqualsTo(alias string, value *godog.DocString) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathEndsWith(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathContains(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathStartsWith(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertWhileIgnoringCaseThatResponsePathEndsWith(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertWhileIgnoringCaseThatResponsePathContains(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertWhileIgnoringCaseThatResponsePathStartsWith(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathMatchesPattern(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathMatchesPattern(alias, k, value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathEndsWith(alias, k, value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathContains(alias, k, value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathStartsWith(alias, k, value string) error {
	return nil
}

func (instance *HttpContext) AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathEndsWith(alias, k, value string) error {
	return nil
}

func (instance *HttpContext) AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathContains(alias, k, value string) error {
	return nil
}

func (instance *HttpContext) AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathStartsWith(alias, k, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsTime(k string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsDate(k string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsDateTime(k string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsTime(alias, k string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsDate(alias, k string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsDateTime(alias, k string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsSame(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsAfter(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsSameOrAfter(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsSameOrBefore(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsBefore(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsSame(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsAfter(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsSameOrAfter(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsSameOrBefore(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsBefore(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathLengthIs(k string, float64 string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathLengthIs(alias, k string, float64 string) error {
	return nil
}

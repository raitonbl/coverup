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

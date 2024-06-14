package http

import (
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"reflect"
)

func newResponseBodyIsEqualToHandler(instance *HttpContext, opts HandlerOpts) any {
	f := func(alias string, value *godog.DocString) error {
		return execOnResponseBodyEqualsToHandler(instance, opts, alias, []byte(value.Content))
	}
	if !opts.isAliasedFunction {
		return func(value *godog.DocString) error {
			return f("", value)
		}
	}
	return f
}

func newResponseBodyIsEqualToFileHandler(instance *HttpContext, opts HandlerOpts) any {
	f := func(alias string, value string) error {
		binary, err := doGetFromURI(instance, fileUriScheme, value)
		if err != nil {
			return err
		}
		return execOnResponseBodyEqualsToHandler(instance, opts, alias, binary)
	}
	if !opts.isAliasedFunction {
		return func(value string) error {
			return f("", value)
		}
	}
	return f
}

func execOnResponseBody(instance *HttpContext, alias string, f func(*HttpRequest, *HttpResponse) error) error {
	return instance.onNamedHttpRequest(alias, func(req *HttpRequest) error {
		if req.response == nil {
			if alias == "" {
				return fmt.Errorf(`%s needs to be submitted before making assertions`, ComponentType)
			} else {
				return fmt.Errorf(`%s["%s"] needs to be submitted before making assertions`, ComponentType, alias)
			}
		}
		return f(req, req.response)
	})
}

func execOnResponseBodyEqualsToHandler(instance *HttpContext, opts HandlerOpts, alias string, binary []byte) error {
	return execOnResponseBody(instance, alias, func(req *HttpRequest, res *HttpResponse) error {
		var predicate func() (bool, error)
		if res.headers["content-type"] == "application/json" {
			predicate = func() (bool, error) {
				fromResponse := map[string]any{}
				if prob := json.Unmarshal(res.body, &fromResponse); prob != nil {
					return false, prob
				}
				fromValue := map[string]any{}
				if prob := json.Unmarshal(binary, &fromValue); prob != nil {
					return false, prob
				}
				return reflect.DeepEqual(fromResponse, fromValue), nil
			}
		} else {
			predicate = func() (bool, error) {
				if string(res.body) != string(binary) {
					return false, nil
				}
				return true, nil
			}
		}
		r, err := predicate()
		if err != nil {
			return err
		}
		if opts.isAffirmationExpected == r {
			return nil
		}

		if opts.isAffirmationExpected {
			return fmt.Errorf("response isn't equal to expectation.\n%s", string(res.body))
		}
		return fmt.Errorf("response  shouldn't match expectation")
	})
}

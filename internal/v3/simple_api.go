package v3

import (
	"fmt"
	"regexp"
	"strings"
)

type HandlerOpts struct {
	isAffirmation    bool
	isAliasAware     bool
	ignoreCase       bool
	interpolateValue bool
}

type HandlerFactory func(instance *HttpContext, opts HandlerOpts) any

func createResponseBodyPathEqualTo(instance *HttpContext, opts HandlerOpts) any {
	f := createResponseBodyPathEqualsTo(instance, opts)
	if opts.isAliasAware {
		return func(expr, alias, compareTo string) error {
			return f.(func(string, string, any) error)(expr, alias, compareTo)
		}
	}
	return func(expr, compareTo string) error {
		return f.(func(string, any) error)(expr, compareTo)
	}
}

func createResponseBodyPathEqualToFloat64(instance *HttpContext, opts HandlerOpts) any {
	f := createResponseBodyPathEqualsTo(instance, opts)
	if opts.isAliasAware {
		return func(expr, alias string, compareTo float64) error {
			return f.(func(string, string, any) error)(expr, alias, compareTo)
		}
	}
	return func(expr string, compareTo float64) error {
		return f.(func(string, any) error)(expr, compareTo)
	}
}

func createResponseBodyPathEqualToBoolean(instance *HttpContext, opts HandlerOpts) any {
	f := createResponseBodyPathEqualsTo(instance, opts)
	if opts.isAliasAware {
		return func(expr, alias string, compareTo string) error {
			valueOf := compareTo == "true"
			return f.(func(string, string, any) error)(expr, alias, valueOf)
		}
	}
	return func(expr string, compareTo string) error {
		valueOf := compareTo == "true"
		return f.(func(string, any) error)(expr, valueOf)
	}
}

func createResponseBodyPathEqualsTo(instance *HttpContext, opts HandlerOpts) any {
	f := func(expr, alias string, compareTo any) error {
		return instance.onNamedHttpRequestResponseBodyPath(expr, alias, func(_ *HttpRequest, response *HttpResponse, value any) error {
			if (value == compareTo) == opts.isAffirmation {
				return nil
			}
			condition := "must"
			if !opts.isAffirmation {
				condition += "n't"
			}
			if alias == "" {
				return fmt.Errorf(`$%s=%v %s be equal to %v`, expr, value, condition, compareTo)
			}
			return fmt.Errorf(`%s.$%s=%v %s be equal to %v`, alias, expr, condition, value, compareTo)
		})
	}
	if opts.isAliasAware {
		return f
	}
	return func(expr string, compareTo any) error {
		return f(expr, "", compareTo)
	}
}

func createResponseBodyPathContains(instance *HttpContext, opts HandlerOpts) any {
	return createResponseBodyPathThenExecuteStringOperation(instance, "contain", opts, strings.Contains)
}

func createResponseBodyPathStartsWith(instance *HttpContext, opts HandlerOpts) any {
	return createResponseBodyPathThenExecuteStringOperation(instance, "starts with", opts, strings.HasPrefix)
}

func createResponseBodyPathMatchesPattern(instance *HttpContext, opts HandlerOpts) any {
	return createResponseBodyPathThenExecuteStringOperation(instance, "matches pattern", opts, func(fromResponse string, value string) bool {
		r, err := regexp.Compile(value)
		if err != nil {
			//TODO LOG ERROR
			return false
		}
		return r.Match([]byte(fromResponse))
	})
}

func createResponseBodyPathEndsWith(instance *HttpContext, opts HandlerOpts) any {
	return createResponseBodyPathThenExecuteStringOperation(instance, "ends with", opts, strings.HasSuffix)
}

func createResponseBodyPathThenExecuteStringOperation(instance *HttpContext, operation string, opts HandlerOpts, predicate func(string, string) bool) any {
	f := func(expr, alias string, c string) error {
		return instance.onNamedHttpRequestResponseBodyPath(expr, alias, func(_ *HttpRequest, response *HttpResponse, value any) error {
			if value == nil {
				if alias == "" {
					return fmt.Errorf(`$%s mustn't be undefined`, expr)
				}
				return fmt.Errorf(`%s.$%s mustn't be undefined`, alias, expr)
			}
			valueOf, isText := value.(string)
			if !isText {
				if alias == "" {
					return fmt.Errorf(`$%s must be a string but got %v`, expr, value)
				}
				return fmt.Errorf(`%s.$%s must be a string but got %v`, alias, expr, value)
			}
			compareTo := c
			if opts.interpolateValue {
				v, err := instance.ctx.GetValue(c)
				if err != nil {
					return err
				}
				compareTo, isText = v.(string)
				if !isText {
					if alias == "" {
						return fmt.Errorf(`%s must be a string but got %v`, c, value)
					}
					return fmt.Errorf(`%s must be a string but got %v`, c, value)
				}
			}
			v1 := valueOf
			v2 := compareTo
			if opts.ignoreCase {
				v1 = strings.ToUpper(v1)
				v2 = strings.ToUpper(v2)
			}
			r := predicate(v1, v2)
			if r == opts.isAffirmation {
				return nil
			}
			condition := "must"
			if !opts.isAffirmation {
				condition += "n't"
			}
			if alias == "" {
				return fmt.Errorf(`$%s=%v %s %s %v`, expr, value, condition, operation, compareTo)
			}
			return fmt.Errorf(`%s.$%s=%v %s %s %v`, alias, expr, condition, operation, value, compareTo)
		})
	}
	if opts.isAliasAware {
		return f
	}
	return func(expr string, compareTo string) error {
		return f(expr, "", compareTo)
	}
}

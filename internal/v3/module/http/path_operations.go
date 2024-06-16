package http

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg/api"
	"github.com/raitonbl/coverup/pkg/checks"
	"strings"
)

const (
	anyStringRegex        = `"([^"]*)"`
	boolRegex             = `(true|false)`
	anyNumber             = `(-?\d+(\.\d+)?)`
	resolvableStringRegex = `"` + valueRegex + `"`
)

type PathOperationSettings struct {
	IgnoreCase  bool
	ValueRegexp string
}

type PathOperations struct {
	Line                 string
	ExpressionPattern    string
	PhraseFactory        func(string) []string
	AliasedPhraseFactory func(string) []string
	ExtractFromResponse  func(response *Response, expr string) (any, error)
}

func (instance *PathOperations) New(ctx api.StepDefinitionContext) {
	instance.enabledEqualsToSupport(ctx)
	instance.enabledEndsWithSupport(ctx)
	instance.enabledStartsWithSupport(ctx)
	instance.enabledContainsSupport(ctx)
	// matches pattern
	// is lesser
	// is greater
	// is lesser or equal to
	// is greater or equal to
}

func (instance *PathOperations) enabledEqualsToSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "be equal to", true, instance.equalsToAssertionFactory)
}

func (instance *PathOperations) enabledStartsWithSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "start with", true, instance.startsWithAssertionFactory)
}

func (instance *PathOperations) enabledEndsWithSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "end with", true, instance.endsWithAssertionFactory)
}

func (instance *PathOperations) enabledContainsSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "contain", true, instance.containsAssertionFactory)
}

func (instance *PathOperations) enabledSupportTo(ctx api.StepDefinitionContext, operation string, allowsIgnoreCase bool, f func(options FactoryOpts[PathOperationSettings]) api.HandlerFactory) {
	verbs := []string{"should", "shouldn't"}
	aliases := []string{"", httpRequestRegex}
	args := []string{anyStringRegex, resolvableStringRegex, valueRegex, anyNumber, boolRegex}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a specific %s response %s %s the specified value", instance.Line, fmt.Sprintf("%s %s", verb, operation), ComponentType),
			Options:     make([]api.Option, 0),
		}
		for _, alias := range aliases {
			for _, arg := range args {
				numberOfOptions := 1
				supportsIgnoreCase := allowsIgnoreCase && (arg == anyStringRegex || arg == resolvableStringRegex || arg == valueRegex)
				if supportsIgnoreCase {
					numberOfOptions = 2
				}
				for i := 0; i < numberOfOptions; i++ {
					isIgnoreCase := i == 1
					var phrases []string
					format := fmt.Sprintf(`%s %s %s %s %s`, instance.Line, instance.ExpressionPattern, verb, operation, arg)
					if alias == aliases[0] {
						phrases = instance.PhraseFactory(format)
					} else {
						phrases = instance.AliasedPhraseFactory(format)
					}
					for _, p := range phrases {
						phrase := p
						if isIgnoreCase {
							phrase += ", ignoring case"
						}
						phrase += "$"
						options := FactoryOpts[PathOperationSettings]{
							Settings: &PathOperationSettings{
								ValueRegexp: arg,
								IgnoreCase:  isIgnoreCase,
							},
							AssertTrue:                  verb == verbs[0],
							AssertAlias:                 alias == aliases[1],
							ResolveValueBeforeAssertion: arg != anyNumber && arg != boolRegex,
						}
						step.Options = append(step.Options, api.Option{
							Regexp:         phrase,
							Description:    step.Description,
							HandlerFactory: f(options),
						})
					}
				}
			}
		}
		ctx.Then(step)
	}
}

func (instance *PathOperations) startsWithAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return instance.stringOperationAssertionFactory(options, strings.HasPrefix)
}

func (instance *PathOperations) endsWithAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return instance.stringOperationAssertionFactory(options, strings.HasSuffix)
}

func (instance *PathOperations) containsAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return instance.stringOperationAssertionFactory(options, strings.Contains)
}

func (instance *PathOperations) stringOperationAssertionFactory(options FactoryOpts[PathOperationSettings], predicate func(string, string) bool) api.HandlerFactory {
	return func(c api.ScenarioContext) any {
		f := func(alias string, expr string, v any) error {
			compareTo := v
			res, err := getHttpResponse(c, alias)
			if err != nil {
				return err
			}
			if addr, isString := compareTo.(string); isString && options.ResolveValueBeforeAssertion {
				valueOf, prob := c.Resolve(addr)
				if prob != nil {
					return prob
				}
				compareTo = valueOf
			}
			fromResponse, err := instance.ExtractFromResponse(res, expr)
			if err != nil {
				return err
			}
			if !checks.IsString(fromResponse) {
				return fmt.Errorf("response.%s[%s]=%v should be a string", instance.Line, expr, fromResponse)
			}
			if !checks.IsString(compareTo) {
				return fmt.Errorf("argument should be a string")
			}
			r := false
			if options.Settings.IgnoreCase {
				r = predicate(strings.ToUpper(fromResponse.(string)), strings.ToUpper(compareTo.(string)))
			} else {
				r = predicate(fromResponse.(string), compareTo.(string))
			}
			if r == options.AssertTrue {
				return nil
			}
			return fmt.Errorf("response.%s[%s]=%v", instance.Line, expr, fromResponse)
		}
		return instance.createHandler(options, f)
	}
}

func (instance *PathOperations) equalsToAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return func(c api.ScenarioContext) any {
		f := func(alias string, expr string, v any) error {
			compareTo := v
			res, err := getHttpResponse(c, alias)
			if err != nil {
				return err
			}
			if addr, isString := compareTo.(string); isString && options.ResolveValueBeforeAssertion {
				valueOf, prob := c.Resolve(addr)
				if prob != nil {
					return prob
				}
				compareTo = valueOf
			}
			valueOf, err := instance.ExtractFromResponse(res, expr)
			if err != nil {
				return err
			}
			isEqualTo := false
			if options.Settings.IgnoreCase && checks.IsString(valueOf) && checks.IsString(compareTo) {
				isEqualTo = strings.ToUpper(valueOf.(string)) == strings.ToUpper(compareTo.(string))
			} else {
				isEqualTo = valueOf == compareTo
			}
			if isEqualTo == options.AssertTrue {
				return nil
			}
			return fmt.Errorf("response.%s[%s]=%v", instance.Line, expr, valueOf)
		}
		return instance.createHandler(options, f)
	}
}

func (instance *PathOperations) createHandler(options FactoryOpts[PathOperationSettings], f any) any {
	if !options.AssertAlias {
		alias := ""
		if options.Settings.ValueRegexp == anyNumber {
			return func(expr string, value float64) error {
				return f.(func(string, string, any) error)(alias, expr, value)
			}
		} else if options.Settings.ValueRegexp == boolRegex {
			return func(expr string, value bool) error {
				return f.(func(string, string, any) error)(alias, expr, value)
			}
		} else {
			return func(expr, value string) error {
				return f.(func(string, string, any) error)(alias, expr, value)
			}
		}
	}
	if options.Settings.ValueRegexp == anyNumber {
		return func(alias string, expr string, value float64) error {
			return f.(func(string, string, any) error)(alias, expr, value)
		}
	} else if options.Settings.ValueRegexp == boolRegex {
		return func(alias string, expr string, value bool) error {
			return f.(func(string, string, any) error)(alias, expr, value)
		}
	} else {
		return func(alias string, expr, value string) error {
			return f.(func(string, string, any) error)(alias, expr, value)
		}
	}
}

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
	GetValue             func(response *Response, expr string) (any, error)
}

func (instance *PathOperations) New(ctx api.StepDefinitionContext) {
	instance.enabledEqualsToSupport(ctx)
	// is equal to [ignore case]
	// starts with [ignore case]
	// ends with [ignore case]
	// contains [ignore case]
	// matches pattern
	// is lesser
	// is greater
	// is lesser or equal to
	// is greater or equal to
}

func (instance *PathOperations) enabledEqualsToSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "equal", instance.equalsToAssertionFactory)
}

func (instance *PathOperations) enabledSupportTo(ctx api.StepDefinitionContext, operation string, f func(options FactoryOpts[PathOperationSettings]) api.HandlerFactory) {
	verbs := []string{"is", "is not"}
	aliases := []string{"", httpRequestRegex}
	args := []string{anyStringRegex, resolvableStringRegex, valueRegex, anyNumber, boolRegex}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a specific %s response %s %s to thes specified value", instance.Line, fmt.Sprintf("%s %s", verb, operation), ComponentType),
			Options:     make([]api.Option, 0),
		}
		for _, alias := range aliases {
			for _, arg := range args {
				var phrases []string
				format := fmt.Sprintf(`%s %s %s %s to %s`, instance.Line, instance.ExpressionPattern, verb, operation, arg)
				if alias == aliases[0] {
					phrases = instance.PhraseFactory(format)
				} else {
					phrases = instance.AliasedPhraseFactory(format)
				}
				supportsIgnoreCase := arg == anyNumber || arg == resolvableStringRegex || arg == valueRegex
				for _, phrase := range phrases {
					numberOfOptions := 1
					if supportsIgnoreCase {
						numberOfOptions = 2
					}
					for i := 0; i < numberOfOptions; i++ {
						options := FactoryOpts[PathOperationSettings]{
							Settings: &PathOperationSettings{
								ValueRegexp: arg,
								IgnoreCase:  i == 1,
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
			valueOf, err := instance.GetValue(res, expr)
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
		return instance.createFn(options, f)
	}
}

func (instance *PathOperations) createFn(options FactoryOpts[PathOperationSettings], f any) any {
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

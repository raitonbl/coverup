package http

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg/api"
	"github.com/raitonbl/coverup/pkg/checks"
)

type ThenHttpResponseStepFactory struct {
	ThenBasicHttpResponseStepFactory
}

func (instance *ThenHttpResponseStepFactory) New(ctx api.StepDefinitionContext) {
	instance.ThenBasicHttpResponseStepFactory.New(ctx)
	instance.enableBodyPathStepSupport(ctx)
}

func (instance *ThenHttpResponseStepFactory) enableBodyPathStepSupport(ctx api.StepDefinitionContext) {
	ops := PathOperations{
		ExpressionPattern:          `\$(\S+)`,
		Line:                       "body",
		ConvertToNumberIfNecessary: false,
		PhraseFactory:              createResponseLinePart,
		AliasedPhraseFactory:       createAliasedResponseLinePart,
		ExtractFromResponse: func(res *Response, expr string) (any, error) {
			if !checks.IsAnyOf(res.headers["content-type"], "application/json", "application/problem+json") {
				return nil, fmt.Errorf(`content-type must be either "application/json" or "application/problem+json"`)
			}
			return res.JSONPath(expr)
		},
	}
	ops.New(ctx)
}

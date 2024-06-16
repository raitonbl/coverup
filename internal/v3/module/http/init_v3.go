package http

import (
	"github.com/raitonbl/coverup/pkg/api"
)

func OnV3(ctx api.StepDefinitionContext) {
	arr := []api.StepFactory{
		&GivenHttpRequestStepFactory{},
		&ThenHttpResponseStepFactory{},
	}
	for _, each := range arr {
		each.New(ctx)
	}
}

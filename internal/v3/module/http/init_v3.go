package http

import (
	"github.com/raitonbl/coverup/pkg/api"
)

func OnV3(ctx api.StepDefinitionContext) {
	g := GivenHttpRequestStepFactory{}
	g.New(ctx)
}

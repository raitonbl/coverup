package aws

import "github.com/raitonbl/coverup/pkg/api"

func On(ctx api.StepDefinitionContext) {
	arr := make([]api.StepFactory, 0)
	for _, each := range arr {
		each.New(ctx)
	}
}

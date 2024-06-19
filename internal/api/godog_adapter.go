package api

import (
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/pkg/api"
)

type GoDogAdapter interface {
	api.StepDefinitionContext
	Configure(*godog.ScenarioContext)
}

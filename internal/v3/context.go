package v3

import (
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
)

type ScenarioContext interface {
	context.Context
	GerkhinContext() *godog.ScenarioContext
}

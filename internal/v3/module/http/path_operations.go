package http

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg/api"
)

type PathOperations struct {
	Line      string
	ValueOpts []string
	GetValue  func(expr string) (any, error)
}

func (instance *PathOperations) New(ctx api.StepDefinitionContext) {
	// is equal to [ignore case]
	// starts with [ignore case]
	// ends with [ignore case]
	// contains [ignore case]
	// matches pattern
	// is lesser
	// is greater
	// is lesser or equal to
	// is greater or equal to
	step := api.StepDefinition{
		Description: fmt.Sprintf("Asserts that a specific %s response header is equal to a specific value", ComponentType),
		Options:     make([]api.Option, 0),
	}
	ctx.Then(step)
}

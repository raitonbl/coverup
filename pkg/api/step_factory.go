package api

type StepFactory interface {
	New(ctx StepDefinitionContext)
}

package api

type StepDefinitionContext interface {
	Step(StepDefinition)
	Given(StepDefinition)
	When(StepDefinition)
	Then(StepDefinition)
}

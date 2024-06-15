package api

type Module struct {
	Name          string
	Description   string
	ComponentType string
	Steps         []StepDefinition
}

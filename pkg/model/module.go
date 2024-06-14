package model

type Module struct {
	Name          string
	Description   string
	ComponentType string
	Steps         []StepDefinition
}

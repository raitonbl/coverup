package pkg

// response body [path] is [arg]
// response body [path] isn't [arg]

type Module struct {
	Name          string
	Description   string
	ComponentType string
	Steps         []StepDefinition
}

type StepDefinition struct {
	Format      string
	Description string
}

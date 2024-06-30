package api

const (
	EntityComponentType     = "Entities"
	PropertiesComponentType = "Properties"
)

const (
	NonLiteralStringExpression = `([^"]*)`
	LiteralStringExpression    = `"` + NonLiteralStringExpression + `"`
)

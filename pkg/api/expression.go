package api

const (
	PropertyExpression = `{{\s*Properties.[a-zA-Z0-9_]+\s*}}`
	ValueExpression    = `{{\s*([a-zA-Z0-9_]+\.)*[a-zA-Z0-9_]+\s*}}`
	EntityExpression   = `{{\s*` + ComponentType + `.([a-zA-Z0-9_]+)\s*}}`
)

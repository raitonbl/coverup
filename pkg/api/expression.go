package api

const (
	PropertyExpression = `{{\s*Properties.[a-zA-Z0-9_]+\s*}}`
	ValueExpression    = `{{\s*.*\s*}}`
	EntityExpression   = `{{\s*` + EntityComponentType + `.([a-zA-Z0-9_]+)\s*}}`
)

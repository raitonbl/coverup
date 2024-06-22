package api

const (
	PropertyExpression = `{{\s*Properties.[a-zA-Z0-9_]+\s*}}`
	ValueExpression    = `{{s*([a-zA-Z]+\.[a-zA-Z0-9.]+)s*}}`
	EntityExpression   = `{{\s*` + EntityComponentType + `.([a-zA-Z0-9_]+)\s*}}`
)

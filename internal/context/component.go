package context

type Component interface {
	GetPathValue(x string) (any, error)
}

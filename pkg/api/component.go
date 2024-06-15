package api

type Component interface {
	GetPathValue(x string) (any, error)
}

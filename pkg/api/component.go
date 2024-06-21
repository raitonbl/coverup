package api

type Component interface {
	ValueFrom(x string) (any, error)
}

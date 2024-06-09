package pkg

type Component interface {
	GetPathValue(x string) (any, error)
}

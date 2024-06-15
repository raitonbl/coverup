package api

import (
	"io/fs"
)

type ScenarioContext interface {
	GetFS() fs.ReadFileFS
	Resolve(value string) (any, error)
	GetGivenComponent(componentType, alias string) (any, error)
	AddGivenComponent(componentType string, ptr Component, alias string) error
	GetValue(namespace, key string) (any, error)
	SetValue(namespace, key string, value any) error
}

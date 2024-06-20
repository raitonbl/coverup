package entities

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg/api"
)

type Entity interface {
	api.Component
	GetName() string
	GetDescription() string
}

type BasicEntity struct {
	Name        string
	Description string
}

func (instance BasicEntity) GetName() string {
	return instance.Name
}

func (instance BasicEntity) GetDescription() string {
	return instance.Description
}

func (instance BasicEntity) GetPathValue(key string) (any, error) {
	switch key {
	case "Name":
		return instance.Name, nil
	case "Description":
		return instance.Description, nil
	default:
		return nil, fmt.Errorf(`%s not defined`, key)
	}
}

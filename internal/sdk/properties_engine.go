package sdk

import (
	"github.com/raitonbl/coverup/pkg/api"
	"gopkg.in/ini.v1"
	"io/fs"
)

type PropertiesEngine struct {
	Props *ini.File
}

func NewPropertiesEngine(fs fs.ReadFileFS, seq ...string) (ValueResolver, error) {
	if seq == nil {
		return NewPropertiesEngine(fs)
	}
	if len(seq) == 0 {
		return &PropertiesEngine{}, nil
	}
	arr := make([]any, len(seq))
	for index, f := range seq {
		binary, err := fs.ReadFile(f)
		if err != nil {
			return nil, err
		}
		arr[index] = binary
	}
	source := arr[0]
	others := make([]any, 0)
	if len(arr) > 1 {
		others = arr[1:]
	}
	props, err := ini.Load(source, others...)
	if err != nil {
		return nil, err
	}
	return &PropertiesEngine{Props: props}, nil
}

func (instance *PropertiesEngine) ToMap() map[string]string {
	return nil
}

func (instance *PropertiesEngine) ValueFrom(x string) (any, error) {
	if instance.Props == nil {
		return nil, nil
	}
	var valueOf any = instance.Props.Section(ini.DEFAULT_SECTION).Key(x).Value()
	if valueOf == "" {
		valueOf = nil
	}
	return valueOf, nil
}

func (instance *PropertiesEngine) GetType() string {
	return api.PropertiesComponentType
}

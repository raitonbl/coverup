package http

import (
	"encoding/json"
	"github.com/thoas/go-funk"
	"k8s.io/client-go/util/jsonpath"
	"reflect"
)

var jsonPathCache = make(map[string]*jsonpath.JSONPath)

type Request struct {
	form      *Form
	path      string
	method    string
	serverURL string
	body      []byte
	response  *Response
	headers   map[string]string
}

func (instance *Request) ValueFrom(x string) (any, error) {
	panic("implement me")
}

type Form struct {
	encType    string
	attributes map[string]any
}

func (instance *Form) ValueFrom(x string) (any, error) {
	panic("implement me")
}

type Response struct {
	body       []byte
	headers    map[string]string
	statusCode float64
	pathCache  map[string]any
	jsonPath   *JsonPathContext
}

func (instance *Response) ValueFrom(x string) (any, error) {
	panic("implement me")
}

func (instance *Response) JSONPath(expr string) (any, error) {
	if instance.jsonPath == nil {
		instance.jsonPath = &JsonPathContext{
			binary: instance.body,
			cache:  make(map[string]any),
		}
	}
	return instance.jsonPath.Get(expr)
}

type JsonPathContext struct {
	body   any
	binary []byte
	cache  map[string]any
}

func (instance *JsonPathContext) Get(expr string) (any, error) {
	if expr[0] == '.' {
		return instance.Get(expr[1:])
	}
	valueOf, hasValue := instance.cache[expr]
	if hasValue {
		return valueOf, nil
	}
	if instance.body == nil {
		var data any
		if err := json.Unmarshal(instance.binary, &data); err != nil {
			return nil, err
		}
		instance.body = data
	}
	var jp *jsonpath.JSONPath
	if item, isCached := jsonPathCache[expr]; isCached {
		jp = item
	} else {
		jp = jsonpath.New(expr)
		if err := jp.Parse("{." + expr + "}"); err != nil {
			return nil, err
		}
		jsonPathCache[expr] = jp
	}
	results, err := jp.FindResults(instance.body)
	if err != nil {
		return nil, err
	}
	array := funk.FlatMap(results, func(value []reflect.Value) []any {
		return funk.Map(value, func(each reflect.Value) any {
			return each.Interface()
		}).([]any)
	}).([]any)

	if len(array) == 0 {
		valueOf = nil
	} else if len(array) == 1 {
		valueOf = array[0]
	} else {
		valueOf = array
	}
	instance.cache[expr] = valueOf
	return valueOf, nil
}

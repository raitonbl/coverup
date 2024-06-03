package internal

import (
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"k8s.io/client-go/util/jsonpath"
	"reflect"
)

var jpCache = make(map[string]*jsonpath.JSONPath)

func extractValueFromResponsePath(contentType string, expr string, binary []byte) (any, error) {
	switch contentType {
	case "application/json":
		return extractJSONPathValue(binary, expr)
	default:
		return nil, fmt.Errorf("cannot extract $body.%s when content-type=%s", expr, contentType)
	}
}

func extractJSONPathValue(binary []byte, expr string) (any, error) {
	var data interface{}
	if err := json.Unmarshal(binary, &data); err != nil {
		return nil, err
	}
	var jp *jsonpath.JSONPath
	if item, hasValue := jpCache[expr]; hasValue {
		jp = item
	} else {
		jp = jsonpath.New(expr)
		if err := jp.Parse("{." + expr + "}"); err != nil {
			return nil, err
		}
		jpCache[expr] = jp
	}
	results, err := jp.FindResults(data)
	if err != nil {
		return nil, err
	}
	array := funk.FlatMap(results, func(value []reflect.Value) []any {
		return funk.Map(value, func(each reflect.Value) any {
			return each.Interface()
		}).([]any)
	}).([]any)

	if len(array) == 0 {
		return nil, nil
	}

	if len(array) == 1 {
		return array[0], nil
	}
	return array, nil
}

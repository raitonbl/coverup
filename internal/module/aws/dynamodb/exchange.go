package dynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"strconv"
	"strings"
)

type GetItemResponse struct {
	cache         map[string]interface{}
	properties    map[string]interface{}
	sdkAttributes map[string]types.AttributeValue
}

func (instance *GetItemResponse) ValueFrom(location string) (interface{}, error) {
	if instance.cache == nil {
		instance.cache = make(map[string]interface{})
	}
	// Check if the value is already in the cache
	if value, found := instance.cache[location]; found {
		return value, nil
	}
	if instance.properties == nil {
		value, err := getGolangValueFrom(&types.AttributeValueMemberM{Value: instance.sdkAttributes})
		if err != nil {
			return nil, err
		}
		instance.sdkAttributes = nil
		instance.properties = value.(map[string]interface{})
	}
	value, err := doValueFrom(instance.properties, location)
	if err != nil {
		return nil, err
	}
	// Store the value in the cache
	instance.cache[location] = value
	return value, nil
}

func doValueFrom(data map[string]interface{}, path string) (interface{}, error) {
	elements := strings.Split(path, ".")
	var current interface{} = data
	for _, element := range elements {
		if strings.Contains(element, "[") && strings.Contains(element, "]") {
			current, element = parseArrayElement(current, element)
			if current == nil {
				return nil, nil
			}
		} else {
			if currentMap, ok := current.(map[string]interface{}); ok {
				current = currentMap[element]
				if current == nil {
					return nil, nil
				}
			} else {
				return nil, nil
			}
		}
	}
	return current, nil
}

func parseArrayElement(current interface{}, element string) (interface{}, string) {
	arrayIndex := strings.Index(element, "[")
	key := element[:arrayIndex]
	indexStr := element[arrayIndex+1 : len(element)-1]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, element
	}
	if currentMap, isObject := current.(map[string]interface{}); isObject {
		if array, hasValue := currentMap[key]; hasValue {
			if index >= 0 && index < reflect.ValueOf(array).Len() {
				return reflect.ValueOf(array).Index(index).Interface(), ""
			}
		}
	}
	return nil, element
}

func getDynamoDbAttributeValueFrom(value interface{}) (types.AttributeValue, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return &types.AttributeValueMemberS{Value: v.String()}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &types.AttributeValueMemberN{Value: strconv.FormatInt(v.Int(), 10)}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &types.AttributeValueMemberN{Value: strconv.FormatUint(v.Uint(), 10)}, nil
	case reflect.Float32, reflect.Float64:
		return &types.AttributeValueMemberN{Value: strconv.FormatFloat(v.Float(), 'f', -1, 64)}, nil
	case reflect.Bool:
		return &types.AttributeValueMemberBOOL{Value: v.Bool()}, nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return &types.AttributeValueMemberB{Value: v.Bytes()}, nil
		}
		values := make([]types.AttributeValue, v.Len())
		for i := 0; i < v.Len(); i++ {
			t, err := getDynamoDbAttributeValueFrom(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			values[i] = t
		}
		return &types.AttributeValueMemberL{Value: values}, nil
	case reflect.Map:
		m := make(map[string]types.AttributeValue)
		for _, key := range v.MapKeys() {
			t, err := getDynamoDbAttributeValueFrom(v.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}
			m[key.String()] = t
		}
		return &types.AttributeValueMemberM{Value: m}, nil
	case reflect.Ptr:
		if v.IsNil() {
			return &types.AttributeValueMemberNULL{Value: true}, nil
		}
		return getDynamoDbAttributeValueFrom(v.Elem().Interface())
	default:
		return &types.AttributeValueMemberNULL{Value: true}, nil
	}
}

func getGolangValueFrom(attr types.AttributeValue) (interface{}, error) {
	switch v := attr.(type) {
	case *types.AttributeValueMemberSS, *types.AttributeValueMemberBS:
		return v.(*types.AttributeValueMemberSS).Value, nil
	case *types.AttributeValueMemberNS:
		return getGolangFloat64ArrayFrom(v.Value)
	case *types.AttributeValueMemberS:
		return v.Value, nil
	case *types.AttributeValueMemberN:
		return getGolangFloat64From(v.Value)
	case *types.AttributeValueMemberBOOL:
		return v.Value, nil
	case *types.AttributeValueMemberB:
		return v.Value, nil
	case *types.AttributeValueMemberL:
		return getGolangArrayFrom(v.Value)
	case *types.AttributeValueMemberM:
		return getGolangMapFrom(v.Value)
	case *types.AttributeValueMemberNULL:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported attribute value type: %T", attr)
	}
}

func getGolangFloat64From(value string) (interface{}, error) {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse number: %v", err)
	}
	return num, nil
}

func getGolangFloat64ArrayFrom(values []string) (interface{}, error) {
	seq := make([]float64, len(values))
	for index, value := range values {
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %v", err)
		}
		seq[index] = num
	}
	return seq, nil
}

func getGolangArrayFrom(values []types.AttributeValue) (interface{}, error) {
	parsed := make([]interface{}, len(values))
	for i, item := range values {
		t, err := getGolangValueFrom(item)
		if err != nil {
			return nil, err
		}
		parsed[i] = t
	}
	return parsed, nil
}

func getGolangMapFrom(value map[string]types.AttributeValue) (interface{}, error) {
	m := make(map[string]interface{})
	for key, item := range value {
		t, err := getGolangValueFrom(item)
		if err != nil {
			return nil, err
		}
		m[key] = t
	}
	return m, nil
}

package dynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"strconv"
	"strings"
)

type GetItemResponse struct {
	properties    map[string]any
	sdkAttributes map[string]types.AttributeValue
}

func (instance *GetItemResponse) ValueFrom(location string) (any, error) {
	if instance.properties == nil {
		value, err := getGolangValueFrom(&types.AttributeValueMemberM{Value: instance.sdkAttributes})
		if err != nil {
			return nil, err
		}
		instance.properties = value.(map[string]any)
	}

	return getValueFromPath(instance.properties, location)
}

func getValueFromPath(data map[string]any, path string) (any, error) {
	elements := strings.Split(path, ".")
	var current interface{} = data

	for _, element := range elements {
		if currentMap, ok := current.(map[string]interface{}); ok {
			if strings.Contains(element, "[") && strings.Contains(element, "]") {
				// Handle arrays
				arrayIndex := strings.Index(element, "[")
				key := element[:arrayIndex]
				indexStr := element[arrayIndex+1 : len(element)-1]
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					return nil, err
				}
				array, ok := currentMap[key].([]interface{})
				if !ok {
					return nil, nil
				}
				if index >= len(array) {
					return nil, fmt.Errorf("index out of range for %s[%d]", key, index)
				}
				current = array[index]
			} else {
				current = currentMap[element]
			}
		} else if currentArray, ok := current.([]interface{}); ok {
			index, err := strconv.Atoi(element)
			if err != nil {
				return nil, fmt.Errorf("expected an integer index, got %s", element)
			}

			if index >= len(currentArray) {
				return nil, fmt.Errorf("index out of range for array")
			}
			current = currentArray[index]
		} else {
			return nil, nil
		}
	}
	return current, nil
}

func getDynamoDbAttributeValueFrom(value any) (types.AttributeValue, error) {
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

func getGolangValueFrom(attr types.AttributeValue) (any, error) {
	switch v := attr.(type) {
	case *types.AttributeValueMemberSS:
		return v.Value, nil
	case *types.AttributeValueMemberNS:
		seq := make([]float64, len(v.Value))
		for index, value := range v.Value {
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse number: %v", err))
			}
			seq[index] = num
		}
		return seq, nil
	case *types.AttributeValueMemberBS:
		return v.Value, nil
	case *types.AttributeValueMemberS:
		return v.Value, nil
	case *types.AttributeValueMemberN:
		num, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse number: %v", err))
		}
		return num, nil
	case *types.AttributeValueMemberBOOL:
		return v.Value, nil
	case *types.AttributeValueMemberB:
		return v.Value, nil
	case *types.AttributeValueMemberL:
		values := make([]interface{}, len(v.Value))
		for i, item := range v.Value {
			t, err := getGolangValueFrom(item)
			if err != nil {
				return nil, err
			}
			values[i] = t
		}
		return values, nil
	case *types.AttributeValueMemberM:
		m := make(map[string]interface{})
		for key, item := range v.Value {
			t, err := getGolangValueFrom(item)
			if err != nil {
				return nil, err
			}
			m[key] = t
		}
		return m, nil
	case *types.AttributeValueMemberNULL:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported attribute value type: %T", attr)
	}
}

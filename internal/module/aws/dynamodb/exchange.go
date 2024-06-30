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
		instance.properties = make(map[string]any)
	}
	value, hasValue := instance.properties[location]
	if hasValue {
		return value, nil
	}
	key, subPath := instance.getKeyAndSubPath(location)
	var v any
	var err error
	if arrayIndex := instance.parseArrayIndex(key); arrayIndex != -1 {
		t := key[:strings.Index(key, "[")]
		v, err = instance.getArrayValue(instance.sdkAttributes[t], arrayIndex, subPath)
	} else {
		v, err = getGolangValueFrom(instance.sdkAttributes[key])
	}
	if err != nil {
		return nil, err
	}
	subPathValue, err := instance.getSubPathValue(v, subPath)
	if err == nil {
		instance.properties[location] = subPathValue
	}
	return subPathValue, err
}

func (instance *GetItemResponse) getSubPathValue(value any, location string) (any, error) {
	if location == "" {
		return value, nil
	}
	key, subPath := instance.getKeyAndSubPath(location)
	if arrayIndex := instance.parseArrayIndex(key); arrayIndex != -1 {
		return instance.getArrayValue(value, arrayIndex, subPath)
	}
	return instance.getObjectValue(value, key, subPath)
}

func (instance *GetItemResponse) getObjectValue(value any, key string, subPath string) (any, error) {
	m, isMap := value.(map[string]any)
	if !isMap {
		return nil, fmt.Errorf("cannot access nested value: %v is not a map", value)
	}
	subPathValue, hasValue := m[key]
	if !hasValue {
		return nil, fmt.Errorf("key %s does not exist in map", key)
	}
	return instance.getSubPathValue(subPathValue, subPath)
}

func (instance *GetItemResponse) getArrayValue(value any, arrayIndex int, subPath string) (any, error) {
	arr, isArray := value.([]any)
	if !isArray {
		return nil, fmt.Errorf("cannot access array value: %v is not an array", value)
	}
	if arrayIndex >= len(arr) {
		return nil, fmt.Errorf("index %d out of bounds for array", arrayIndex)
	}
	return instance.getSubPathValue(arr[arrayIndex], subPath)
}

func (instance *GetItemResponse) getKeyAndSubPath(location string) (string, string) {
	indexOf := strings.Index(location, ".")
	if indexOf == -1 {
		return location, ""
	}
	return location[:indexOf], location[indexOf+1:]
}

func (instance *GetItemResponse) parseArrayIndex(key string) int {
	if strings.HasSuffix(key, "]") {
		start := strings.Index(key, "[")
		if start != -1 {
			index, err := strconv.Atoi(key[start+1 : len(key)-1])
			if err == nil {
				return index
			}
		}
	}
	return -1
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

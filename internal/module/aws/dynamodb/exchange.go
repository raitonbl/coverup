package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"strconv"
)

type GetItemResponse struct {
	properties    map[string]any
	sdkAttributes map[string]types.AttributeValue
}

func (instance *GetItemResponse) ValueFrom(x string) (any, error) {
	//TODO implement me
	panic("implement me")
}

func createDynamoDbAttributeFrom(value any) (types.AttributeValue, error) {
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
			t, err := createDynamoDbAttributeFrom(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			values[i] = t
		}
		return &types.AttributeValueMemberL{Value: values}, nil
	case reflect.Map:
		m := make(map[string]types.AttributeValue)
		for _, key := range v.MapKeys() {
			t, err := createDynamoDbAttributeFrom(v.MapIndex(key).Interface())
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
		return createDynamoDbAttributeFrom(v.Elem().Interface())
	default:
		return &types.AttributeValueMemberNULL{Value: true}, nil
	}
}

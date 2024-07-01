package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
)

func TestValueFrom(t *testing.T) {
	response := GetItemResponse{
		sdkAttributes: map[string]types.AttributeValue{
			"id":           &types.AttributeValueMemberS{Value: "96cd487b-b225-4074-9bfb-85b857b016db"},
			"name":         &types.AttributeValueMemberS{Value: "RaitonBL"},
			"age":          &types.AttributeValueMemberN{Value: "21"},
			"tags":         &types.AttributeValueMemberSS{Value: []string{"Platform=Golang", "Project=Coverup"}},
			"luck_numbers": &types.AttributeValueMemberNS{Value: []string{"7", "11"}},
			"is_developer": &types.AttributeValueMemberBOOL{Value: true},
			"cli_opts": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "--alertsPerDay"},
				&types.AttributeValueMemberN{Value: "1"},
			}},
			"console_opts": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"template":    &types.AttributeValueMemberS{Value: "default"},
				"screen_mode": &types.AttributeValueMemberS{Value: "dark"},
				"zoom_in":     &types.AttributeValueMemberN{Value: "75"},
			},
			},
		},
	}

	tests := []struct {
		location string
		expected any
	}{
		{"id", "96cd487b-b225-4074-9bfb-85b857b016db"},
		{"name", "RaitonBL"},
		{"age", 21.0},
		{"tags[0]", "Platform=Golang"},
		{"tags[1]", "Project=Coverup"},
		{"luck_numbers[0]", 7.0},
		{"luck_numbers[1]", 11.0},
		{"is_developer", true},
		{"cli_opts[0]", "--alertsPerDay"},
		{"cli_opts[1]", 1.0},
		{"console_opts.template", "default"},
		{"console_opts.screen_mode", "dark"},
		{"console_opts.zoom_in", 75.0},
		{"tags[0]", "Platform=Golang"},
	}

	for _, tt := range tests {
		t.Run(tt.location, func(t *testing.T) {
			value, err := response.ValueFrom(tt.location)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestInvalidLocation(t *testing.T) {
	response := GetItemResponse{
		sdkAttributes: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "96cd487b-b225-4074-9bfb-85b857b016db"},
		},
	}
	tests := []struct {
		location string
	}{
		{"unknown"},
		{"tags[2]"},
		{"console_opts.invalid_key"},
	}
	for _, tt := range tests {
		t.Run(tt.location, func(t *testing.T) {
			valueOf, err := response.ValueFrom(tt.location)
			assert.Nil(t, err)
			assert.Nil(t, valueOf)
		})
	}
}

func TestGetDynamoDbAttributeValueFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected types.AttributeValue
		hasError bool
	}{
		{
			name:     "String",
			input:    "test",
			expected: &types.AttributeValueMemberS{Value: "test"},
			hasError: false,
		},
		{
			name:     "Int",
			input:    123,
			expected: &types.AttributeValueMemberN{Value: "123"},
			hasError: false,
		},
		{
			name:     "Float",
			input:    123.456,
			expected: &types.AttributeValueMemberN{Value: strconv.FormatFloat(123.456, 'f', -1, 64)},
			hasError: false,
		},
		{
			name:     "Bool",
			input:    true,
			expected: &types.AttributeValueMemberBOOL{Value: true},
			hasError: false,
		},
		{
			name:     "ByteSlice",
			input:    []byte{1, 2, 3},
			expected: &types.AttributeValueMemberB{Value: []byte{1, 2, 3}},
			hasError: false,
		},
		{
			name:  "Slice",
			input: []interface{}{"a", 1, true},
			expected: &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "a"},
				&types.AttributeValueMemberN{Value: "1"},
				&types.AttributeValueMemberBOOL{Value: true},
			}},
			hasError: false,
		},
		{
			name:  "Map",
			input: map[string]interface{}{"a": "1", "b": 2, "c": false},
			expected: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"a": &types.AttributeValueMemberS{Value: "1"},
				"b": &types.AttributeValueMemberN{Value: "2"},
				"c": &types.AttributeValueMemberBOOL{Value: false},
			}},
			hasError: false,
		},
		{
			name:     "Pointer",
			input:    func() *int { i := 123; return &i }(),
			expected: &types.AttributeValueMemberN{Value: "123"},
			hasError: false,
		},
		{
			name:     "NilPointer",
			input:    (*int)(nil),
			expected: &types.AttributeValueMemberNULL{Value: true},
			hasError: false,
		},
		{
			name:     "UnsupportedType",
			input:    struct{}{},
			expected: &types.AttributeValueMemberNULL{Value: true},
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getDynamoDbAttributeValueFrom(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

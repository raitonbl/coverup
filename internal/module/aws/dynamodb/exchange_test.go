package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
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

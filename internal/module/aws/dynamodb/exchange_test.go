package dynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
)

func TestGetItemResponse_ValueFrom(t *testing.T) {
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
			}},
		},
	}
	valueOf, err := response.ValueFrom("id")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(valueOf)
}

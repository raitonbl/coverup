package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/raitonbl/coverup/internal/module/aws/config"
	"github.com/raitonbl/coverup/pkg/api"
	"strings"
)

const (
	ComponentType = "DynamoDbItem"
)

type StepFactory struct {
}

func (instance *StepFactory) New(ctx api.StepDefinitionContext) {
	variantOpts := []string{"", fmt.Sprintf(`, known as %s`, api.NonLiteralStringExpression)}
	opts := []string{api.ValueExpression, fmt.Sprintf(`"%s"`, api.ValueExpression), api.LiteralStringExpression}
	description := "Checks if an item with Given definition exists in a DynamoDb Table"
	for _, variant := range variantOpts {
		isGetOperation := variant == variantOpts[1]
		if isGetOperation {
			description = "Gets a  DynamoDb Item, with specified primary key, from a DynamoDb Table"
		}
		step := api.StepDefinition{
			Description: description,
			Options:     make([]api.Option, 0),
		}
		for _, dynamoDbTableOpt := range opts {
			step.Options = append(step.Options, api.Option{
				Description:    description,
				Regexp:         fmt.Sprintf(` DynamoDb Table %s has item(([^)]+))%s`, dynamoDbTableOpt, variant),
				HandlerFactory: instance.createDynamoDbOperationHandlerFactory(isGetOperation),
			})
		}
		ctx.Step(step)
	}
}

func (instance *StepFactory) createDynamoDbOperationHandlerFactory(isGetOperation bool) api.HandlerFactory {
	return func(c api.ScenarioContext) any {
		if isGetOperation {
			return func(dynamoDbTable, itemDefinition, alias string) error {
				return instance.doDynamoDbGetItem(c, dynamoDbTable, itemDefinition, alias)
			}
		}
		return func(dynamoDbTable, itemDefinition string) error {
			return instance.doDynamoDbGetItem(c, dynamoDbTable, itemDefinition, "")
		}
	}
}

func (instance *StepFactory) doDynamoDbGetItem(c api.ScenarioContext, name, definition, alias string) error {
	clientRegistry, err := c.GetGivenComponent(config.ClientRegistryComponentType, "")
	if err != nil {
		return err
	}
	t, err := c.Resolve(name)
	if err != nil {
		return err
	}
	dynamoDbTable, isString := t.(string)
	if !isString {
		return fmt.Errorf("no such dynamoDb Table %v", t)
	}
	dynamoDbClient, err := clientRegistry.(config.ClientRegistry).GetClient("", dynamoDbTable)
	if err != nil {
		return err
	}
	t, err = c.Resolve(definition)
	if err != nil {
		return err
	}
	itemDefinition, isString := t.(string)
	if !isString {
		return fmt.Errorf("no such Item (%v) on dynamoDb Table %s", itemDefinition, dynamoDbTable)
	}
	request, err := instance.createGetItemRequest(c, dynamoDbTable, itemDefinition)
	if err != nil {
		return err
	}
	response, err := dynamoDbClient.(*dynamodb.Client).GetItem(context.Background(), request)
	if err != nil {
		return err
	}
	if response.Item == nil {
		return fmt.Errorf("no such Item (%s) on dynamoDb Table %s", itemDefinition, dynamoDbTable)
	}
	if alias == "" {
		return nil
	}
	return c.AddGivenComponent(ComponentType, &GetItemResponse{sdkAttributes: response.Item}, alias)
}

func (instance *StepFactory) createGetItemRequest(c api.ScenarioContext, table string, filterDefinition string) (*dynamodb.GetItemInput, error) {
	attributeDefinitions := strings.Split(filterDefinition, ",")
	request := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       make(map[string]types.AttributeValue),
	}
	if len(attributeDefinitions) < 1 || len(attributeDefinitions) > 2 {
		return nil, fmt.Errorf("unsupported format for attribute definition (%s)", filterDefinition)
	}
	for _, each := range attributeDefinitions {
		definition := strings.TrimLeft(each, " ")
		splitIndex := strings.Index(definition, "=")
		name := definition[:splitIndex]
		expr := strings.TrimLeft(definition[splitIndex+1:], " ")
		rawValue, err := c.Resolve(expr)
		if err != nil {
			return nil, fmt.Errorf(`error caught during %s mapping into dynamoDbAttribute. caused by:\n%v`, name, err)
		}
		value, err := getDynamoDbAttributeValueFrom(rawValue)
		if err != nil {
			return nil, fmt.Errorf(`error caught during %s mapping into dynamoDbAttribute. caused by:\n%v`, name, err)
		}
		request.Key[name] = value
	}
	return request, nil
}

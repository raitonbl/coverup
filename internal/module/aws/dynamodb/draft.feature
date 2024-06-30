Feature: Design

  Scenario:
    Given DynamoDb {{Properties.DynamoDb.name}} item named GetItem
    And {{DynamoDbItem.GetItem}}.id is equal to "a7459242-f8cf-4291-af86-351a70bcebdb"
    Then {{DynamoDbItem.GetItem}}.name is equal to "Gamdias Case"

  Scenario:
    Given DynamoDb {{Properties.DynamoDb.name}} item named GetItem
    And {{DynamoDbItem.GetItem}}.id is equal to "a7459242-f8cf-4291-af86-351a70bcebdb"
    Then {{DynamoDbItem.GetItem}}.name is equal to "Gamdias Case"

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And method is POST
    And path is /vouchers
    And body is:
          """
                  {
                      "benefit": "PSN 100 UK",
                      "promo-code": "raitonbl.com"
                  }
          """
    And DynamoDb Table  {{Properties.DynamoDb.name}} has item(id= "ac522d71-13a4-4287-9816-7f5c4b21b54d" )
    And DynamoDb Table  {{Properties.DynamoDb.name}} has item(id= "ac522d71-13a4-4287-9816-7f5c4b21b54d" ), known as GetItem
    And  the http response body $.benefit is equal to {{DynamoDbItem.GetItem.benefit}}

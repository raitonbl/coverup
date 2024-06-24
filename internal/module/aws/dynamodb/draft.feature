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
    Then Fetch item GetItem from DynamodbTable {{Properties.DynamoDb.name}}
    Then  the http response body $.benefit is equal to {{DynamoDbItem.GetItem.benefit}}


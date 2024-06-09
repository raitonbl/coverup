Feature: Design

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And the Method is GET
    And server url is https://localhost:8443
    And path is /persons/27258303-9ebc-4b84-a17e-f886161ab2f5
    When the client submits the HttpRequest
    Then the response status code is 200
    And the response content-type is "application/json"
    And the response headers contains:
      | content-type      | application/json                                                             |
      | x-amazon-trace-id | Root=1-5a969f52-0b48d5e712d3d3a6b1c8ad89; Parent=53995c3f42cd8ad8; Sampled=1 |
    And the response body respects schema file://schemas/person.json
    And the response body $.name is "John"
    And the response body $.age is 30
    And the response body $.is_eligible is true
    And the response body $.id is "27258303-9ebc-4b84-a17e-f886161ab2f5"

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And the method is GET
    And path is /persons/{{Properties.entities.default.id}}
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the response content-type is "application/json"
    And the response headers contains:
      | content-type      | application/json                                                             |
      | x-amazon-trace-id | Root=1-5a969f52-0b48d5e712d3d3a6b1c8ad89; Parent=53995c3f42cd8ad8; Sampled=1 |
    And the response body respects schema file://schemas/person.json
    And the response body $.name is "John"
    And the response body $.age is 30
    And the response body $.is_eligible is true
    And the response body $.id is "{{Properties.entities.default.id}}"

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And method is GET
    And path is /persons/{{Properties.entities.default.id}}
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the response content-type is "application/json"
    And the response body respects schema file://schemas/person.json
    And the body is file://person.response.json

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And method is GET
    And path is /persons/{{Properties.entities.default.id}}
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the response content-type is "application/json"
    And the response body respects schema file://schemas/person.json
    And the body is file://{{Properties.File.Response}}

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And method is GET
    And path is /persons/{{Properties.entities.default.id}}
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the response content-type is "application/json"
    And the response body respects schema file://schemas/person.json
    And the response body is:
    """
    {
      "id": "27258303-9ebc-4b84-a17e-f886161ab2f5",
      "name": "John",
      "age": 30,
      "is_eligible": true
    }
    """

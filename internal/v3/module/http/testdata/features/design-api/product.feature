Feature: Design

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | "application/json" |
    And method is GET
    And path is /products/27258303-9ebc-4b84-a17e-f886161ab2f5
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the content-type is "application/json"
    And the response body $.id is "27258303-9ebc-4b84-a17e-f886161ab2f5"
    And response body $.id starts with "27258303"
    And response body $.id ends with "f886161ab2f5"
    And response body $.id contains "4b84-a17e"
    And response body $.id ignoring case starts with "27258303"
    And response body $.id ignoring case ends with "f886161ab2f5"
    And response body $.id ignoring case contains "4b84-a17e"
    And response body $.id matches pattern "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | "application/json" |
    And method is GET
    And path is /products/27258303-9ebc-4b84-a17e-f886161ab2f5
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the content-type is "application/json"
    And the response body $.offer_created_at is DateTime
    And the response body $.offer_expires_at is before {{Properties.clock.endOfLife}}
    And the response body $.offer_expires_at is after {{HttpRequest.Current.Body.offer_created_at}}
    And the response body $.offer_created_at is before {{HttpRequest.Current.Body.offer_expires_at}}

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | "application/json" |
    And method is GET
    And path is /products/27258303-9ebc-4b84-a17e-f886161ab2f5
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the content-type is "application/json"
    And response body $.about length is 3
    And response body $.about[0] starts with "One Touch SSD"
    And response body $.tags[0].name is "IT"
    And response body $.tags[*].name is "IT"

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | "application/json" |
    And method is GET
    And path is /products/27258303-9ebc-4b84-a17e-f886161ab2f5
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the content-type is "application/json"
    And response body $.warranty.amount is greater than 1
    And response body $.warranty.unit is part of ["years","months","days"]
    And the response body $.price.amount is greater or equal to 200
    And the response Body $.characteristics.capacity.amount is lesser or equal to 1
    And response body $.characteristics.hard_disk_form_factor.amount is greater than 2.6


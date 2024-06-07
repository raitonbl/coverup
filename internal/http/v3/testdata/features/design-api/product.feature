Feature: Design

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And method is GET
    And path is /products/27258303-9ebc-4b84-a17e-f886161ab2f5
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the content-type is application/json
    And $body.id is "27258303-9ebc-4b84-a17e-f886161ab2f5"
    And $body.id starts with "27258303"
    And $body.id ends with "f886161ab2f5"
    And $body.id contains with "4b84-a17e"
    And $body.id starts with "27258303", ignoring case
    And $body.id ends with "f886161ab2f5", ignoring case
    And $body.id contains "4b84-a17e", ignoring case
    And $body.id matches pattern "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And method is GET
    And path is /products/27258303-9ebc-4b84-a17e-f886161ab2f5
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the content-type is application/json
    And the $body.offer_created_at is DateTime
    And the $body.offer_created_at is before $body.offer_expires_at
    And the $body.offer_expires_at is after $body.offer_created_at
    And the $body.offer_expires_at is before {{Properties.clock.endOfLife}}

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And method is GET
    And path is /products/27258303-9ebc-4b84-a17e-f886161ab2f5
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the content-type is application/json
    And the $body.about length is 3
    And the $body.about is not empty
    And the $body.about[0] starts with "One Touch SSD"
    And the $body.tags[0].name is "IT"
    And the $body.tags[*].name is "IT"

  Scenario:
    Given a HttpRequest
    And the headers:
      | content-type | application/json |
    And method is GET
    And path is /products/27258303-9ebc-4b84-a17e-f886161ab2f5
    And server url is https://localhost:8443
    When the client submits the HttpRequest
    Then the response status code is 200
    And the content-type is application/json
    And the $body.warranty.amount is greater than 1
    And the $body.warranty.unit is part of ["years","months","days"}
    And the $body.price.amount is greater of equal to 200
    And the $body.characteristics.capacity.amount is lesser or equal to 1
    And the $body.characteristics.hard_disk_form_factor.amount than 2.6


Feature: Design

  Scenario:
    Given a HttpRequest
      And the headers:
        | content-type | application/json |
      And method is GET
      And path is /vouchers/27258303-9ebc-4b84-a17e-f886161ab2f5
      And server url is https://localhost:8443
      When the client submits the HttpRequest
    Then the response status code is 200

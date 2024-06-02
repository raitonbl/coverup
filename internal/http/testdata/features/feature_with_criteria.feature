Feature: GetBreads
    Scenario: No query parameters
        Given a HttpRequest
            And the headers:
                | content-type  | application/json |
        When GET "/breeds/f9643a80-af1d-422a-9f15-18d466822053"
        Then the response statusCode is 200
            And the $body.data.id is equal to "f9643a80-af1d-422a-9f15-18d466822053"
            And the $body.data.type is equal to "breed"
            And the $body.data.name starts with "Caucasian"
            And the $body.data.name ends with "Dog"
            And the $body.data.attributes.min_life is greater or equal to 15
            And the $body.data.attributes.max_life is lesser or equal to 20
            And the $body.data.attributes.description matches "^The$"
            And the $headers.content-type starts with "application/"
            And the $headers.content-type ends with "/json"
            And the $headers.content-type is equal to "application/json"
            And the $headers.x-ratelimit-limit is greater or equal to 1
            And the $headers.x-ratelimit-remaining is lesser or equal to 1625074800

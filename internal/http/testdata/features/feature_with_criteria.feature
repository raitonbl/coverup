Feature: GetBreads
    Scenario: No query parameters
        Given a HttpRequest
            And the headers:
                | content-type  | application/json |
        When GET "/breeds/f9643a80-af1d-422a-9f15-18d466822053"
        Then the response statusCode is 200
            And the response headers:
                | content-type  | application/json |
            And the response body:
                """
                {
                   "data":{
                      "id":"f9643a80-af1d-422a-9f15-18d466822053",
                      "type":"breed",
                      "attributes":{
                         "name":"Caucasian Shepherd Dog",
                         "min_life":15,
                         "max_life":20,
                         "description":"The Caucasian Shepherd dog is a serious guardian breed and should never be taken lightly.",
                         "hypoallergenic":false
                      }
                   },
                   "links":{
                      "self":"https://dogapi.dog/api/v2/breeds/f9643a80-af1d-422a-9f15-18d466822053"
                   }
                }
                """
            And the $.data.id is equal to "f9643a80-af1d-422a-9f15-18d466822053"
            And the $.data.type is equal to "breed"
            And the $.data.name starts with "Caucasian"
            And the $.data.name ends with "Dog"
            And the $.data.attributes.min_life is greater or equal to 15
            And the $.data.attributes.max_life is greater or equal to 20
            And the $.data.attributes.description matches "^The$"

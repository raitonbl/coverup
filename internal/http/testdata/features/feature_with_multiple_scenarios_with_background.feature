Feature: GetBreeds

    Background: On Get Breeds
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        When GET "/breeds"

    Scenario: Has response body
        Then the response statusCode is 200
            And the response headers:
                | content-type  | application/json |
            And the response body:
                """
                {
                   "data":[
                      {
                         "id":"f9643a80-af1d-422a-9f15-18d466822053",
                         "type":"breed",
                         "attributes":{
                            "name":"Caucasian Shepherd Dog",
                            "description":"The Caucasian Shepherd dog is a serious guardian breed and should never be taken lightly. ",
                            "hypoallergenic":false
                         }
                      },
                      {
                         "id":"dc5e84f8-9151-4624-836c-25b4e313118b",
                         "type":"breed",
                         "attributes":{
                            "name":"Bouvier des Flandres",
                            "description":"They don't build 'em like this anymore.",
                            "hypoallergenic":false
                         }
                      }
                   ],
                   "links":{
                      "self":"https://dogapi.dog/api/v2/breeds",
                      "current":"https://dogapi.dog/api/v2/breeds?page[number]=1",
                      "next":"https://dogapi.dog/api/v2/breeds?page[number]=2",
                      "last":"https://dogapi.dog/api/v2/breeds?page[number]=2"
                   }
                }
                """

    Scenario: Has response body uri(File)
        Then the response statusCode is 200
        And the response headers:
            | content-type  | application/json |
        And the response body uri is: file://GetBreeds.json

    Scenario: Has response body uri(http)
        Then the response statusCode is 200
        And the response headers:
            | content-type  | application/json |
        And the response body uri is: http://localhost:9999/GetBreeds.json
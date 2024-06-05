Feature: Design
    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And method is GET
        And path is /vouchers/27258303-9ebc-4b84-a17e-f886161ab2f5
        And server url is https://localhost:8443

    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And method is POST
        And path is /vouchers
        And body is:
                """
                {
                    "benefit": "PSN 100 UK",
                    "promo-code": "raitonbl.com"
                }
                """

    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And the method is POST
        And path is /vouchers
        And server url is https://localhost:8443
        And accept is "application/json"

    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And method is POST
        And path is /vouchers
        And server url is https://localhost:8443
        And content-type is "application/json"

    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And method is POST
        And path is https://localhost:8443/vouchers
        And content-type is "application/json"
        And the header x-amazon-trace-id is "Root=1-5a969f52-0b48d5e712d3d3a6b1c8ad89; Parent=53995c3f42cd8ad8; Sampled=1"

    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And method is POST
        And path is /vouchers
        And server url is https://localhost:8443
        And content-type is "application/json"
        And body is file://request.json

    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And method is POST
        And path is /forms
        And server url is https://localhost:8443
        And form enctype is multipart/form-data
        And form attribute "full_name" is "RaitonBL"
        And form attribute "full_name" is "{{Entities.administrator.name}}"
        And form attribute "picture" is file://image.png
        And form attribute "picture" is file://{{Properties.files.picture}}

    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And method is POST
        And path is /forms
        And server url is https://localhost:8443
        And form enctype is multipart/form-data
        And form attribute "full_name" is "{{Entities.administrator.name}}"
        And form attribute "picture" is file://image.png
        And form attribute "picture" is file://{{Properties.files.picture}}

    Scenario:
        Given a HttpRequest
        And the headers:
            | content-type  | application/json |
        And method is POST
        And path is /forms
        And server url is https://localhost:8443
        And form method is POST
        And form enctype is application/x-www-form-urlencoded
        And form enctype is multipart/form-data
        And form attribute "full_name" is "RaitonBL"
        And form attribute "full_name" is "{{Entities.administrator.name}}"

Feature: Buy voucher
    Scenario: Buy a PSN 100 UK voucher
        Given a HttpRequest named SendVoucherRequest
            And the headers:
                | content-type  | application/json |
            And operation POST /vouchers
            # And Timeout 3 seconds
            And body:
                """
                {
                    "benefit": "PSN 100 UK",
                    "promo-code": "raitonbl.com"
                }
                """
        Then the response statusCode is 200
            And the {{HttpRequest.SendVoucherRequest.StatusCode}} is 200
            And the $body complies with schema file://response.schema.json
            And the $body.benefit is equal to "PSN 100 UK"
            And the $body.price.amount is equal to 85
            And the $body.price.currency is equal to GBP
            And the $body.has_discount is true
            And the {{HttpRequest.SendVoucherRequest.body}}.id is defined
            # Force a 2 second wait
            And wait 2 seconds

        Given a HttpRequest named AssertVoucherHasBeenPurchased
            And Server https://www.api.psn.co.uk
            And the headers:
                | content-type  | application/json |
            And Operation GET /vouchers/{{HttpRequest.SendVoucherRequest.Response.Body.id}}
        Then the response statusCode is 200
            And the $body complies with schema file://psn-response-schema
            And the $body.id is equal to {{HttpRequest.SendVoucherRequest.Response.Body.id}}

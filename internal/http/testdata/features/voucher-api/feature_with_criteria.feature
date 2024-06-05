Feature: Buy voucher
    Scenario: Buy a PSN 100 UK voucher
        Given a HttpRequest named SendVoucherRequest
            And the headers:
                | content-type  | application/json |
            And method is POST
            And uri is /vouchers
           # And server url is https://localhost:8443
            And body is:
                """
                {
                    "benefit": "PSN 100 UK",
                    "promo-code": "raitonbl.com"
                }
                """

           # And accept is "application/json"
           # And content-type is "application/json"

           # And body is file://request.json
           # And body is http://request.json
           # And body is htts://request.json

           # And form method is POST
           # And form enctype is multipart/form-data
           # And form enctype is multipart/form-data
           # And form field "full_name" is "RaitonBL"
           # And form field "full_name" is "{{Entities.administrator.name}}"
           # And form field "picture" is file://image.png
           # And form field "picture" is file://{{Properties.files.picture}}

           # And form method is POST
           # And form enctype is application/x-www-form-urlencoded
           # And form enctype is multipart/form-data
           # And form field "full_name" is "RaitonBL"
           # And form field "full_name" is "{{Entities.administrator.name}}"

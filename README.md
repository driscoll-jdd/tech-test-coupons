# tech-test-coupons
This is a simple micro service I was asked to create as a technical exercise as part of a job application.

## Install

    go get github.com/driscoll-jdd/tech-test-coupons

## Usage

The main application is contained within the file TechTest.go. Running (or compiling and executing) this will open a web server on port 8080, which can be used to interact with the microservice.

## Endpoints

The service contains the following endpoints:

    /api/coupons/create

This accepts a json POST body containing a slice of Coupon (Structures/Coupon.go) structs. Coupons are uniquely identified and stored by the id field (json: id, struct: UUID) so you can edit a coupon by populating this field with an existing value.

In response to calls to this endpoint, you get an instance of the CreationResponse struct (Structures/CreationResponse.go) which will provide information on the success of your request, the number of coupons created or updated, and a slice of the arrays which have been manipulated - including UUID values for coupons which have been created, which you can use later to edit these coupons.

HTTP status codes are as follows:

 * 200 - This is for a successful operation
 * 400 - This indicates a bad request where either no json content was received or no valid coupons were received

Coupons can be checked before submission with an inbuilt method IsValid() which will simply validate a coupon and return a boolean indicator of validity.

    /api/coupons/fetch

This returns a slice of coupons, either all of them by default or a filtered list if you supply GET parameters. Currently supported are:

 * name
 * brand

package main

import (
	"fmt"
	"net/http"
	"TechTest/Structures"
	"io/ioutil"
	"encoding/json"
	"strings"
)

var cache Structures.Cache

	func main() {

		// Set up our internal cache - a little library I wrote which implements similar functionality to Redis, but within go
		// By default cache coupons for an hour
		cache = Structures.CreateCache(3600)

		// This method will create new coupons, but if a UUID is supplied in the json, coupons will be updated instead
		http.HandleFunc("/api/coupons/create", UpdateCoupons)

		// Fetch coupons
		http.HandleFunc("/api/coupons/fetch", ListCoupons)

		// Open a basic web server
		err := http.ListenAndServe(":8080", nil)

		if(err != nil) {

			fmt.Println("Error opening socket: " + err.Error())
		}
	}

	// Send back a json encoded array of tokens - all if no filters are applied, or else only those which match the given filters
	func ListCoupons(w http.ResponseWriter, r *http.Request) {

		// Start by getting all the coupons from cache
		couponsJson := cache.ReadAll()
		coupons := make([]Structures.Coupon, 0)

		// Unpack our cached json into structs
		for _, rawCoupon := range couponsJson {

			myCoupon := Structures.Coupon{}
			json.Unmarshal([]byte(rawCoupon), &myCoupon)

			coupons = append(coupons, myCoupon)
		}

		/*
			This becomes more complicated when you have potentially thousands of coupons in a database, and a percentage of those in cache. For that scenario it is necessary to hand this over to a database here
			instead of purely relying on RAM which this proof of concept can get away with
		*/

		// At this point coupons is an array of coupon structs. Apply filters if the user has defined them
		if(len(r.URL.Query().Get("brand")) > 0) {

			for id, coupon := range coupons {

				if(strings.ToLower(coupon.Brand) != strings.ToLower(r.URL.Query().Get("brand"))) {

					// Cut this specific coupon out of our results basket
					if(len(coupons) > id) {

						coupons = append(coupons[:id], coupons[id + 1:]...)

					} else {

						coupons = append(coupons[:id])
					}
				}
			}
		}

		if(len(r.URL.Query().Get("name")) > 0) {

			for id, coupon := range coupons {

				if(strings.ToLower(coupon.Name) != strings.ToLower(r.URL.Query().Get("name"))) {

					// Cut this specific coupon out of our results basket
					coupons = append(coupons[:id], coupons[id + 1:]...)
				}
			}
		}

		// Convert our results basket into json
		resultsJson, _ := json.Marshal(coupons)

		// Dump this out
		w.WriteHeader(200)
		w.Write(resultsJson)
	}

	// Create or update an array of coupons, supplied in json format as the raw POST content for this POST request
	func UpdateCoupons(w http.ResponseWriter, r *http.Request) {

		// Get our POST content - should be a json block of coupons
		jsonContent, err := ioutil.ReadAll(r.Body)

		// Create a response struct here
		response := Structures.CreationResponse{}

		if(err != nil) {

			// By default everything in the response is blank and the outcome is false, which is perfect. But we should give a message about what happened here.
			response.Message = "No POST json content supplied"

			// Turn into bytes
			responseBytes, _ := json.Marshal(response)

			// Our response
			w.WriteHeader(400)
			w.Write(responseBytes)
			return
		}

		// Make a block of coupons to create
		coupons := make([]Structures.Coupon, 0)

		// Unmarshal our json into a slice of coupons
		json.Unmarshal(jsonContent, &coupons)

		// Flip through these, adding them to our internal state
		created := 0
		updated := 0
		myCoupons := make([]Structures.Coupon, 0)
		for _, coupon := range coupons {

			if(!coupon.IsValid()) {

				// Don't bother with this one
			}

			isCreation := false

			// In this simple version of an API, assume we are updating a coupon if the coupon has a UUID. Later on, might be wise to check the cache to determine if this is an update
			if(len(coupon.UUID) < 1) {

				isCreation = true
			}

			// Store this coupon into our cache, effectively 'creating' it
			cache.Write(coupon.GetUUID(), coupon.Marshal(), 0)

			// Keep a copy of this for our response. For updates, this copy now includes the UUID end users will use to update the coupon in future.
			myCoupons = append(myCoupons, coupon)

			// Announce this coupon has been created / updated to other active instances of this microservice. This might be through something like RabbitMQ, or a custom solution with ZeroMQ, or equivalent.
			// Could be that you simply store this in Redis and use that as your high speed storage, instead of RAM. Anyway, put this into the high speed storage here
			go cacheCoupon(coupon)

			// At this point, create a background thread to store this in our database / redis instance / sqlite file / other great storage medium that is slower than RAM
			go storeCoupon(coupon)

			// Increment our creation counter
			if(isCreation) {

				created++

			} else {

				updated++
			}
		}

		// Did we create / update anything
		if(created < 1 && updated < 1) {

			// This hasn't worked out. The data supplied must have been invalid
			response.Message = "No valid coupons supplied"

			// Turn into bytes
			responseBytes, _ := json.Marshal(response)

			// Our response
			w.WriteHeader(400)
			w.Write(responseBytes)
			return
		}

		// Prepare a response for the user
		response.Success = true
		response.Created = created
		response.Updated = updated
		response.Coupons = myCoupons

		// Convert our response to json
		responseBytes, _ := json.Marshal(response)

		// The set of created / updated coupons supplied back to the end user with some information about the number of units created / updated
		w.WriteHeader(200)
		w.Write(responseBytes)
	}



/*
	================
	Helper Functions
	================
*/

	// Store this in your high speed data store, such as a messaging queue to sync with other instances, or Redis or something like that
	func cacheCoupon(coupon Structures.Coupon) {


	}

	// Store a coupon in a slow(er) data store such as a database
	func storeCoupon(coupon Structures.Coupon) {


	}
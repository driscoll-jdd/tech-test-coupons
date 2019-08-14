package main

import (
	"TechTest/Structures"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// A holder for the coupons we will use to test the api
var coupons []Structures.Coupon

	// Test the creation of coupons - should return a json array of coupons that match the ones we sent in, but with UUIDs
	func TestCreate(t *testing.T) {

		// Create some basic coupons
		coupons = append(coupons, Structures.Coupon{ Name: "Special Offer", Brand: "Driscoll", Value: 4.99, Created: time.Now(), Expiry: time.Now().Add(time.Hour * 72)})
		coupons = append(coupons, Structures.Coupon{ Name: "Exclusive Bargain", Brand: "TryNSave", Value: 19.99, Created: time.Now(), Expiry: time.Now().Add(time.Hour * 48)})

		// Get a json version of our coupons array
		couponsJson, _ := json.Marshal(coupons)

		// Try to create these coupons
		client := http.Client{}
		req, _ := http.NewRequest("POST","http://127.0.0.1:8080/api/coupons/create", bytes.NewBuffer(couponsJson))
		req.Header.Set("Content-Type","application/x-www-form-urlencoded")

		resp, err := client.Do(req)

		if(err != nil) {

			t.Fail()
			t.Logf("Error making post request for token creation")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		if(err != nil) {

			t.Fail()
			t.Log("Error reading response body from token creation")
			return
		}

		// Unpack and check our response
		response := Structures.CreationResponse{}
		json.Unmarshal(body, &response)

		if(response.Created != 2) {

			t.Fail()
			t.Log("Amount created did not equal two. Got this amount: [" + strconv.Itoa(response.Created) + "]")
			return
		}

		// Go through each of the new coupons, saving our UUIDs which we can use later to update the coupons
		for pos, newCoupon := range response.Coupons {

			if(newCoupon.Name != coupons[pos].Name || newCoupon.Brand != coupons[pos].Brand) {

				t.Fail()
				return
			}

			coupons[pos].UUID = newCoupon.UUID
		}
	}

	// Test our ability to modify a coupon
	func TestModify(t *testing.T) {

		// Take one of our coupons from before and make a modification
		myCoupon := coupons[1]
		myCoupon.Brand = "BurgerKong"
		coupons[1] = myCoupon

		// Wrap this into an array of coupons
		myCoupons := make([]Structures.Coupon, 0)
		myCoupons = append(myCoupons, myCoupon)

		// Convert this to json
		myCouponJson, _ := json.Marshal(myCoupons)

		// Send this in for modification
		client := http.Client{}
		req, _ := http.NewRequest("POST","http://127.0.0.1:8080/api/coupons/create", bytes.NewBuffer(myCouponJson))
		req.Header.Set("Content-Type","application/x-www-form-urlencoded")

		resp, err := client.Do(req)

		if(err != nil) {

			t.Fail()
			t.Logf("Error making post request for token modification")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		if(err != nil) {

			t.Fail()
			t.Log("Error reading response body from token modification")
			return
		}

		// Unpack our response into a struct
		response := Structures.CreationResponse{}
		json.Unmarshal(body, &response)

		if(response.Updated != 1) {

			t.Fail()
			t.Log("Modified coupon count was [" + strconv.Itoa(response.Updated) + "]")
			return
		}

		// Check the modified coupon is echoed back to us
		if(response.Coupons[0].Brand != myCoupon.Brand) {

			t.Fail()
			t.Log("Brand did not match")
			return
		}
	}

	// Test the ability to search for certain coupons
	func TestSearch(t *testing.T) {

		client := http.Client{}
		req, _ := http.NewRequest("GET","http://127.0.0.1:8080/api/coupons/fetch?brand=BurgerKong", nil)

		resp, err := client.Do(req)

		if(err != nil) {

			t.Fail()
			t.Logf("Error making post request for token modification")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		if(err != nil) {

			t.Fail()
			t.Log("Error reading response body from token modification")
			return
		}

		// Unpack our coupon(s) - hopefully just one
		myCoupons := make([]Structures.Coupon, 0)
		json.Unmarshal(body, &myCoupons)

		// OK so the first result should match our second coupon
		if(myCoupons[0].Brand != coupons[1].Brand) {

			t.Fail()
			t.Log("Search did not give the expected result")
		}
	}
package Structures

import (
	"time"
	"github.com/satori/go.uuid"
	"encoding/json"
)

	type Coupon struct {

		UUID string							`json:"id"`
		Name string							`json:"name"`
		Brand string						`json:"brand"`
		Value float64						`json:"value"`
		Created time.Time					`json:"createdAt"`
		Expiry time.Time					`json:"expiry"`
	}

	// Basic validation for a coupon
	func (c *Coupon) IsValid() bool {

		if(len(c.Name) > 0 && c.Value > 0 && c.Expiry.After(time.Now())) {

			return true
		}

		return false
	}

	// Fetch the UUID for this coupon - ones sent in via the API might not have a UUID, so generate one if necessary
	func (c *Coupon) GetUUID() string {

		if(len(c.UUID) < 1) {

			c.UUID = uuid.Must(uuid.NewV4()).String()
		}

		return c.UUID
	}

	// Function to automatically get the json of this coupon. Includes a check that we have a UUID, as
	// people sending in json of coupons via the API won't necessarily come up with a UUID for their proposed coupon
	func (c *Coupon) Marshal() string {

		// Check we have a UUID - calls to the API to create a coupon would not necessarily include this
		c.GetUUID()

		// Marshal ourselves and send back a string representation
		couponBytes, _ := json.Marshal(c)
		return string(couponBytes)
	}

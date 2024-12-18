package models

import "time"

// Counter represents a counter for a merchant
type Counter struct {
	ID            string    `bson:"_id,omitempty"`  // MongoDB document ID
	MerchantPhone string    `bson:"merchant_phone"` // Merchant's phone number
	Name          string    `bson:"name"`           // Counter name
	Location      string    `bson:"location"`       // Counter location
	Description   string    `bson:"description"`    // Counter description
	CreatedAt     time.Time `bson:"created_at"`     // Creation timestamp
	UpdatedAt     time.Time `bson:"updated_at"`     // Update timestamp
}

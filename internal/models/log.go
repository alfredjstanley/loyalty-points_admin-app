package models

import "time"

// Log represents a MongoDB document for logging requests and responses
type Log struct {
	ID            string    `bson:"_id,omitempty"`
	UserPhone     string    `bson:"user_phone"`
	MerchantPhone string    `bson:"merchant_phone"`
	Amount        float64   `bson:"amount"`
	InvoiceID     string    `bson:"invoice_id"`
	PaymentMode   string    `bson:"payment_mode"`
	Status        string    `bson:"status"` // "SUCCESS" or "FAILURE"
	Response      any       `bson:"response"`
	CreatedAt     time.Time `bson:"created_at"`
}

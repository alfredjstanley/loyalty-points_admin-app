package models

import "time"

// Report represents the data needed for the merchant reports
type Report struct {
	ID                string    `json:"id" bson:"_id,omitempty"`
	MerchantName      string    `json:"merchant_name" bson:"merchant_name"`
	TotalSales        float64   `json:"total_sales" bson:"total_sales"`
	TotalTransactions int       `json:"total_transactions" bson:"total_transactions"`
	PointsEarned      float64   `json:"points_earned" bson:"points_earned"` // Can be calculated as total sales / 25
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
}

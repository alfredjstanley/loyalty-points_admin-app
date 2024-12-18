package models

type Report struct {
	MerchantPhone     string  `bson:"_id" json:"phone_number"`
	TotalSales        float64 `bson:"total_sales" json:"total_sales"`
	TotalTransactions int     `bson:"total_transactions" json:"total_transactions"`
	PointsEarned      float64 `bson:"points_earned" json:"points_earned"`
	StoreName         string  `json:"store_name,omitempty"`
	Location          string  `json:"location,omitempty"`
}

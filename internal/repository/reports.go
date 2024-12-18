package repository

import (
	"context"

	"wac-offline-payment/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMerchantReports() ([]models.Report, error) {
	// Aggregate reports from logs collection
	logsCollection := client.Database("wac-points").Collection("logs")
	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.M{
				"_id":                "$merchant_phone",
				"total_sales":        bson.M{"$sum": "$amount"},
				"total_transactions": bson.M{"$sum": 1},
				"points_earned":      bson.M{"$sum": bson.M{"$divide": []interface{}{"$amount", 25}}},
			}},
		},
	}

	cursor, err := logsCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode aggregated reports
	var reports []models.Report
	if err := cursor.All(context.Background(), &reports); err != nil {
		return nil, err
	}

	// Extract phone numbers
	phoneNumbers := []string{}
	for _, report := range reports {
		phoneNumbers = append(phoneNumbers, report.MerchantPhone)
	}

	// Fetch user details from the users repository
	users, err := GetMerchantsByPhoneNumbers(phoneNumbers)
	if err != nil {
		return nil, err
	}

	// Map user details by phone number
	userDetailsMap := make(map[string]models.Merchant)
	for _, user := range users {
		userDetailsMap[user.PhoneNumber] = user
	}

	// Merge user details into reports
	for i, report := range reports {
		if user, exists := userDetailsMap[report.MerchantPhone]; exists {
			reports[i].StoreName = user.StoreName
			reports[i].Location = user.Location
		} else {
			reports[i].StoreName = "Unknown"
			reports[i].Location = "Unknown"
		}
	}

	return reports, nil
}

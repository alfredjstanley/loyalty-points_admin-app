package repository

import (
	"context"

	"wac-offline-payment/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func SearchMerchants(query string) ([]models.Merchant, error) {
	collection := client.Database("olopo-points").Collection("users")

	// Search by name, location, or phone number
	filter := bson.M{
		"$or": []bson.M{
			{"store_name": bson.M{"$regex": query, "$options": "i"}},
			{"location": bson.M{"$regex": query, "$options": "i"}},
			{"phone_number": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var merchants []models.Merchant
	if err := cursor.All(context.Background(), &merchants); err != nil {
		return nil, err
	}

	return merchants, nil
}

func GetMerchantsByPhoneNumbers(phoneNumbers []string) ([]models.Merchant, error) {
	// Reference to the users collection
	collection := client.Database("olopo-points").Collection("users")

	// Build filter for phone numbers
	filter := bson.M{"phone_number": bson.M{"$in": phoneNumbers}}

	// Execute the query
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode the results into a slice of Merchant
	var merchants []models.Merchant
	if err := cursor.All(context.Background(), &merchants); err != nil {
		return nil, err
	}

	return merchants, nil
}

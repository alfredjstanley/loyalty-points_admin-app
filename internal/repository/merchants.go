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

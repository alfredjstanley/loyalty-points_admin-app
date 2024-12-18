package repository

import (
	"context"
	"time"

	"wac-offline-payment/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

// AddCounter adds a new counter for a merchant
func AddCounter(counter models.Counter) error {
	collection := client.Database("wac-points").Collection("counters")

	counter.CreatedAt = time.Now()
	counter.UpdatedAt = time.Now()

	_, err := collection.InsertOne(context.Background(), counter)
	return err
}

// GetCountersByMerchant fetches counters for a specific merchant
func GetCountersByMerchant(merchantPhone string) ([]models.Counter, error) {
	collection := client.Database("wac-points").Collection("counters")

	filter := bson.M{"merchant_phone": merchantPhone}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var counters []models.Counter
	if err := cursor.All(context.Background(), &counters); err != nil {
		return nil, err
	}

	return counters, nil
}

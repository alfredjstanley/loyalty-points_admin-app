package repository

import (
	"context"

	"wac-offline-payment/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetLogsWithPagination fetches transaction logs with pagination
func GetLogsWithPagination(page, limit int) ([]models.Log, int, error) {
	collection := client.Database("wac-points").Collection("logs")
	skip := (page - 1) * limit

	// Pagination options
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{"created_at", -1}})

	// Query logs
	cursor, err := collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Decode logs into a slice
	var logs []models.Log
	if err := cursor.All(context.Background(), &logs); err != nil {
		return nil, 0, err
	}

	// Get total log count
	total, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil
}

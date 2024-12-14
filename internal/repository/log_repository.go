package repository

import (
	"context"

	"wac-offline-payment/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetSuccessTransactionCount() (int, error) {
	collection := client.Database("wac-points").Collection("logs")
	count, err := collection.CountDocuments(context.Background(), bson.M{"status": "SUCCESS"})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func GetTotalTransactionAmount() (float64, error) {
	collection := client.Database("wac-points").Collection("logs")
	matchStage := bson.D{{"$match", bson.D{{"status", "SUCCESS"}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", nil}, {"totalAmount", bson.D{{"$sum", "$amount"}}}}}}

	cursor, err := collection.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.Background())

	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		return 0, err
	}
	if len(result) > 0 {
		return result[0]["totalAmount"].(float64), nil
	}
	return 0, nil
}

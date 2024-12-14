package repository

import (
	"context"
	"strconv"

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
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: "SUCCESS"}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "totalAmount", Value: bson.D{{Key: "$sum", Value: "$amount"}}}}}}

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

func SearchTransactionLogs(query string) ([]models.Log, error) {
	collection := client.Database("wac-points").Collection("logs")

	// Use regex to match query in UserPhone, MerchantPhone, InvoiceID, or Status
	filter := bson.M{
		"$or": []bson.M{
			{"user_phone": bson.M{"$regex": query, "$options": "i"}},
			{"merchant_phone": bson.M{"$regex": query, "$options": "i"}},
			{"amount": bson.M{"$regex": query, "$options": "i"}},
			{"invoice_id": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	// Check if query is a valid number (for amount field search)
	if amount, err := strconv.ParseFloat(query, 64); err == nil {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"amount": amount})
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var logs []models.Log
	if err := cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

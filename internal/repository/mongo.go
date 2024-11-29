package repository

import (
	"context"
	"log"
	"time"

	"wac-offline-payment/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Initialize MongoDB connection
func InitMongo(uri string) {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB!")
}

// SaveLog saves a log entry to MongoDB
func SaveLog(logEntry models.Log) error {
	collection := client.Database("wac-points").Collection("logs")
	logEntry.CreatedAt = time.Now()
	_, err := collection.InsertOne(context.TODO(), logEntry)
	return err
}

func SaveUser(user models.User) error {
	// log user to stdout
	log.Printf("Saving user: %v", user)
	collection := client.Database("olopo-points").Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Error inserting user into MongoDB: %v", err)
		return err
	}
	return nil
}

func FindUserByPhone(phoneNumber string) (*models.User, error) {
	collection := client.Database("olopo-points").Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), map[string]interface{}{
		"phone_number": phoneNumber,
	}).Decode(&user)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	collection := client.Database("olopo-points").Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUsersWithPagination(page, limit int) ([]models.User, int, error) {
	collection := client.Database("olopo-points").Collection("users")
	skip := (page - 1) * limit

	// Count total users
	total, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated users
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

package repository

import (
	"context"
	"log"
	"time"

	"wac-offline-payment/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func findUserByPhone(phoneNumber string) (*models.User, error) {
	collection := client.Database("olopo-points").Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), map[string]interface{}{
		"phone_number": phoneNumber,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No user found for phone number: %s", phoneNumber)
			return nil, nil
		}
		log.Printf("Error finding user: %v", err)
		return nil, err
	}
	return &user, nil
}

func findCounterByUsername(username string) (*models.Counter, error) {
	collection := client.Database("wac-points").Collection("counters")
	var counter models.Counter
	err := collection.FindOne(context.Background(), map[string]interface{}{
		"username": username,
	}).Decode(&counter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No counter found for username: %s", username)
			return nil, nil
		}
		log.Printf("Error finding counter: %v", err)
		return nil, err
	}
	return &counter, nil
}

func FindUserOrCounter(identifier string) (*models.User, *models.Counter, error) {
	user, err := findUserByPhone(identifier)
	if err != nil {
		return nil, nil, err
	}

	if user != nil {
		// User was found
		return user, nil, nil
	}

	// No user found, attempt to find counter
	counter, err := findCounterByUsername(identifier)
	if err != nil {
		return nil, nil, err
	}

	return nil, counter, nil
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

	// Fetch paginated users in reverse order (newest first)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}}) // Sort by `_id` in descending order
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

func UpdateMerchant(id primitive.ObjectID, update models.MerchantUpdate) error {
	collection := client.Database("olopo-points").Collection("users")

	updateFields := bson.M{}
	if update.StoreName != "" {
		updateFields["store_name"] = update.StoreName
	}
	if update.Location != "" {
		updateFields["location"] = update.Location
	}
	if update.Password != "" {
		updateFields["password"] = update.Password
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
	)
	return err
}

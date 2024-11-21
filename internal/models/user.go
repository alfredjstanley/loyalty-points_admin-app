package models

type User struct {
	StoreName   string `json:"store_name" bson:"store_name"`
	Location    string `json:"location" bson:"location"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Password    string `json:"password" bson:"password"`
}

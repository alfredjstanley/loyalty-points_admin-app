package models

type User struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	StoreName   string `json:"store_name" bson:"store_name"`
	Location    string `json:"location" bson:"location"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Password    string `json:"password" bson:"password"`
}

type Merchant struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	StoreName   string `json:"store_name" bson:"store_name"`
	Location    string `json:"location" bson:"location"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

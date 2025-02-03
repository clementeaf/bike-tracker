package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	Name           string              `bson:"name"`
	Email          string              `bson:"email"`
	Password       string              `bson:"password"`
	WalletBalance  float64             `bson:"wallet_balance"`
	LastSession    time.Time           `bson:"last_session"`
	LastBikeUsedID *primitive.ObjectID `bson:"last_bike_used_id,omitempty"`
}

func NewUser(name, email, password string) User {
	return User{
		ID:             primitive.NewObjectID(),
		Name:           name,
		Email:          email,
		Password:       password,
		WalletBalance:  0,
		LastSession:    time.Now(),
		LastBikeUsedID: nil,
	}
}

package wallet

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	WalletID  primitive.ObjectID `bson:"wallet_id" json:"wallet_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Type      string             `bson:"type" json:"type"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type Wallet struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Balance     float64            `bson:"balance" json:"balance"`
	LastUpdated time.Time          `bson:"last_updated" json:"last_updated"`
}

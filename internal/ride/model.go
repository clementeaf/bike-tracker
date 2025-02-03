package ride

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ride struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	BikeID      primitive.ObjectID `bson:"bike_id" json:"bike_id"`
	StartCoords []float64          `bson:"start_coords" json:"start_coords"`
	EndCoords   []float64          `bson:"end_coords,omitempty" json:"end_coords,omitempty"`
	Status      bool               `bson:"status" json:"status"` // true = on going, false = ended
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	FinalCost   float64            `bson:"final_cost,omitempty" json:"final_cost,omitempty"`
	BatteryLeft float64            `bson:"battery_left,omitempty" json:"battery_left,omitempty"`
}

type RideRequest struct {
	BikeID      string    `json:"bike_id"`
	StartCoords []float64 `json:"start_coords"`
}

type EndRideRequest struct {
	RideID    string    `json:"ride_id"`
	EndCoords []float64 `json:"end_coords"`
	Battery   float64   `json:"battery"`
}

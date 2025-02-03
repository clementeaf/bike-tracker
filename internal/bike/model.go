package bike

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bike struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	BatteryLevel      float64              `bson:"battery_level" json:"battery_level"`
	Latitude          float64              `bson:"latitude" json:"latitude"`
	Longitude         float64              `bson:"longitude" json:"longitude"`
	Status            int                  `bson:"status" json:"status"`
	LastUsedAt        time.Time            `bson:"last_used_at" json:"last_used_at"`
	UserHistory       []primitive.ObjectID `bson:"user_history,omitempty" json:"user_history"`
	TotalUsageMinutes float64              `bson:"total_usage_minutes" json:"total_usage_minutes"`
	TotalEarnings     float64              `bson:"total_earnings" json:"total_earnings"`
	LastMaintenance   time.Time            `bson:"last_maintenance" json:"last_maintenance"`
	NextMaintenance   time.Time            `bson:"next_maintenance" json:"next_maintenance"`
	OperationalSince  time.Time            `bson:"operational_since" json:"operational_since"`
}

type TripCost struct {
	DistanceMeters float64 `json:"distance_meters"`
	TotalCost      float64 `json:"total_cost"`
}

type BikeRequest struct {
	BikeID primitive.ObjectID `bson:"bike_id" json:"bike_id"`
	Status int                `bson:"status" json:"status"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
}

const (
	StatusFree        = 1 // Free
	StatusInUse       = 2 // Being used
	StatusMaintenance = 3 // Under maintenance
	StatusNoBattery   = 4 // Batery 0
	StatusReserved    = 5 // Reserved
)

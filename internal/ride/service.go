package ride

import (
	"context"
	"errors"
	"time"

	"github.com/clementeaf/bike-tracker/internal/bike"
	"github.com/clementeaf/bike-tracker/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// New ride in Database
func insertRide(ride Ride) error {
	_, err := database.GetCollection("rides").InsertOne(context.Background(), ride)
	return err
}

// Get ride by ID
func getRideByID(rideID string) (Ride, error) {
	var ride Ride
	rideObjectID, _ := primitive.ObjectIDFromHex(rideID)

	err := database.GetCollection("rides").FindOne(context.Background(), bson.M{"_id": rideObjectID}).Decode(&ride)
	if err != nil {
		return ride, errors.New("viaje no encontrado")
	}

	return ride, nil
}

// Get all rides
func getAllRides() ([]Ride, error) {
	var rides []Ride

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.GetCollection("rides").Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("error al consultar rides: " + err.Error())
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rides); err != nil {
		return nil, errors.New("error al procesar rides: " + err.Error())
	}

	return rides, nil
}

// Get rides by status
func getRidesByStatus(status bool) ([]Ride, error) {
	var rides []Ride
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.GetCollection("rides").Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &rides); err != nil {
		return nil, err
	}

	return rides, nil
}

// Calculate ride cost by time
func calculateCost(minutes float64) float64 {
	return minutes * 0.5 // $0.5 por minuto
}

// Validate if bike is available
func validateBike(bikeID string) (bike.Bike, error) {
	var bicycle bike.Bike
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bikeObjectID, err := primitive.ObjectIDFromHex(bikeID)
	if err != nil {
		return bicycle, errors.New("ID de bicicleta inválido")
	}

	err = database.GetCollection("bikes").FindOne(ctx, bson.M{"_id": bikeObjectID}).Decode(&bicycle)
	if err != nil {
		return bicycle, errors.New("bicicleta no encontrada")
	}

	if bicycle.Status != 1 {
		return bicycle, errors.New("bicicleta no está disponible (no está libre)")
	}

	if bicycle.BatteryLevel < 20 {
		return bicycle, errors.New("bicicleta inactiva por nivel de batería bajo")
	}

	return bicycle, nil
}

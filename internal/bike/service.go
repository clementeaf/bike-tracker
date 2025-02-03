package bike

import (
	"context"
	"errors"
	"time"

	"github.com/clementeaf/bike-tracker/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a new bike in the database
func RegisterBike() (*Bike, error) {
	bike := &Bike{
		ID:                primitive.NewObjectID(),
		BatteryLevel:      100,
		Latitude:          0,
		Longitude:         0,
		Status:            StatusFree, // Siempre un entero (1)
		LastUsedAt:        time.Time{},
		UserHistory:       []primitive.ObjectID{},
		TotalUsageMinutes: 0,
		TotalEarnings:     0,
		LastMaintenance:   time.Time{},
		NextMaintenance:   time.Time{},
		OperationalSince:  time.Now(),
	}

	collection := database.GetCollection("bikes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, bike)
	if err != nil {
		return nil, err
	}

	return bike, nil
}

// Get all bikes if state = 1
func GetAvailableBikes() ([]Bike, error) {
	var bikes []Bike

	// Contexto con timeout para la consulta
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Consulta para buscar bicicletas con estado "Libre" (StatusFree)
	cursor, err := database.GetCollection("bikes").Find(ctx, bson.M{"status": StatusFree})
	if err != nil {
		return nil, errors.New("error al consultar bicicletas disponibles: " + err.Error())
	}
	defer cursor.Close(ctx)

	// Decodificar bicicletas en el slice
	if err := cursor.All(ctx, &bikes); err != nil {
		return nil, errors.New("error al procesar bicicletas disponibles: " + err.Error())
	}

	return bikes, nil
}

// Update bike status
func UpdateBikeStatus(bikeID string, userID string, status int) error {
	if status < StatusFree || status > StatusReserved {
		return errors.New("estado inválido")
	}

	bikeCollection := database.GetCollection("bikes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bikeObjectID, err := primitive.ObjectIDFromHex(bikeID)
	if err != nil {
		return errors.New("ID de bicicleta inválido")
	}

	var bike Bike
	err = bikeCollection.FindOne(ctx, bson.M{"_id": bikeObjectID}).Decode(&bike)
	if err != nil {
		return errors.New("bicicleta no encontrada")
	}

	updateFields := bson.M{
		"status":       status,
		"last_used_at": time.Now(),
	}

	if status == StatusInUse {
		userObjectID, _ := primitive.ObjectIDFromHex(userID)
		updateFields["user_history"] = append(bike.UserHistory, userObjectID)
	}

	_, err = bikeCollection.UpdateOne(ctx, bson.M{"_id": bikeObjectID}, bson.M{"$set": updateFields})
	if err != nil {
		return errors.New("no se pudo actualizar el estado de la bicicleta")
	}

	return nil
}

// Trip cost by time
func CalculateTripCost(durationMinutes float64) TripCost {
	const ratePerMinute = 0.50 // Tarifa en USD por minuto
	totalCost := durationMinutes * ratePerMinute

	return TripCost{
		DistanceMeters: durationMinutes * 100, // Suponiendo 100 metros por minuto
		TotalCost:      totalCost,
	}
}

// Get all bikes
func GetAllBikes() ([]Bike, error) {
	var bikes []Bike

	// Crear un contexto con timeout para la consulta
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtener todas las bicicletas (sin filtro)
	cursor, err := database.GetCollection("bikes").Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("error al consultar bicicletas: " + err.Error())
	}
	defer cursor.Close(ctx)

	// Decodificar los resultados en la slice de bicicletas
	if err := cursor.All(ctx, &bikes); err != nil {
		return nil, errors.New("error al procesar bicicletas: " + err.Error())
	}

	return bikes, nil
}

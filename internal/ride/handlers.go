package ride

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/clementeaf/bike-tracker/internal/bike"
	"github.com/clementeaf/bike-tracker/internal/wallet"
	"github.com/clementeaf/bike-tracker/pkg/auth"
	"github.com/clementeaf/bike-tracker/pkg/database"
	httpresponse "github.com/clementeaf/bike-tracker/pkg/http"
	"github.com/clementeaf/bike-tracker/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST New Ride initiate
func handleStartRide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método no permitido"})
		return
	}

	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "No autorizado"})
		logger.Error("handleStartRide - Usuario no autorizado", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var rideRequest RideRequest
	if err := json.NewDecoder(r.Body).Decode(&rideRequest); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
		logger.Error("handleStartRide - Error al decodificar JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if rideRequest.BikeID == "" || len(rideRequest.StartCoords) != 2 {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Faltan datos requeridos (BikeID o Coordenadas)"})
		logger.Error("handleStartRide - Campos faltantes", map[string]interface{}{
			"bike_id":      rideRequest.BikeID,
			"start_coords": rideRequest.StartCoords,
		})
		return
	}

	bikeObject, err := validateBike(rideRequest.BikeID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		logger.Error("handleStartRide - Bicicleta no válida", map[string]interface{}{
			"bike_id": rideRequest.BikeID,
			"error":   err.Error(),
		})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "ID de usuario no válido"})
		logger.Error("handleStartRide - Error al convertir userID", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return
	}

	if err := wallet.DeductRideFee(userObjectID.Hex()); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusPaymentRequired, map[string]string{"error": "Saldo insuficiente en la wallet"})
		logger.Error("handleStartRide - Wallet insuficiente", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return
	}

	if err := bike.UpdateBikeStatus(bikeObject.ID.Hex(), userID, bike.StatusInUse); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error al actualizar el estado de la bicicleta"})
		logger.Error("handleStartRide - Error al actualizar estado de la bicicleta", map[string]interface{}{
			"bike_id": bikeObject.ID.Hex(),
			"user_id": userID,
			"status":  bike.StatusInUse,
			"error":   err.Error(),
		})
		return
	}

	ride := Ride{
		ID:          primitive.NewObjectID(),
		UserID:      userObjectID,
		BikeID:      bikeObject.ID,
		StartCoords: rideRequest.StartCoords,
		Status:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := insertRide(ride); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error al iniciar el viaje"})
		logger.Error("handleStartRide - Error al crear ride", map[string]interface{}{
			"ride":  ride,
			"error": err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusCreated, ride)
	logger.Info("handleStartRide - Ride iniciado exitosamente", map[string]interface{}{
		"ride_id":   ride.ID.Hex(),
		"user_id":   userID,
		"bike_id":   bikeObject.ID.Hex(),
		"start_lat": rideRequest.StartCoords[0],
		"start_lon": rideRequest.StartCoords[1],
	})
}

// POST Ride end
func handleEndRide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	// Validar y obtener el ID del usuario autenticado desde el token JWT
	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "No autorizado",
		})
		logger.Error("handleEndRide - Usuario no autenticado", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var req EndRideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Datos inválidos",
		})
		logger.Error("handleEndRide - JSON inválido", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	rideID, err := primitive.ObjectIDFromHex(req.RideID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "ID de viaje inválido",
		})
		logger.Error("handleEndRide - ID de viaje inválido", map[string]interface{}{
			"ride_id": req.RideID,
			"error":   err.Error(),
		})
		return
	}

	var ride Ride
	if err := database.GetCollection("rides").FindOne(context.Background(), bson.M{"_id": rideID}).Decode(&ride); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "Viaje no encontrado",
		})
		logger.Error("handleEndRide - Viaje no encontrado", map[string]interface{}{
			"ride_id": req.RideID,
			"error":   err.Error(),
		})
		return
	}

	// Validar que el viaje pertenece al usuario autenticado
	if ride.UserID.Hex() != userID {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "No autorizado para finalizar este viaje",
		})
		logger.Error("handleEndRide - Usuario no autorizado para finalizar el viaje", map[string]interface{}{
			"ride_id": ride.ID.Hex(),
			"user_id": userID,
		})
		return
	}

	duration := time.Since(ride.CreatedAt).Minutes()
	finalCost := calculateCost(duration)

	var bike bike.Bike
	if err := database.GetCollection("bikes").FindOne(context.Background(), bson.M{"_id": ride.BikeID}).Decode(&bike); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Bicicleta no encontrada",
		})
		logger.Error("handleEndRide - Bicicleta no encontrada", map[string]interface{}{
			"bike_id": ride.BikeID.Hex(),
			"error":   err.Error(),
		})
		return
	}

	batteryConsumptionPerMinute := 2.0
	totalBatteryConsumed := batteryConsumptionPerMinute * duration
	batteryLeft := bike.BatteryLevel - totalBatteryConsumed
	if batteryLeft < 0 {
		batteryLeft = 0
	}

	updateRide := bson.M{
		"$set": bson.M{
			"end_coords":   req.EndCoords,
			"status":       false,
			"updated_at":   time.Now(),
			"final_cost":   finalCost,
			"battery_left": batteryLeft,
		},
	}

	if _, err := database.GetCollection("rides").UpdateOne(context.Background(), bson.M{"_id": rideID}, updateRide); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "No se pudo actualizar el viaje",
		})
		logger.Error("handleEndRide - Error al actualizar viaje", map[string]interface{}{
			"ride_id": ride.ID.Hex(),
			"error":   err.Error(),
		})
		return
	}

	updateBike := bson.M{
		"$set": bson.M{
			"status":        1,
			"battery_level": batteryLeft,
			"latitude":      req.EndCoords[0],
			"longitude":     req.EndCoords[1],
			"last_used_at":  time.Now(),
		},
	}

	if _, err := database.GetCollection("bikes").UpdateOne(context.Background(), bson.M{"_id": ride.BikeID}, updateBike); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "No se pudo actualizar la bicicleta",
		})
		logger.Error("handleEndRide - Error al actualizar bicicleta", map[string]interface{}{
			"bike_id": ride.BikeID.Hex(),
			"error":   err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status": "finalizado",
	})

	logger.Info("handleEndRide - Viaje finalizado con éxito", map[string]interface{}{
		"ride_id": ride.ID.Hex(),
		"bike_id": ride.BikeID.Hex(),
	})
}

// GET all rides
func handleGetAllRides(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	rides, err := getAllRides()
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		logger.Error("handleGetAllRides - Error al obtener todos los rides", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, rides)
	logger.Info("handleGetAllRides - Rides obtenidos exitosamente", map[string]interface{}{
		"rides_count": len(rides),
	})
}

// GET ride by ID.
func handleGetRideByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	rideID := strings.TrimPrefix(r.URL.Path, "/rides/")
	if rideID == "" {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Falta el ID del viaje",
		})
		return
	}

	ride, err := getRideByID(rideID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		logger.Error("handleGetRideByID - Error al obtener el ride", map[string]interface{}{
			"ride_id": rideID,
			"error":   err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, ride)
	logger.Info("handleGetRideByID - Ride obtenido exitosamente", map[string]interface{}{
		"ride_id": rideID,
	})
}

// GET rides if status = true
func handleGetActiveRides(w http.ResponseWriter, r *http.Request) {
	rides, err := getRidesByStatus(true)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al obtener viajes en curso",
		})
		logger.Error("handleGetActiveRides - Error al consultar viajes en curso", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, rides)
	logger.Info("handleGetActiveRides - Viajes activos obtenidos exitosamente", map[string]interface{}{
		"active_rides_count": len(rides),
	})
}

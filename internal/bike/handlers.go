package bike

import (
	"encoding/json"
	"net/http"

	"github.com/clementeaf/bike-tracker/pkg/auth"
	httpresponse "github.com/clementeaf/bike-tracker/pkg/http"
	"github.com/clementeaf/bike-tracker/pkg/logger"
)

// Header & content types
const (
	ContentType  = "application/json"
	UserIDHeader = "Authenticated-User-ID"
)

// POST New Bike
func HandleRegisterBike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	bike, err := RegisterBike()
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al generar la bicicleta",
		})
		logger.Error("POST /bikes/new - Error al generar bicicleta", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusCreated, bike)
	logger.Info("POST /bikes/new - Bicicleta generada exitosamente", map[string]interface{}{
		"bike_id": bike.ID.Hex(),
	})
}

// GET Available Bikes - Status = 1
func HandleGetAvailableBikes(w http.ResponseWriter, r *http.Request) {
	// Validar método HTTP
	if r.Method != http.MethodGet {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	// Llamar al servicio para obtener las bicicletas disponibles
	bikes, err := GetAvailableBikes()
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al obtener bicicletas disponibles",
		})
		logger.Error("GET /bikes/available - Error al consultar servicio", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Manejar el caso de que no haya bicicletas disponibles
	if len(bikes) == 0 {
		httpresponse.SendJSONResponse(w, http.StatusOK, []interface{}{}) // Respuesta vacía
		logger.Info("GET /bikes/available - No hay bicicletas disponibles", map[string]interface{}{
			"available_bikes": 0,
		})
		return
	}

	// Enviar la lista de bicicletas disponibles
	httpresponse.SendJSONResponse(w, http.StatusOK, bikes)
	logger.Info("GET /bikes/available - Bicicletas disponibles devueltas", map[string]interface{}{
		"bikes_count": len(bikes),
	})
}

// PUT Bike status
func HandleUpdateBikeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "No autorizado",
		})
		logger.Error("PUT /bikes/status - Usuario no autenticado", nil)
		return
	}

	var input struct {
		BikeID string `json:"bike_id"`
		Status int    `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Error en el formato del JSON",
		})
		logger.Error("PUT /bikes/status - JSON inválido", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if input.BikeID == "" || input.Status <= 0 {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Faltan campos requeridos (bike_id, status)",
		})
		logger.Error("PUT /bikes/status - Campos faltantes", map[string]interface{}{
			"bike_id": input.BikeID,
			"status":  input.Status,
		})
		return
	}

	err = UpdateBikeStatus(input.BikeID, userID, input.Status)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "No se pudo actualizar el estado de la bicicleta",
		})
		logger.Error("PUT /bikes/status - Error al actualizar bicicleta", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Estado actualizado correctamente",
	})
	logger.Info("PUT /bikes/status - Estado actualizado exitosamente", map[string]interface{}{
		"bike_id": input.BikeID,
		"user_id": userID,
	})
}

// GET Bikes
func HandleGetAllBikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	bikes, err := GetAllBikes()
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al obtener bicicletas",
		})
		logger.Error("GET /bikes/all - Error al consultar servicio", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if len(bikes) == 0 {
		httpresponse.SendJSONResponse(w, http.StatusOK, []interface{}{})
		logger.Info("GET /bikes/all - No hay bicicletas registradas", map[string]interface{}{
			"total_bikes": 0,
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, bikes)
	logger.Info("GET /bikes/all - Bicicletas devueltas", map[string]interface{}{
		"total_bikes": len(bikes),
	})
}

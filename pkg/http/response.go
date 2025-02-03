package httpresponse

import (
	"encoding/json"
	"net/http"

	"github.com/clementeaf/bike-tracker/pkg/logger"
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("SendJSONResponse - Error al codificar respuesta JSON", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

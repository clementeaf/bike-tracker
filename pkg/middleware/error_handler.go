package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/clementeaf/bike-tracker/pkg/logger"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func JSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := ErrorResponse{
		Status:  status,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		logger.Error("Error al codificar la respuesta JSON de error", map[string]interface{}{
			"original_error": message,
			"encoding_error": err.Error(),
		})
	}
}

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Error inesperado en la aplicación", map[string]interface{}{
					"error": err,
				})
				JSONError(w, http.StatusInternalServerError, "Ocurrió un error interno en el servidor")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/clementeaf/bike-tracker/pkg/auth"
	"github.com/clementeaf/bike-tracker/pkg/logger"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No autorizado: falta el token JWT", http.StatusUnauthorized)
			logger.Error("Solicitud no autorizada: falta el token JWT", nil)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "No autorizado: formato de token inválido", http.StatusUnauthorized)
			logger.Error("Solicitud no autorizada: formato de token inválido", nil)
			return
		}

		claims, err := auth.ValidateToken(tokenParts[1])
		if err != nil {
			http.Error(w, "No autorizado: token inválido o expirado", http.StatusUnauthorized)
			logger.Error("Solicitud no autorizada: token inválido", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		r.Header.Set("Authenticated-User-ID", claims.UserID)
		next.ServeHTTP(w, r)
	})
}

func ApplyMiddlewares(handler http.Handler) http.Handler {
	return ErrorMiddleware(handler)
}

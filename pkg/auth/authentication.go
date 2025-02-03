package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GET User id from JWT
func GetAuthenticatedUserID(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("token no encontrado")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("formato de token inválido")
	}
	tokenString := parts[1]

	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", errors.New("token inválido o expirado")
	}

	return claims.UserID, nil
}

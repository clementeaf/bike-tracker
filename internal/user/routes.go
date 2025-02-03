package user

import (
	"net/http"

	"github.com/clementeaf/bike-tracker/pkg/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {
	// Public routes (no jwt)
	mux.HandleFunc("/users/register", handleRegister)
	mux.HandleFunc("/users/login", handleLogin)

	// Protected routes (with jwt)
	mux.Handle("/users/me", middleware.AuthMiddleware(http.HandlerFunc(handleGetMe)))
	mux.Handle("/users/me/update", middleware.AuthMiddleware(http.HandlerFunc(handleUpdateUser)))
	mux.Handle("/users/me/delete", middleware.AuthMiddleware(http.HandlerFunc(handleDeleteUser)))
}

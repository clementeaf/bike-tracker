package api

import (
	"net/http"

	"github.com/clementeaf/bike-tracker/internal/bike"
	"github.com/clementeaf/bike-tracker/internal/ride"
	"github.com/clementeaf/bike-tracker/internal/user"
	"github.com/clementeaf/bike-tracker/internal/wallet"
	"github.com/clementeaf/bike-tracker/pkg/middleware"
)

// NewRouter crea un enrutador HTTP y registra todas las rutas principales
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Registrar rutas de usuarios
	user.RegisterRoutes(mux)

	// Registrar rutas de viajes
	ride.RegisterRoutes(mux)

	// Registrar rutas de wallet
	wallet.RegisterRoutes(mux)

	// Registrar rutas de bicicletas
	bike.RegisterRoutes(mux)

	// Ruta raíz (para manejar rutas no encontradas)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error": "Ruta no encontrada"}`, http.StatusNotFound)
	})

	// Aplicar middleware de logging global
	return middleware.ApplyMiddlewares(mux)
}

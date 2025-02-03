package api

import (
	"log"
	"net/http"
	"os"
)

// StartServer inicializa y arranca el servidor HTTP
func StartServer() {
	// Crear el enrutador
	router := NewRouter()

	// Leer el puerto desde las variables de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar el servidor
	log.Printf("🚀 Servidor corriendo en el puerto %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

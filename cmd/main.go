package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/clementeaf/bike-tracker/internal/api"
	"github.com/clementeaf/bike-tracker/pkg/config"
	"github.com/clementeaf/bike-tracker/pkg/database"
	"github.com/clementeaf/bike-tracker/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}
	// Cargar variables de entorno
	config.LoadEnv()

	// Conectar a MongoDB
	database.ConnectMongo()

	// Inicializar logger
	logger.InitLogger()

	// Crear el enrutador principal
	router := api.NewRouter() // Ahora retorna http.Handler

	// Arrancar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("🚀 Servidor corriendo en el puerto", port)
	http.ListenAndServe(":"+port, router) // Se mantiene porque http.ListenAndServe acepta http.Handler
}

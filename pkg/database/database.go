package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectMongo() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("La variable de entorno MONGO_URI no está configurada. Verifica tu archivo .env")
	}

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("No se pudo realizar el ping a MongoDB: %v", err)
	}

	fmt.Println("Conexión exitosa a MongoDB")
	Client = client
}

func GetCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		log.Fatal("La variable de entorno MONGO_DB_NAME no está configurada. Verifica tu archivo .env")
	}

	return Client.Database(dbName).Collection(collectionName)
}

func DisconnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Fatalf("Error al desconectar de MongoDB: %v", err)
	}
	log.Println("Desconexión exitosa de MongoDB")
}

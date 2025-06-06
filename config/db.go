
package config

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoURI apunta a tu servidor de Mongo. 
// Puedes cambiar "mongodb://localhost:27017" por la URI de tu Atlas o tu Docker local.
const MongoURI = "mongodb://localhost:27017" 
const DBName = "api_rest_db"

// ClientMongo es la instancia global del cliente de MongoDB
var ClientMongo *mongo.Client

// ConnectBD inicia la conexión a MongoDB
func ConnectDB() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(MongoURI)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("Error al conectar a MongoDB: ", err)
    }

    // es para verificar la conexión
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Error al hacer ping a MongoDB: ", err)
    }

    fmt.Println("Conectado a MongoDB en ", MongoURI)
    ClientMongo = client
}


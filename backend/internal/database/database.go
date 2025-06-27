package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"home-server-hub/internal/config"
)

// connectDB estabelece conexão com o MongoDB
func Connect(dbConfig config.DatabaseConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dbConfig.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verificar conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Conectado ao MongoDB!")
	return client.Database(dbConfig.DatabaseName), nil
}

// disconnectDB encerra a conexão com o MongoDB
func Disconnect(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Client().Disconnect(ctx); err != nil {
		log.Fatalf("Erro ao desconectar do MongoDB: %v", err)
	}
	log.Println("Desconectado do MongoDB")
}
